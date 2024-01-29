// Package gen allows generating Go structs from avro schemas.
package protogen

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io"
	"strings"
	"text/template"

	"github.com/ettle/strcase"
	"github.com/kjuulh/avro/v2"
)

// Config configures the code generation.
type Config struct {
	PackageName string
}

const outputTemplate = `syntax = "proto3"

// Code generated by avro/gen. DO NOT EDIT.

package {{ .PackageName }}

{{- range .Imports }}
import "{{ . }}";
{{- end }}
{{- range .ThirdPartyImports }}
import "{{ . }}";
{{- end }}

{{- range .Typedefs }}
// {{ .Name }} is a generated struct.
message {{ .Name }} {
	{{- range $index, $field := .Fields }}
		{{ $field.Type }} {{ $field.Name }} = {{ add $index 1 }};
	{{- end }}
}
{{ end }}`

var primitiveMappings = map[avro.Type]string{
	"string":  "string",
	"bytes":   "bytes",
	"int":     "int32",
	"long":    "int64",
	"float":   "float",
	"double":  "double",
	"boolean": "bool",
}

// Struct generates Go structs based on the schema and writes them to w.
func Struct(s string, w io.Writer, cfg Config) error {
	schema, err := avro.Parse(s)
	if err != nil {
		return err
	}
	return StructFromSchema(schema, w, cfg)
}

// StructFromSchema generates Go structs based on the schema and writes them to w.
func StructFromSchema(schema avro.Schema, w io.Writer, cfg Config) error {
	rec, ok := schema.(*avro.RecordSchema)
	if !ok {
		return errors.New("can only generate Go code from Record Schemas")
	}

	opts := []OptsFunc{}
	g := NewGenerator(cfg.PackageName, opts...)
	g.Parse(rec)

	buf := &bytes.Buffer{}
	if err := g.Write(buf); err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("could not format code: %w", err)
	}

	_, err = w.Write(formatted)
	return err
}

// OptsFunc is a function that configures a generator.
type OptsFunc func(*Generator)

// Generator generates Go structs from schemas.
type Generator struct {
	imports           []string
	thirdPartyImports []string
	typedefs          []typedef

	pkg       string
	nameCaser *strcase.Caser
}

