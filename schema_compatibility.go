package avro

import (
	"errors"
	"fmt"
	"sync"
)

type recursionError struct{}

func (e recursionError) Error() string {
	return ""
}

type compatKey struct {
	reader [32]byte
	writer [32]byte
}

// SchemaCompatibility determines the compatibility of schemas.
type SchemaCompatibility struct {
	cache sync.Map // map[compatKey]error
}

// NewSchemaCompatibility creates a new schema compatibility instance.
func NewSchemaCompatibility() *SchemaCompatibility {
	return &SchemaCompatibility{}
}

// Compatible determines the compatibility if the reader and writer schemas.
func (c *SchemaCompatibility) Compatible(reader, writer Schema) error {
	return c.compatible(reader, writer)
}

func (c *SchemaCompatibility) compatible(reader, writer Schema) error {
	key := compatKey{reader: reader.Fingerprint(), writer: writer.Fingerprint()}
	if err, ok := c.cache.Load(key); ok {
		if _, ok := err.(recursionError); ok {
			// Break the recursion here.
			return nil
		}

		if err == nil {
			return nil
		}

		return err.(error)
	}

	c.cache.Store(key, recursionError{})
	err := c.match(reader, writer)
	if err != nil {
		// We dont want to pay the cost of fmt.Errorf every time
		err = errors.New(err.Error())
	}
	c.cache.Store(key, err)
	return err
}

