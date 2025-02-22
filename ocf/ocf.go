// Package ocf implements encoding and decoding of Avro Object Container Files as defined by the Avro specification.
//
// See the Avro specification for an understanding of Avro: http://avro.apache.org/docs/current/
package ocf

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kjuulh/avro/v2"
	"github.com/kjuulh/avro/v2/internal/bytesx"
)

const (
	schemaKey = "avro.schema"
	codecKey  = "avro.codec"
)

var magicBytes = [4]byte{'O', 'b', 'j', 1}

// HeaderSchema is the Avro schema of a container file header.
var HeaderSchema = avro.MustParse(`{
	"type": "record", 
	"name": "org.apache.avro.file.Header",
	"fields": [
		{"name": "magic", "type": {"type": "fixed", "name": "Magic", "size": 4}},
		{"name": "meta", "type": {"type": "map", "values": "bytes"}},
		{"name": "sync", "type": {"type": "fixed", "name": "Sync", "size": 16}}
	]
}`)

// Header represents an Avro container file header.
type Header struct {
	Magic [4]byte           `avro:"magic"`
	Meta  map[string][]byte `avro:"meta"`
	Sync  [16]byte          `avro:"sync"`
}

// Decoder reads and decodes Avro values from a container file.
type Decoder struct {
	reader      *avro.Reader
	resetReader *bytesx.ResetReader
	decoder     *avro.Decoder
	meta        map[string][]byte
	sync        [16]byte

	codec Codec

	count int64
}

// NewDecoder returns a new decoder that reads from reader r.
func NewDecoder(r io.Reader) (*Decoder, error) {
	reader := avro.NewReader(r, 1024)

	h, err := readHeader(reader)
	if err != nil {
		return nil, fmt.Errorf("decoder: %w", err)
	}

	decReader := bytesx.NewResetReader([]byte{})

	return &Decoder{
		reader:      reader,
		resetReader: decReader,
		decoder:     avro.NewDecoderForSchema(h.Schema, decReader),
		meta:        h.Meta,
		sync:        h.Sync,
		codec:       h.Codec,
	}, nil
}

// Metadata returns the header metadata.
func (d *Decoder) Metadata() map[string][]byte {
	return d.meta
}

// HasNext determines if there is another value to read.
func (d *Decoder) HasNext() bool {
	if d.count <= 0 {
		count := d.readBlock()
		d.count = count
	}

	if d.reader.Error != nil {
		return false
	}

	return d.count > 0
}

// Decode reads the next Avro encoded value from its input and stores it in the value pointed to by v.
func (d *Decoder) Decode(v any) error {
	if d.count <= 0 {
		return errors.New("decoder: no data found, call HasNext first")
	}

	d.count--

	return d.decoder.Decode(v)
}

// Error returns the last reader error.
func (d *Decoder) Error() error {
	if errors.Is(d.reader.Error, io.EOF) {
		return nil
	}

	return d.reader.Error
}

func (d *Decoder) readBlock() int64 {
	count := d.reader.ReadLong()
	size := d.reader.ReadLong()

	if count > 0 {
		data := make([]byte, size)
		d.reader.Read(data)

		data, err := d.codec.Decode(data)
		if err != nil {
			d.reader.Error = err
		}

		d.resetReader.Reset(data)
	}

	var sync [16]byte
	d.reader.Read(sync[:])
	if d.sync != sync && !errors.Is(d.reader.Error, io.EOF) {
		d.reader.Error = errors.New("decoder: invalid block")
	}

	return count
}

type encoderConfig struct {
	BlockLength      int
	CodecName        CodecName
	CodecCompression int
	Metadata         map[string][]byte
	Sync             [16]byte
	EncodingConfig   avro.API
}

// EncoderFunc represents an configuration function for Encoder.
type EncoderFunc func(cfg *encoderConfig)

// WithBlockLength sets the block length on the encoder.
func WithBlockLength(length int) EncoderFunc {
	return func(cfg *encoderConfig) {
		cfg.BlockLength = length
	}
}

// WithCodec sets the compression codec on the encoder.
func WithCodec(codec CodecName) EncoderFunc {
	return func(cfg *encoderConfig) {
		cfg.CodecName = codec
	}
}

// WithCompressionLevel sets the compression codec to deflate and
// the compression level on the encoder.
func WithCompressionLevel(compLvl int) EncoderFunc {
	return func(cfg *encoderConfig) {
		cfg.CodecName = Deflate
		cfg.CodecCompression = compLvl
	}
}

// WithMetadata sets the metadata on the encoder header.
func WithMetadata(meta map[string][]byte) EncoderFunc {
	return func(cfg *encoderConfig) {
		cfg.Metadata = meta
	}
}

// WithSyncBlock sets the sync block.
func WithSyncBlock(sync [16]byte) EncoderFunc {
	return func(cfg *encoderConfig) {
		cfg.Sync = sync
	}
}

// WithEncodingConfig sets the value encoder config on the OCF encoder.
func WithEncodingConfig(wCfg avro.API) EncoderFunc {
	return func(cfg *encoderConfig) {
		cfg.EncodingConfig = wCfg
	}
}

// Encoder writes Avro container file to an output stream.
type Encoder struct {
	writer  *avro.Writer
	buf     *bytes.Buffer
	encoder *avro.Encoder
	sync    [16]byte

	codec Codec

	blockLength int
	count       int
}

