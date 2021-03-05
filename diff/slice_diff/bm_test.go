package slice_diff

import (
	"strconv"
	"testing"
)

func BenchmarkChanSlice(b *testing.B) {
	cs := NewChanSlice(10, b.N)
	cs.Start()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cs.Append(strconv.Itoa(i))
	}
}

func BenchmarkLockSlice(b *testing.B) {
	cs := NewLockSlice(b.N)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cs.Append(strconv.Itoa(i))
	}
}
