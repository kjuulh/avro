package something

// Code generated by avro/gen. DO NOT EDIT.

import (
	"math/big"
	"time"

	"github.com/kjuulh/avro/v2"
)

// InnerRecord is a generated struct.
type InnerRecord struct {
	InnerJustBytes                   []byte    `avro:"innerJustBytes"`
	InnerPrimitiveNullableArrayUnion *[]string `avro:"innerPrimitiveNullableArrayUnion"`
}

var schemaInnerRecord = avro.MustParse(`{"name":"a.c.InnerRecord","type":"record","fields":[{"name":"innerJustBytes","type":"bytes"},{"name":"innerPrimitiveNullableArrayUnion","type":["null",{"type":"array","items":"string"}]}]}`)

// Schema returns the schema for InnerRecord.
func (o *InnerRecord) Schema() avro.Schema {
	return schemaInnerRecord
}

// Unmarshal decodes b into the receiver.
func (o *InnerRecord) Unmarshal(b []byte) error {
	return avro.Unmarshal(o.Schema(), b, o)
}

// Marshal encodes the receiver.
func (o *InnerRecord) Marshal() ([]byte, error) {
	return avro.Marshal(o.Schema(), o)
}

// RecordInMap is a generated struct.
type RecordInMap struct {
	Name string `avro:"name"`
}

var schemaRecordInMap = avro.MustParse(`{"name":"a.b.RecordInMap","type":"record","fields":[{"name":"name","type":"string"}]}`)

// Schema returns the schema for RecordInMap.
func (o *RecordInMap) Schema() avro.Schema {
	return schemaRecordInMap
}

// Unmarshal decodes b into the receiver.
func (o *RecordInMap) Unmarshal(b []byte) error {
	return avro.Unmarshal(o.Schema(), b, o)
}

// Marshal encodes the receiver.
func (o *RecordInMap) Marshal() ([]byte, error) {
	return avro.Marshal(o.Schema(), o)
}

// RecordInArray is a generated struct.
type RecordInArray struct {
	AString string `avro:"aString"`
}

var schemaRecordInArray = avro.MustParse(`{"name":"a.b.recordInArray","type":"record","fields":[{"name":"aString","type":"string"}]}`)

// Schema returns the schema for RecordInArray.
func (o *RecordInArray) Schema() avro.Schema {
	return schemaRecordInArray
}

// Unmarshal decodes b into the receiver.
func (o *RecordInArray) Unmarshal(b []byte) error {
	return avro.Unmarshal(o.Schema(), b, o)
}

// Marshal encodes the receiver.
func (o *RecordInArray) Marshal() ([]byte, error) {
	return avro.Marshal(o.Schema(), o)
}

// RecordInNullableUnion is a generated struct.
type RecordInNullableUnion struct {
	AString string `avro:"aString"`
}

var schemaRecordInNullableUnion = avro.MustParse(`{"name":"a.b.recordInNullableUnion","type":"record","fields":[{"name":"aString","type":"string"}]}`)

// Schema returns the schema for RecordInNullableUnion.
func (o *RecordInNullableUnion) Schema() avro.Schema {
	return schemaRecordInNullableUnion
}

// Unmarshal decodes b into the receiver.
func (o *RecordInNullableUnion) Unmarshal(b []byte) error {
	return avro.Unmarshal(o.Schema(), b, o)
}

// Marshal encodes the receiver.
func (o *RecordInNullableUnion) Marshal() ([]byte, error) {
	return avro.Marshal(o.Schema(), o)
}

// Record1InNonNullableUnion is a generated struct.
type Record1InNonNullableUnion struct {
	AString string `avro:"aString"`
}

var schemaRecord1InNonNullableUnion = avro.MustParse(`{"name":"a.b.record1InNonNullableUnion","type":"record","fields":[{"name":"aString","type":"string"}]}`)

// Schema returns the schema for Record1InNonNullableUnion.
func (o *Record1InNonNullableUnion) Schema() avro.Schema {
	return schemaRecord1InNonNullableUnion
}

// Unmarshal decodes b into the receiver.
func (o *Record1InNonNullableUnion) Unmarshal(b []byte) error {
	return avro.Unmarshal(o.Schema(), b, o)
}

