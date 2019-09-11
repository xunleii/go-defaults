package defaults_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/xunleii/go-defaults"
)

type Parent struct {
	Children []Child `default:"[{\"name\": \"alice\", \"age\": 10},{\"name\": \"bob\", \"age\": 2}]"`
}

type Child struct {
	Name string
	Age  int `default:"10"`
}

type ExampleBasic struct {
	Bool       bool    `default:"true"`
	Integer    int     `default:"33"`
	Integer8   int8    `default:"8"`
	Integer16  int16   `default:"16"`
	Integer32  int32   `default:"32"`
	Integer64  int64   `default:"64"`
	UInteger   uint    `default:"11"`
	UInteger8  uint8   `default:"18"`
	UInteger16 uint16  `default:"116"`
	UInteger32 uint32  `default:"132"`
	UInteger64 uint64  `default:"164"`
	String     string  `default:"foo"`
	Bytes      []byte  `default:"bar"`
	Float32    float32 `default:"3.2"`
	Float64    float64 `default:"6.4"`
	Struct     struct {
		Bool    bool `default:"true"`
		Integer int  `default:"33"`
	}
	Children []Child
	Duration time.Duration `default:"1s"`
}

type ExampleNested struct {
	Struct ExampleBasic
}

func TestSetDefaultsBasic(t *testing.T) {
	basic := &ExampleBasic{}
	defaults.SetDefaults(basic)

	assertTypes(t, basic)
}

func TestSetDefaultsNested(t *testing.T) {
	nested := &ExampleNested{}
	defaults.SetDefaults(nested)

	assertTypes(t, &nested.Struct)
}

func TestSetDefaultsWithValues(t *testing.T) {
	basic := &ExampleBasic{
		Integer:  55,
		UInteger: 22,
		Float32:  9.9,
		String:   "bar",
		Bytes:    []byte("foo"),
		Children: []Child{{Name: "alice"}, {Name: "bob", Age: 2}},
	}

	defaults.SetDefaults(basic)

	assert.Equal(t, 55, basic.Integer)
	assert.Equal(t, uint(22), basic.UInteger)
	assert.Equal(t, float32(9.9), basic.Float32)
	assert.Equal(t, "bar", basic.String)
	assert.Equal(t, "foo", string(basic.Bytes))
	assert.Equal(t, 10, basic.Children[0].Age)
	assert.Equal(t, 2, basic.Children[1].Age)
}

func TestRegisterCustomDefault(t *testing.T) {
	defaults.RegisterCustomDefault(reflect.TypeOf([]Child{}), func(field *defaults.FieldData) {
		var v []Child
		_ = json.Unmarshal([]byte(field.TagValue), &v)
		field.Value.Set(reflect.ValueOf(v))
	})

	parent := &Parent{}
	defaults.SetDefaults(parent)

	assert.Len(t, parent.Children, 2)
	assert.Equal(t, "alice", parent.Children[0].Name)
	assert.Equal(t, 10, parent.Children[0].Age)
	assert.Equal(t, "bob", parent.Children[1].Name)
	assert.Equal(t, 2, parent.Children[1].Age)
}

func assertTypes(t *testing.T, foo *ExampleBasic) {
	assert.Equal(t, true, foo.Bool)
	assert.Equal(t, 33, foo.Integer)
	assert.Equal(t, int8(8), foo.Integer8)
	assert.Equal(t, int16(16), foo.Integer16)
	assert.Equal(t, int32(32), foo.Integer32)
	assert.Equal(t, int64(64), foo.Integer64)
	assert.Equal(t, uint(11), foo.UInteger)
	assert.Equal(t, uint8(18), foo.UInteger8)
	assert.Equal(t, uint16(116), foo.UInteger16)
	assert.Equal(t, uint32(132), foo.UInteger32)
	assert.Equal(t, uint64(164), foo.UInteger64)
	assert.Equal(t, "foo", foo.String)
	assert.Equal(t, "bar", string(foo.Bytes))
	assert.Equal(t, float32(3.2), foo.Float32)
	assert.Equal(t, 6.4, foo.Float64)
	assert.Equal(t, true, foo.Struct.Bool)
	assert.Equal(t, time.Second, foo.Duration)
	assert.Nil(t, foo.Children)
}
