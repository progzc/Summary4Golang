package inline_function

import "testing"

var Result1 int

func BenchmarkMax1(b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		r = max1(-1, i)
	}
	Result1 = r
}

var Result2 int

func BenchmarkMax2(b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		r = max2(-1, i)
	}
	Result2 = r
}
