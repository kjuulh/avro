package avro

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestName_NameAndNamespace(t *testing.T) {
	n, err := newName("bar", "foo", nil)
	require.NoError(t, err)

	assert.Equal(t, "bar", n.Name())
	assert.Equal(t, "foo", n.Namespace())
	assert.Equal(t, "foo.bar", n.FullName())
}

func TestName_QualifiedName(t *testing.T) {
	n, err := newName("foo.bar", "test", nil)
	require.NoError(t, err)

	assert.Equal(t, "bar", n.Name())
	assert.Equal(t, "foo", n.Namespace())
	assert.Equal(t, "foo.bar", n.FullName())
}

func TestName_NameAndNamespaceAndAlias(t *testing.T) {
	n, err := newName("bar", "foo", []string{"baz", "test.bat"})
	require.NoError(t, err)

	assert.Equal(t, "bar", n.Name())
	assert.Equal(t, "foo", n.Namespace())
	assert.Equal(t, "foo.bar", n.FullName())
	assert.Equal(t, []string{"foo.baz", "test.bat"}, n.Aliases())
}

func TestName_EmpryName(t *testing.T) {
	_, err := newName("", "foo", nil)

	assert.Error(t, err)
}

func TestName_InvalidNameFirstChar(t *testing.T) {
	_, err := newName("+bar", "foo", nil)

	assert.Error(t, err)
}

func TestName_InvalidNameOtherChar(t *testing.T) {
	_, err := newName("bar+", "foo", nil)

	assert.Error(t, err)
}

func TestName_InvalidNamespaceFirstChar(t *testing.T) {
	_, err := newName("bar", "+foo", nil)

	assert.Error(t, err)
}

func TestName_InvalidNamespaceOtherChar(t *testing.T) {
	_, err := newName("bar", "foo+", nil)

	assert.Error(t, err)
}

func TestName_InvalidAliasFirstChar(t *testing.T) {
	_, err := newName("bar", "foo", []string{"+bar"})

	assert.Error(t, err)
}

func TestName_InvalidAliasOtherChar(t *testing.T) {
	_, err := newName("bar", "foo", []string{"bar+"})

	assert.Error(t, err)
}

func TestName_InvalidAliasFQNFirstChar(t *testing.T) {
	_, err := newName("bar", "foo", []string{"test.+bar"})

	assert.Error(t, err)
}

func TestName_InvalidAliasFQNOtherChar(t *testing.T) {
	_, err := newName("bar", "foo", []string{"test.bar+"})

	assert.Error(t, err)
}

func TestProperties_PropGetsFromEmptySet(t *testing.T) {
	p := properties{}

	assert.Nil(t, p.Prop("test"))
}

