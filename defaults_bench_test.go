package defaults_test

import (
	"testing"

	mdefaults "github.com/mcuadros/go-defaults"
	xdefaults "github.com/xunleii/go-defaults"
)

func BenchmarkLogic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		foo := &ExampleBasic{}
		xdefaults.SetDefaults(foo)
	}
}

func BenchmarkMcuadrosLogic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		foo := &ExampleBasic{}
		mdefaults.SetDefaults(foo)
	}
}