// Marshal encodes the receiver.
func (o *Record1InNonNullableUnion) Marshal() ([]byte, error) {
	return avro.Marshal(o.Schema(), o)
}

// Record2InNonNullableUnion is a generated struct.
type Record2InNonNullableUnion struct {
	AString string `avro:"aString"`
}

var schemaRecord2InNonNullableUnion = avro.MustParse(`{"name":"a.b.record2InNonNullableUnion","type":"record","fields":[{"name":"aString","type":"string"}]}`)

// Schema returns the schema for Record2InNonNullableUnion.
func (o *Record2InNonNullableUnion) Schema() avro.Schema {
	return schemaRecord2InNonNullableUnion
}

// Unmarshal decodes b into the receiver.
func (o *Record2InNonNullableUnion) Unmarshal(b []byte) error {
	return avro.Unmarshal(o.Schema(), b, o)
}

// Marshal encodes the receiver.
func (o *Record2InNonNullableUnion) Marshal() ([]byte, error) {
	return avro.Marshal(o.Schema(), o)
}

// Record1InNullableUnion is a generated struct.
type Record1InNullableUnion struct {
	AString string `avro:"aString"`
}

var schemaRecord1InNullableUnion = avro.MustParse(`{"name":"a.b.record1InNullableUnion","type":"record","fields":[{"name":"aString","type":"string"}]}`)

// Schema returns the schema for Record1InNullableUnion.
func (o *Record1InNullableUnion) Schema() avro.Schema {
	return schemaRecord1InNullableUnion
}

// Unmarshal decodes b into the receiver.
func (o *Record1InNullableUnion) Unmarshal(b []byte) error {
	return avro.Unmarshal(o.Schema(), b, o)
}

// Marshal encodes the receiver.
func (o *Record1InNullableUnion) Marshal() ([]byte, error) {
	return avro.Marshal(o.Schema(), o)
}

// Record2InNullableUnion is a generated struct.
type Record2InNullableUnion struct {
	AString string `avro:"aString"`
}

var schemaRecord2InNullableUnion = avro.MustParse(`{"name":"a.b.record2InNullableUnion","type":"record","fields":[{"name":"aString","type":"string"}]}`)

// Schema returns the schema for Record2InNullableUnion.
func (o *Record2InNullableUnion) Schema() avro.Schema {
	return schemaRecord2InNullableUnion
}

// Unmarshal decodes b into the receiver.
func (o *Record2InNullableUnion) Unmarshal(b []byte) error {
	return avro.Unmarshal(o.Schema(), b, o)
}

// Marshal encodes the receiver.
func (o *Record2InNullableUnion) Marshal() ([]byte, error) {
	return avro.Marshal(o.Schema(), o)
}

// Test is a generated struct.
type Test struct {
	AString                         string                 `avro:"aString"`
	ABoolean                        bool                   `avro:"aBoolean"`
	AnInt                           int                    `avro:"anInt"`
	AFloat                          float32                `avro:"aFloat"`
	ADouble                         float64                `avro:"aDouble"`
	ALong                           int64                  `avro:"aLong"`
	JustBytes                       []byte                 `avro:"justBytes"`
	PrimitiveNullableArrayUnion     *[]string              `avro:"primitiveNullableArrayUnion"`
	InnerRecord                     InnerRecord            `avro:"innerRecord"`
	AnEnum                          string                 `avro:"anEnum"`
	AFixed                          [7]byte                `avro:"aFixed"`
	ALogicalFixed                   avro.LogicalDuration   `avro:"aLogicalFixed"`
	AnotherLogicalFixed             avro.LogicalDuration   `avro:"anotherLogicalFixed"`
	MapOfStrings                    map[string]string      `avro:"mapOfStrings"`
	MapOfRecords                    map[string]RecordInMap `avro:"mapOfRecords"`
	ADate                           time.Time              `avro:"aDate"`
	ADuration                       time.Duration          `avro:"aDuration"`
	ALongTimeMicros                 time.Duration          `avro:"aLongTimeMicros"`
	ALongTimestampMillis            time.Time              `avro:"aLongTimestampMillis"`
	ALongTimestampMicro             time.Time              `avro:"aLongTimestampMicro"`
	ABytesDecimal                   *big.Rat               `avro:"aBytesDecimal"`
	ARecordArray                    []RecordInArray        `avro:"aRecordArray"`
	NullableRecordUnion             *RecordInNullableUnion `avro:"nullableRecordUnion"`
	NonNullableRecordUnion          any                    `avro:"nonNullableRecordUnion"`
	NullableRecordUnionWith3Options any                    `avro:"nullableRecordUnionWith3Options"`
	Ref                             Record2InNullableUnion `avro:"ref"`
	UUID                            string                 `avro:"uuid"`
}