func TestIsValidDefault(t *testing.T) {
	tests := []struct {
		name     string
		schemaFn func() Schema
		def      any
		want     any
		wantOk   bool
	}{

		{
			name: "Null",
			schemaFn: func() Schema {
				return &NullSchema{}
			},
			def:    nil,
			want:   nullDefault,
			wantOk: true,
		},
		{
			name: "Null Invalid Type",
			schemaFn: func() Schema {
				return &NullSchema{}
			},
			def:    "test",
			wantOk: false,
		},
		{
			name: "String",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(String, nil)
			},
			def:    "test",
			want:   "test",
			wantOk: true,
		},
		{
			name: "String Invalid Type",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(String, nil)
			},
			def:    1,
			wantOk: false,
		},
		{
			name: "Bytes",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Bytes, nil)
			},
			def:    "test",
			want:   []byte("test"),
			wantOk: true,
		},
		{
			name: "Bytes Invalid Type",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Bytes, nil)
			},
			def:    1,
			wantOk: false,
		},
		{
			name: "Enum",
			schemaFn: func() Schema {
				s, _ := NewEnumSchema("foo", "", []string{"BAR"})
				return s
			},
			def:    "BAR",
			want:   "BAR",
			wantOk: true,
		},
		{
			name: "Enum Invalid Default",
			schemaFn: func() Schema {
				s, _ := NewEnumSchema("foo", "", []string{"BAR"})
				return s
			},
			def:    "BUP",
			wantOk: false,
		},
		{
			name: "Enum Empty string",
			schemaFn: func() Schema {
				s, _ := NewEnumSchema("foo", "", []string{"BAR"})
				return s
			},
			def:    "",
			wantOk: false,
		},
		{
			name: "Enum Invalid Type",
			schemaFn: func() Schema {
				s, _ := NewEnumSchema("foo", "", []string{"BAR"})
				return s
			},
			def:    1,
			wantOk: false,
		},
		{
			name: "Fixed",
			schemaFn: func() Schema {
				s, _ := NewFixedSchema("foo", "", 4, nil)
				return s
			},
			def:    "test",
			want:   [4]byte{'t', 'e', 's', 't'},
			wantOk: true,
		},
		{
			name: "Fixed Invalid Type",
			schemaFn: func() Schema {
				s, _ := NewFixedSchema("foo", "", 1, nil)
				return s
			},
			def:    1,
			wantOk: false,
		},
		{
			name: "Boolean",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Boolean, nil)
			},
			def:    true,
			want:   true,
			wantOk: true,
		},
		{
			name: "Boolean Invalid Type",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Boolean, nil)
			},
			def:    1,
			wantOk: false,
		},
		{
			name: "Int",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Int, nil)
			},
			def:    1,
			want:   1,
			wantOk: true,
		},
		{
			name: "Int Int8",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Int, nil)
			},
			def:    int8(1),
			want:   1,
			wantOk: true,
		},
		{
			name: "Int Int16",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Int, nil)
			},
			def:    int16(1),
			want:   1,
			wantOk: true,
		},
		{
			name: "Int Int32",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Int, nil)
			},
			def:    int32(1),
			want:   1,
			wantOk: true,
		},
		{
			name: "Int Float64",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Int, nil)
			},
			def:    float64(1),
			want:   1,
			wantOk: true,
		},
		{
			name: "Int Invalid Type",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Int, nil)
			},
			def:    "test",
			wantOk: false,
		},
		{
			name: "Long",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Long, nil)
			},
			def:    int64(1),
			want:   int64(1),
			wantOk: true,
		},
		{
			name: "Long Float64",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Long, nil)
			},
			def:    float64(1),
			want:   int64(1),
			wantOk: true,
		},
		{
			name: "Long Invalid Type",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Long, nil)
			},
			def:    "test",
			wantOk: false,
		},
		{
			name: "Float",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Float, nil)
			},
			def:    float32(1),
			want:   float32(1),
			wantOk: true,
		},
		{
			name: "Float Float64",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Float, nil)
			},
			def:    float64(1),
			want:   float32(1),
			wantOk: true,
		},
		{
			name: "Float Invalid Type",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Float, nil)
			},
			def:    "test",
			wantOk: false,
		},
		{
			name: "Double",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Double, nil)
			},
			def:    float64(1),
			want:   float64(1),
			wantOk: true,
		},
		{
			name: "Double Invalid Type",
			schemaFn: func() Schema {
				return NewPrimitiveSchema(Double, nil)
			},
			def:    "test",
			wantOk: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, ok := isValidDefault(test.schemaFn(), test.def)

			assert.Equal(t, test.wantOk, ok)
			if ok {
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestRecursionError_Error(t *testing.T) {
	err := recursionError{}

	assert.Equal(t, "", err.Error())
}

func TestSchema_FingerprintUsingCaches(t *testing.T) {
	schema := NewPrimitiveSchema(String, nil)

	want, _ := schema.FingerprintUsing(CRC64Avro)

	got, _ := schema.FingerprintUsing(CRC64Avro)

	value, ok := schema.cache.Load(CRC64Avro)
	require.True(t, ok)
	assert.Equal(t, want, value)
	assert.Equal(t, want, got)
}

func TestSchema_IsPromotable(t *testing.T) {
	tests := []struct {
		typ    Type
		wantOk bool
	}{
		{
			typ:    Int,
			wantOk: true,
		},
		{
			typ:    Long,
			wantOk: true,
		},
		{
			typ:    Float,
			wantOk: true,
		},
		{
			typ:    String,
			wantOk: true,
		},
		{
			typ:    Bytes,
			wantOk: true,
		},
		{
			typ:    Double,
			wantOk: false,
		},
		{
			typ:    Boolean,
			wantOk: false,
		},
		{
			typ:    Null,
			wantOk: false,
		},
	}

	for i, test := range tests {
		test := test
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ok := isPromotable(test.typ)
			assert.Equal(t, test.wantOk, ok)
		})
	}
}

func TestSchema_IsNative(t *testing.T) {
	tests := []struct {
		typ    Type
		wantOk bool
	}{
		{
			typ:    Null,
			wantOk: true,
		},
		{
			typ:    Boolean,
			wantOk: true,
		},
		{
			typ:    Int,
			wantOk: true,
		},
		{
			typ:    Long,
			wantOk: true,
		},

		{
			typ:    Float,
			wantOk: true,
		},
		{
			typ:    Double,
			wantOk: true,
		},

		{
			typ:    Bytes,
			wantOk: true,
		},
		{
			typ:    String,
			wantOk: true,
		},
		{
			typ:    Record,
			wantOk: false,
		},
		{
			typ:    Array,
			wantOk: false,
		},
		{
			typ:    Map,
			wantOk: false,
		},
		{
			typ:    Fixed,
			wantOk: false,
		},
		{
			typ:    Enum,
			wantOk: false,
		},
		{
			typ:    Union,
			wantOk: false,
		},
	}

	for i, test := range tests {
		test := test
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ok := isNative(test.typ)
			assert.Equal(t, test.wantOk, ok)
		})
	}
}

