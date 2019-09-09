package defaults_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/xunleii/go-defaults"
)

type Parent struct {
	Children []Child
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
	Float32    float32 `default:"3.2"`
	Float64    float64 `default:"6.4"`
	Struct     struct {
		Bool    bool `default:"true"`
		Integer int  `default:"33"`
	}
	Duration time.Duration `default:"1s"`
	Children []Child
}

type ExampleNested struct {
	Struct ExampleBasic
}

func TestSetDefaultsBasic(t *testing.T) {
	foo := &ExampleBasic{}
	defaults.SetDefaults(foo)

	assertTypes(t, foo)
}

func TestSetDefaultsNested(t *testing.T) {
	foo := &ExampleNested{}
	defaults.SetDefaults(foo)

	assertTypes(t, &foo.Struct)
}

func TestSetDefaultsWithValues(t *testing.T) {
	foo := &ExampleBasic{
		Integer:  55,
		UInteger: 22,
		Float32:  9.9,
		String:   "bar",
		Children: []Child{{Name: "alice"}, {Name: "bob", Age: 2}},
	}

	defaults.SetDefaults(foo)

	assert.Equal(t, 55, foo.Integer)
	assert.Equal(t, uint(22), foo.UInteger)
	assert.Equal(t, float32(9.9), foo.Float32)
	assert.Equal(t, "bar", foo.String)
	assert.Equal(t, 10, foo.Children[0].Age)
	assert.Equal(t, 2, foo.Children[1].Age)
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
	assert.Equal(t, float32(3.2), foo.Float32)
	assert.Equal(t, 6.4, foo.Float64)
	assert.Equal(t, true, foo.Struct.Bool)     // not work
	assert.Equal(t, time.Second, foo.Duration) // not work
	assert.Nil(t, foo.Children)
}

func BenchmarkLogic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		foo := &ExampleBasic{}
		defaults.SetDefaults(foo)
	}
}