func (c *SchemaCompatibility) match(reader, writer Schema) error {
	// If the schema is a reference, get the actual schema
	if reader.Type() == Ref {
		reader = reader.(*RefSchema).Schema()
	}
	if writer.Type() == Ref {
		writer = writer.(*RefSchema).Schema()
	}

	if reader.Type() != writer.Type() {
		if writer.Type() == Union {
			// Reader must be compatible with all types in writer
			for _, schema := range writer.(*UnionSchema).Types() {
				if err := c.compatible(reader, schema); err != nil {
					return err
				}
			}

			return nil
		}

		if reader.Type() == Union {
			// Writer must be compatible with at least one reader schema
			var err error
			for _, schema := range reader.(*UnionSchema).Types() {
				err = c.compatible(schema, writer)
				if err == nil {
					return nil
				}
			}

			return fmt.Errorf("reader union lacking writer schema %s", writer.Type())
		}

		switch writer.Type() {
		case Int:
			if reader.Type() == Long || reader.Type() == Float || reader.Type() == Double {
				return nil
			}

		case Long:
			if reader.Type() == Float || reader.Type() == Double {
				return nil
			}

		case Float:
			if reader.Type() == Double {
				return nil
			}

		case String:
			if reader.Type() == Bytes {
				return nil
			}

		case Bytes:
			if reader.Type() == String {
				return nil
			}
		}

		return fmt.Errorf("reader schema %s not compatible with writer schema %s", reader.Type(), writer.Type())
	}

	switch reader.Type() {
	case Array:
		return c.compatible(reader.(*ArraySchema).Items(), writer.(*ArraySchema).Items())

	case Map:
		return c.compatible(reader.(*MapSchema).Values(), writer.(*MapSchema).Values())

	case Fixed:
		r := reader.(*FixedSchema)
		w := writer.(*FixedSchema)

		if err := c.checkSchemaName(r, w); err != nil {
			return err
		}

		if err := c.checkFixedSize(r, w); err != nil {
			return err
		}

	case Enum:
		r := reader.(*EnumSchema)
		w := writer.(*EnumSchema)

		if err := c.checkSchemaName(r, w); err != nil {
			return err
		}

		if err := c.checkEnumSymbols(r, w); err != nil {
			if r.HasDefault() {
				return nil
			}
			return err
		}

	case Record:
		r := reader.(*RecordSchema)
		w := writer.(*RecordSchema)

		if err := c.checkSchemaName(r, w); err != nil {
			return err
		}

		if err := c.checkRecordFields(r, w); err != nil {
			return err
		}

	case Union:
		for _, schema := range writer.(*UnionSchema).Types() {
			if err := c.compatible(reader, schema); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *SchemaCompatibility) checkSchemaName(reader, writer NamedSchema) error {
	if reader.FullName() != writer.FullName() {
		if c.contains(reader.Aliases(), writer.FullName()) {
			return nil
		}
		return fmt.Errorf("reader schema %s and writer schema %s  names do not match", reader.FullName(), writer.FullName())
	}

	return nil
}

func (c *SchemaCompatibility) checkFixedSize(reader, writer *FixedSchema) error {
	if reader.Size() != writer.Size() {
		return fmt.Errorf("%s reader and writer fixed sizes do not match", reader.FullName())
	}

	return nil
}

func (c *SchemaCompatibility) checkEnumSymbols(reader, writer *EnumSchema) error {
	for _, symbol := range writer.Symbols() {
		if !c.contains(reader.Symbols(), symbol) {
			return fmt.Errorf("reader %s is missing symbol %s", reader.FullName(), symbol)
		}
	}

	return nil
}

func (c *SchemaCompatibility) checkRecordFields(reader, writer *RecordSchema) error {
	for _, field := range reader.Fields() {
		f, ok := c.getField(writer.Fields(), field, func(gfo *getFieldOptions) {
			gfo.fieldAlias = true
		})
		if !ok {
			if field.HasDefault() {
				continue
			}

			return fmt.Errorf("reader field %s is missing in writer schema and has no default", field.Name())
		}

		if err := c.compatible(field.Type(), f.Type()); err != nil {
			return err
		}
	}

	return nil
}

func (c *SchemaCompatibility) contains(a []string, s string) bool {
	for _, str := range a {
		if str == s {
			return true
		}
	}

	return false
}

type getFieldOptions struct {
	fieldAlias bool
	elemAlias  bool
}

func (c *SchemaCompatibility) getField(a []*Field, f *Field, optFns ...func(*getFieldOptions)) (*Field, bool) {
	opt := getFieldOptions{}
	for _, fn := range optFns {
		fn(&opt)
	}
	for _, field := range a {
		if field.Name() == f.Name() {
			return field, true
		}
		if opt.fieldAlias {
			if c.contains(f.Aliases(), field.Name()) {
				return field, true
			}
		}
		if opt.elemAlias {
			if c.contains(field.Aliases(), f.Name()) {
				return field, true
			}
		}
	}

	return nil, false
}

// Resolve returns a composite schema that allows decoding data written by the writer schema,
// and makes necessary adjustments to support the reader schema.
//
// It fails if the writer and reader schemas are not compatible.
func (c *SchemaCompatibility) Resolve(reader, writer Schema) (Schema, error) {
	if err := c.compatible(reader, writer); err != nil {
		return nil, err
	}

	return c.resolve(reader, writer)
}

// resolve requires the reader's schema to be already compatible with the writer's.
func (c *SchemaCompatibility) resolve(reader, writer Schema) (Schema, error) {
	if reader.Type() == Ref {
		reader = reader.(*RefSchema).Schema()
	}
	if writer.Type() == Ref {
		writer = writer.(*RefSchema).Schema()
	}

	if writer.Type() != reader.Type() {
		if reader.Type() == Union {
			for _, schema := range reader.(*UnionSchema).Types() {
				sch, err := c.Resolve(schema, writer)
				if err != nil {
					continue
				}
				return sch, nil
			}

			return nil, fmt.Errorf("reader union lacking writer schema %s", writer.Type())
		}

		if writer.Type() == Union {
			schemas := make([]Schema, 0)
			for _, schema := range writer.(*UnionSchema).Types() {
				sch, err := c.Resolve(reader, schema)
				if err != nil {
					return nil, err
				}
				schemas = append(schemas, sch)
			}
			return NewUnionSchema(schemas)
		}

		if isPromotable(writer.Type()) {
			r := NewPrimitiveSchema(reader.Type(), reader.(*PrimitiveSchema).Logical())
			r.actual = writer.Type()
			return r, nil
		}

		return nil, fmt.Errorf("failed to resolve composite schema for %s and %s", reader.Type(), writer.Type())
	}

	if isNative(writer.Type()) {
		return reader, nil
	}

	if writer.Type() == Enum {
		r := reader.(*EnumSchema)
		w := writer.(*EnumSchema)
		if err := c.checkEnumSymbols(r, w); err != nil {
			if r.HasDefault() {
				enum, _ := NewEnumSchema(r.Name(), r.Namespace(), r.Symbols(),
					WithAliases(r.Aliases()),
					WithDefault(r.Default()),
				)
				enum.actual = w.Symbols()
				return enum, nil
			}

			return nil, err
		}
		return reader, nil
	}

	if writer.Type() == Fixed {
		return reader, nil
	}

	if writer.Type() == Union {
		schemas := make([]Schema, 0)
		for _, schema := range writer.(*UnionSchema).Types() {
			sch, err := c.Resolve(reader, schema)
			if err != nil {
				return nil, err
			}
			schemas = append(schemas, sch)
		}
		return NewUnionSchema(schemas)
	}

	if writer.Type() == Array {
		schema, err := c.resolve(reader.(*ArraySchema).Items(), writer.(*ArraySchema).Items())
		if err != nil {
			return nil, err
		}
		return NewArraySchema(schema), nil
	}

	if writer.Type() == Map {
		schema, err := c.resolve(reader.(*MapSchema).Values(), writer.(*MapSchema).Values())
		if err != nil {
			return nil, err
		}
		return NewMapSchema(schema), nil
	}

	if writer.Type() == Record {
		return c.resolveRecord(reader, writer)
	}

	return nil, fmt.Errorf("failed to resolve composite schema for %s and %s", reader.Type(), writer.Type())
}

func (c *SchemaCompatibility) resolveRecord(reader, writer Schema) (Schema, error) {
	w := writer.(*RecordSchema)
	r := reader.(*RecordSchema)

	fields := make([]*Field, 0)
	seen := make(map[string]struct{})

	for _, wf := range w.Fields() {
		rf, ok := c.getField(r.Fields(), wf, func(gfo *getFieldOptions) {
			gfo.elemAlias = true
		})
		if !ok {
			f, _ := NewField(wf.Name(), wf.Type(), WithAliases(wf.aliases), WithOrder(wf.order))
			// I believe def is read only it can be copied even if it's a like-pointer type;
			// data race should not occur.
			f.def = wf.def
			f.hasDef = wf.hasDef
			f.action = FieldIgnore
			fields = append(fields, f)
			continue
		}

		ft, err := c.resolve(rf.Type(), wf.Type())
		if err != nil {
			return nil, err
		}
		f, _ := NewField(rf.Name(), ft, WithAliases(rf.aliases), WithOrder(rf.order))
		f.def = rf.def
		f.hasDef = rf.hasDef
		fields = append(fields, f)

		seen[rf.Name()] = struct{}{}
	}

	for _, rf := range r.Fields() {
		// check if seen in writer's record
		if _, ok := seen[rf.Name()]; ok {
			continue
		}

		f, _ := NewField(rf.Name(), rf.Type(), WithAliases(rf.aliases), WithOrder(rf.order))
		f.def = rf.def
		f.hasDef = rf.hasDef
		f.action = FieldSetDefault
		fields = append(fields, f)
	}

	return NewRecordSchema(r.Name(), r.Namespace(), fields, WithAliases(r.Aliases()))
}