// NewGenerator returns a generator.
func NewGenerator(pkgName string, opts ...OptsFunc) *Generator {
	g := &Generator{
		pkg:       pkgName,
		nameCaser: strcase.NewCaser(false, map[string]bool{}, nil),
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// Reset reset the generator.
func (g *Generator) Reset() {
	g.imports = g.imports[:0]
	g.thirdPartyImports = g.thirdPartyImports[:0]
	g.typedefs = g.typedefs[:0]
}

// Parse parses an avro schema into Go types.
func (g *Generator) Parse(schema avro.Schema) {
	_ = g.generate(schema)
}

func (g *Generator) generate(schema avro.Schema) string {
	switch s := schema.(type) {
	case *avro.RefSchema:
		return g.resolveRefSchema(s)
	case *avro.RecordSchema:
		return g.resolveRecordSchema(s)
	case *avro.PrimitiveSchema:
		typ := primitiveMappings[s.Type()]
		if ls := s.Logical(); ls != nil {
			typ = g.resolveLogicalSchema(ls.Type())
		}
		return typ
	case *avro.ArraySchema:
		return "repeated " + g.generate(s.Items())
	case *avro.EnumSchema:
		// TODO: maybe do enums as well
		return "string"
	case *avro.FixedSchema:
		typ := "bytes"
		if ls := s.Logical(); ls != nil {
			typ = g.resolveLogicalSchema(ls.Type())
		}
		return typ
	case *avro.MapSchema:
		return "map<string, " + g.generate(s.Values()) + ">"
	case *avro.UnionSchema:
		return g.resolveUnionTypes(s)
	default:
		return ""
	}
}

func (g *Generator) resolveTypeName(s avro.NamedSchema) string {
	return g.nameCaser.ToPascal(s.Name())
}

func (g *Generator) resolveRecordSchema(schema *avro.RecordSchema) string {
	fields := make([]field, len(schema.Fields()))
	for i, f := range schema.Fields() {
		typ := g.generate(f.Type())
		tag := f.Name()
		fields[i] = g.newField(g.nameCaser.ToPascal(f.Name()), typ, tag)
	}

	typeName := g.resolveTypeName(schema)
	if !g.hasTypeDef(typeName) {
		g.typedefs = append(g.typedefs, newType(typeName, fields, schema.String()))
	}
	return typeName
}

func (g *Generator) hasTypeDef(name string) bool {
	for _, def := range g.typedefs {
		if def.Name != name {
			continue
		}
		return true
	}
	return false
}

func (g *Generator) resolveRefSchema(s *avro.RefSchema) string {
	if sx, ok := s.Schema().(*avro.RecordSchema); ok {
		return g.resolveTypeName(sx)
	}
	return g.generate(s.Schema())
}

func (g *Generator) resolveUnionTypes(s *avro.UnionSchema) string {
	types := make([]string, 0, len(s.Types()))
	for _, elem := range s.Types() {
		if _, ok := elem.(*avro.NullSchema); ok {
			continue
		}
		types = append(types, g.generate(elem))
	}
	if s.Nullable() {
		// protobuf version 3 is optional by default
		return types[0]
	}
	return "any"
}

func (g *Generator) resolveLogicalSchema(logicalType avro.LogicalType) string {
	var typ string
	switch logicalType {
	case "date", "timestamp-millis", "timestamp-micros":
		typ = "google.protobuf.Timestamp"
	case "time-millis", "time-micros":
		typ = "google/protobuf/duration.proto"
	case "decimal":
		typ = "*big.Rat"
	case "duration":
		typ = "google/protobuf/duration.proto"
	case "uuid":
		typ = "string"
	}
	if strings.Contains(typ, "Timestamp") {
		g.addImport("google/protobuf/timestamp.proto")
	}
	if strings.Contains(typ, "Duration") {
		g.addImport("google/protobuf/timestamp.proto")
	}
	if strings.Contains(typ, "big") {
		g.addImport("math/big")
	}
	return typ
}

func (g *Generator) newField(name, typ, tag string) field {
	tagLine := fmt.Sprintf(`avro:"%s"`, tag)
	return field{
		Name: name,
		Type: typ,
		Tag:  fmt.Sprintf("`%s`", tagLine),
	}
}

func formatTag(tag string) string {
	return tag
}

func (g *Generator) addImport(pkg string) {
	for _, p := range g.imports {
		if p == pkg {
			return
		}
	}
	g.imports = append(g.imports, pkg)
}

func (g *Generator) addThirdPartyImport(pkg string) {
	for _, p := range g.thirdPartyImports {
		if p == pkg {
			return
		}
	}
	g.thirdPartyImports = append(g.thirdPartyImports, pkg)
}

// Write writes Go code from the parsed schemas.
func (g *Generator) Write(w io.Writer) error {
	parsed, err := template.New("out").Funcs(template.FuncMap{"add": func(numbers ...int) int {
		sum := 0
		for _, num := range numbers {
			sum += num
		}
		return sum
	}}).Parse(outputTemplate)
	if err != nil {
		return err
	}

	data := struct {
		PackageName string

		Imports           []string
		ThirdPartyImports []string
		Typedefs          []typedef
	}{
		PackageName:       g.pkg,
		Imports:           g.imports,
		ThirdPartyImports: g.thirdPartyImports,
		Typedefs:          g.typedefs,
	}
	return parsed.Execute(w, data)
}

type typedef struct {
	Name   string
	Fields []field
	Schema string
}

func newType(name string, fields []field, schema string) typedef {
	return typedef{
		Name:   name,
		Fields: fields,
		Schema: schema,
	}
}

type field struct {
	Name string
	Type string
	Tag  string
}