// NewEncoder returns a new encoder that writes to w using schema s.
//
// If the writer is an existing ocf file, it will append data using the
// existing schema.
func NewEncoder(s string, w io.Writer, opts ...EncoderFunc) (*Encoder, error) {
	cfg := encoderConfig{
		BlockLength:      100,
		CodecName:        Null,
		CodecCompression: -1,
		Metadata:         map[string][]byte{},
		EncodingConfig:   avro.DefaultConfig,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	switch file := w.(type) {
	case nil:
		return nil, errors.New("writer cannot be nil")
	case *os.File:
		info, err := file.Stat()
		if err != nil {
			return nil, err
		}

		if info.Size() > 0 {
			reader := avro.NewReader(file, 1024)
			h, err := readHeader(reader)
			if err != nil {
				return nil, err
			}
			if err = skipToEnd(reader, h.Sync); err != nil {
				return nil, err
			}

			writer := avro.NewWriter(w, 512, avro.WithWriterConfig(cfg.EncodingConfig))
			buf := &bytes.Buffer{}
			e := &Encoder{
				writer:      writer,
				buf:         buf,
				encoder:     cfg.EncodingConfig.NewEncoder(h.Schema, buf),
				sync:        h.Sync,
				codec:       h.Codec,
				blockLength: cfg.BlockLength,
			}
			return e, nil
		}
	}

	schema, err := avro.Parse(s)
	if err != nil {
		return nil, err
	}

	cfg.Metadata[schemaKey] = []byte(schema.String())
	cfg.Metadata[codecKey] = []byte(cfg.CodecName)
	header := Header{
		Magic: magicBytes,
		Meta:  cfg.Metadata,
	}
	header.Sync = cfg.Sync
	if header.Sync == [16]byte{} {
		_, _ = rand.Read(header.Sync[:])
	}

	codec, err := resolveCodec(cfg.CodecName, cfg.CodecCompression)
	if err != nil {
		return nil, err
	}

	writer := avro.NewWriter(w, 512, avro.WithWriterConfig(cfg.EncodingConfig))
	writer.WriteVal(HeaderSchema, header)
	if err = writer.Flush(); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	e := &Encoder{
		writer:      writer,
		buf:         buf,
		encoder:     cfg.EncodingConfig.NewEncoder(schema, buf),
		sync:        header.Sync,
		codec:       codec,
		blockLength: cfg.BlockLength,
	}
	return e, nil
}

// Write v to the internal buffer. This method skips the internal encoder and
// therefore the caller is responsible for encoding the bytes. No error will be
// thrown if the bytes does not conform to the schema given to NewEncoder, but
// the final ocf data will be corrupted.
func (e *Encoder) Write(p []byte) (n int, err error) {
	n, err = e.buf.Write(p)
	if err != nil {
		return n, err
	}

	e.count++
	if e.count >= e.blockLength {
		if err = e.writerBlock(); err != nil {
			return n, err
		}
	}

	return n, e.writer.Error
}

// Encode writes the Avro encoding of v to the stream.
func (e *Encoder) Encode(v any) error {
	if err := e.encoder.Encode(v); err != nil {
		return err
	}

	e.count++
	if e.count >= e.blockLength {
		if err := e.writerBlock(); err != nil {
			return err
		}
	}

	return e.writer.Error
}

// Flush flushes the underlying writer.
func (e *Encoder) Flush() error {
	if e.count == 0 {
		return nil
	}

	if err := e.writerBlock(); err != nil {
		return err
	}

	return e.writer.Error
}

// Close closes the encoder, flushing the writer.
func (e *Encoder) Close() error {
	return e.Flush()
}

func (e *Encoder) writerBlock() error {
	e.writer.WriteLong(int64(e.count))

	b := e.codec.Encode(e.buf.Bytes())

	e.writer.WriteLong(int64(len(b)))
	_, _ = e.writer.Write(b)

	_, _ = e.writer.Write(e.sync[:])

	e.count = 0
	e.buf.Reset()
	return e.writer.Flush()
}

type ocfHeader struct {
	Schema avro.Schema
	Codec  Codec
	Meta   map[string][]byte
	Sync   [16]byte
}

func readHeader(reader *avro.Reader) (*ocfHeader, error) {
	var h Header
	reader.ReadVal(HeaderSchema, &h)
	if reader.Error != nil {
		return nil, fmt.Errorf("unexpected error: %w", reader.Error)
	}

	if h.Magic != magicBytes {
		return nil, errors.New("invalid avro file")
	}
	schema, err := avro.Parse(string(h.Meta[schemaKey]))
	if err != nil {
		return nil, err
	}

	codec, err := resolveCodec(CodecName(h.Meta[codecKey]), -1)
	if err != nil {
		return nil, err
	}

	return &ocfHeader{
		Schema: schema,
		Codec:  codec,
		Meta:   h.Meta,
		Sync:   h.Sync,
	}, nil
}

func skipToEnd(reader *avro.Reader, sync [16]byte) error {
	for {
		_ = reader.ReadLong()
		if errors.Is(reader.Error, io.EOF) {
			return nil
		}
		size := reader.ReadLong()
		reader.SkipNBytes(int(size))
		if reader.Error != nil {
			return reader.Error
		}

		var synMark [16]byte
		reader.Read(synMark[:])
		if sync != synMark && !errors.Is(reader.Error, io.EOF) {
			reader.Error = errors.New("invalid block")
		}
	}
}
