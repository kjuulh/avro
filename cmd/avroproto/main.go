package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/kjuulh/avro/v2"
	"github.com/kjuulh/avro/v2/protogen"
)

type config struct {
	Out     string
	PkgName string
}

func main() {
	os.Exit(realMain(os.Args, os.Stdout, os.Stderr))
}

func realMain(args []string, stdout, stderr io.Writer) int {
	var cfg config
	flgs := flag.NewFlagSet("avroproto", flag.ExitOnError)
	flgs.SetOutput(stderr)
	flgs.StringVar(&cfg.Out, "o", "", "The output file path to write to instead of stdout.")
	flgs.StringVar(&cfg.PkgName, "p", "", "The package name for which the protobuf file should include")
	flgs.Usage = func() {
		_, _ = fmt.Fprintln(stderr, "Usage: avroproto [options] schemas")
		_, _ = fmt.Fprintln(stderr, "Options:")
		flgs.PrintDefaults()
	}
	if err := flgs.Parse(args[1:]); err != nil {
		return 1
	}

	if err := validateOpts(flgs.NArg(), cfg); err != nil {
		_, _ = fmt.Fprintln(stderr, "Error: "+err.Error())
		return 1
	}

	opts := []protogen.OptsFunc{}
	g := protogen.NewGenerator(cfg.PkgName, opts...)
	for _, file := range flgs.Args() {
		schema, err := avro.ParseFiles(filepath.Clean(file))
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "Error: %v\n", err)
			return 2
		}
		g.Parse(schema)
	}

	var buf bytes.Buffer
	if err := g.Write(&buf); err != nil {
		_, _ = fmt.Fprintf(stderr, "Error: could not generate code: %v\n", err)
		log.Printf("Error: could not generate code: %v\n", err)
		return 3
	}

	writer := stdout
	if cfg.Out != "" {
		file, err := os.Create(cfg.Out)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "Error: could not create output file: %v\n", err)
			return 4
		}
		defer func() { _ = file.Close() }()

		writer = file
	}

	if _, err := writer.Write(buf.Bytes()); err != nil {
		_, _ = fmt.Fprintf(stderr, "Error: could not write code: %v\n", err)
		return 4
	}

	return 0
}

func validateOpts(nargs int, cfg config) error {
	if nargs < 1 {
		return fmt.Errorf("at least one schema is required")
	}

	return nil
}
