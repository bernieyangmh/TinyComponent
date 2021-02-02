package main

import (
	"testing"
	"time"
)

func BenchmarkMutilChan(b *testing.B) {
	mc := NewMutilChan(time.Millisecond*100, 100)
	mc.Start()

	go func() {
		i := 0
		for {
			mc.Send(i)
			i++
			time.Sleep(time.Millisecond * 1)
		}
	}()

	go func() {
		for _ = range mc.SendC {
			//fmt.Println(fmt.Sprintf("recv item(%d)", i))

		}
	}()
	time.Sleep(time.Second * 3)
	mc.Close()
}