var schemaTest = avro.MustParse(`{"name":"a.b.test","type":"record","fields":[{"name":"aString","type":"string"},{"name":"aBoolean","type":"boolean"},{"name":"anInt","type":"int"},{"name":"aFloat","type":"float"},{"name":"aDouble","type":"double"},{"name":"aLong","type":"long"},{"name":"justBytes","type":"bytes"},{"name":"primitiveNullableArrayUnion","type":["null",{"type":"array","items":"string"}]},{"name":"innerRecord","type":{"name":"a.c.InnerRecord","type":"record","fields":[{"name":"innerJustBytes","type":"bytes"},{"name":"innerPrimitiveNullableArrayUnion","type":["null",{"type":"array","items":"string"}]}]}},{"name":"anEnum","type":{"name":"a.b.Cards","type":"enum","symbols":["SPADES","HEARTS","DIAMONDS","CLUBS"]}},{"name":"aFixed","type":{"name":"a.b.fixedField","type":"fixed","size":7}},{"name":"aLogicalFixed","type":{"name":"a.b.logicalDuration","type":"fixed","size":12,"logicalType":"duration"}},{"name":"anotherLogicalFixed","type":"a.b.logicalDuration"},{"name":"mapOfStrings","type":{"type":"map","values":"string"}},{"name":"mapOfRecords","type":{"type":"map","values":{"name":"a.b.RecordInMap","type":"record","fields":[{"name":"name","type":"string"}]}}},{"name":"aDate","type":{"type":"int","logicalType":"date"}},{"name":"aDuration","type":{"type":"int","logicalType":"time-millis"}},{"name":"aLongTimeMicros","type":{"type":"long","logicalType":"time-micros"}},{"name":"aLongTimestampMillis","type":{"type":"long","logicalType":"timestamp-millis"}},{"name":"aLongTimestampMicro","type":{"type":"long","logicalType":"timestamp-micros"}},{"name":"aBytesDecimal","type":{"type":"bytes","logicalType":"decimal","precision":4,"scale":2}},{"name":"aRecordArray","type":{"type":"array","items":{"name":"a.b.recordInArray","type":"record","fields":[{"name":"aString","type":"string"}]}}},{"name":"nullableRecordUnion","type":["null",{"name":"a.b.recordInNullableUnion","type":"record","fields":[{"name":"aString","type":"string"}]}]},{"name":"nonNullableRecordUnion","type":[{"name":"a.b.record1InNonNullableUnion","type":"record","fields":[{"name":"aString","type":"string"}]},{"name":"a.b.record2InNonNullableUnion","type":"record","fields":[{"name":"aString","type":"string"}]}]},{"name":"nullableRecordUnionWith3Options","type":["null",{"name":"a.b.record1InNullableUnion","type":"record","fields":[{"name":"aString","type":"string"}]},{"name":"a.b.record2InNullableUnion","type":"record","fields":[{"name":"aString","type":"string"}]}]},{"name":"ref","type":"a.b.record2InNullableUnion"},{"name":"uuid","type":{"type":"string","logicalType":"uuid"}}]}`)

// Schema returns the schema for Test.
func (o *Test) Schema() avro.Schema {
	return schemaTest
}

// Unmarshal decodes b into the receiver.
func (o *Test) Unmarshal(b []byte) error {
	return avro.Unmarshal(o.Schema(), b, o)
}

// Marshal encodes the receiver.
func (o *Test) Marshal() ([]byte, error) {
	return avro.Marshal(o.Schema(), o)
}