func TestSchema_FieldEncodeDefault(t *testing.T) {
	schema := MustParse(`{
		"type": "record",
		"name": "test",
		"fields" : [
			{"name": "a", "type": "string", "default": "bar"},
			{"name": "b", "type": "boolean"}
		]
	}`).(*RecordSchema)

	fooEncoder := func(a any) ([]byte, error) {
		return []byte("foo"), nil
	}
	barEncoder := func(a any) ([]byte, error) {
		return []byte("bar"), nil
	}

	assert.Equal(t, nil, schema.fields[0].encodedDef.Load())

	_, err := schema.fields[0].encodeDefault(nil)
	assert.Error(t, err)

	_, err = schema.fields[1].encodeDefault(fooEncoder)
	assert.Error(t, err)

	def, err := schema.fields[0].encodeDefault(fooEncoder)
	assert.NoError(t, err)
	assert.Equal(t, []byte("foo"), def)

	def, err = schema.fields[0].encodeDefault(barEncoder)
	assert.NoError(t, err)
	assert.Equal(t, []byte("foo"), def)
}

func TestSchema_CacheFingerprint(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		cacheFingerprint := cacheFingerprinter{}
		assert.Panics(t, func() {
			cacheFingerprint.fingerprint([]any{func() {}})
		})
	})

	t.Run("promoted", func(t *testing.T) {
		schema := NewPrimitiveSchema(Long, nil)
		assert.Equal(t, schema.Fingerprint(), schema.CacheFingerprint())

		schema = NewPrimitiveSchema(Long, nil)
		schema.actual = Int
		assert.NotEqual(t, schema.Fingerprint(), schema.CacheFingerprint())
	})

	t.Run("enum", func(t *testing.T) {
		schema1 := MustParse(`{
			"type": "enum",
			"name": "test.enum",
			"symbols": ["foo"]
		}`).(*EnumSchema)

		schema2 := MustParse(`{
			"type": "enum",
			"name": "test.enum",
			"symbols": ["foo"],
			"default": "foo"
			}`).(*EnumSchema)
		schema2.actual = []string{"boo"}

		assert.Equal(t, schema1.Fingerprint(), schema1.CacheFingerprint())
		assert.NotEqual(t, schema1.CacheFingerprint(), schema2.CacheFingerprint())
	})

	t.Run("record", func(t *testing.T) {
		schema1 := MustParse(`{
			"type": "record",
			"name": "test",
			"fields" : [
				{"name": "a", "type": "string"},
				{"name": "b", "type": "boolean"}
			]
		}`).(*RecordSchema)

		schema2 := MustParse(`{
			"type": "record",
			"name": "test2",
			"fields" : [
				{"name": "a", "type": "string", "default": "bar"},
				{"name": "b", "type": "boolean", "default": false}
			]
		}`).(*RecordSchema)

		assert.Equal(t, schema1.Fingerprint(), schema1.CacheFingerprint())
		assert.NotEqual(t, schema1.CacheFingerprint(), schema2.CacheFingerprint())
	})
}

func TestEnumSchema_GetSymbol(t *testing.T) {
	tests := []struct {
		schemaFn func() *EnumSchema
		idx      int
		want     any
		wantOk   bool
	}{
		{
			schemaFn: func() *EnumSchema {
				enum, _ := NewEnumSchema("foo", "", []string{"BAR"})
				return enum
			},
			idx:    0,
			wantOk: true,
			want:   "BAR",
		},
		{
			schemaFn: func() *EnumSchema {
				enum, _ := NewEnumSchema("foo", "", []string{"BAR"})
				return enum
			},
			idx:    1,
			wantOk: false,
		},
		{
			schemaFn: func() *EnumSchema {
				enum, _ := NewEnumSchema("foo", "", []string{"FOO"}, WithDefault("FOO"))
				return enum
			},
			idx:    1,
			wantOk: false,
		},
		{
			schemaFn: func() *EnumSchema {
				enum, _ := NewEnumSchema("foo", "", []string{"FOO"})
				enum.actual = []string{"FOO", "BAR"}
				return enum
			},
			idx:    1,
			wantOk: false,
		},
		{
			schemaFn: func() *EnumSchema {
				enum, _ := NewEnumSchema("foo", "", []string{"FOO"}, WithDefault("FOO"))
				enum.actual = []string{"FOO", "BAR"}
				return enum
			},
			idx:    1,
			wantOk: true,
			want:   "FOO",
		},
		{
			schemaFn: func() *EnumSchema {
				enum, _ := NewEnumSchema("foo", "", []string{"FOO", "BAR"})
				enum.actual = []string{"FOO"}
				return enum
			},
			idx:    0,
			wantOk: true,
			want:   "FOO",
		},
	}

	for i, test := range tests {
		test := test
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, ok := test.schemaFn().Symbol(test.idx)
			assert.Equal(t, test.wantOk, ok)
			if ok {
				assert.Equal(t, test.want, got)
			}
		})
	}
}
