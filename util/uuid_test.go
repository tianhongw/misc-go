package util

import (
	"testing"

	"github.com/tianhongw/misc-go/util/assert"
)

func TestRandomString(t *testing.T) {
	for i := 0; i < 100; i++ {
		rs := RandomString(i)
		assert.Equal(t, i, len(rs))
	}
}

func benchmarkRandomString(n int, b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RandomString(n)
	}
}

func BenchmarkRandomString10(b *testing.B) {
	benchmarkRandomString(10, b)
}

func BenchmarkRandomString100(b *testing.B) {
	benchmarkRandomString(100, b)
}

func BenchmarkRandomString1000(b *testing.B) {
	benchmarkRandomString(1000, b)
}

func BenchmarkRandomStringUnsafe10(b *testing.B) {
	benchmarkRandomStringUnsafe(10, b)
}

func BenchmarkRandomStringUnsafe100(b *testing.B) {
	benchmarkRandomStringUnsafe(100, b)
}

func BenchmarkRandomStringUnsafe1000(b *testing.B) {
	benchmarkRandomStringUnsafe(1000, b)
}

func benchmarkRandomStringUnsafe(n int, b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RandStringBytesMaskImprSrcUnsafe(n)
	}
}
