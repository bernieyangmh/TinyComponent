package ringBuf

import (
	"testing"
)

func BenchmarkMutilChan(b *testing.B) {
	mc := NewRingBuf(10)
	for i:=0; i<b.N;i++{
		mc.Write("123")
	}
}

