package str_rune

import (
	"testing"
)

func BenchmarkStrFormat(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		StrFormat("123", "abc", "你好")
	}
}

func BenchmarkStrPlus(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		StrPlus("123", "abc", "你好")
	}
}

func BenchmarkStrBuilder(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		StrBuilder("123", "abc", "你好")
	}
}
