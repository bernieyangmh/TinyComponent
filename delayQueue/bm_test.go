package delayQueue

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func BenchmarkHello(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	dq := NewDelayQueue(200, 0, time.Millisecond*100)
	recordMap := new(sync.Map)

	go func() {
		for i := range dq.readChan {
			d, ok := recordMap.Load(i.(int))
			if ok {
				if d.(DelayItem).T.UnixNano() > time.Now().UnixNano() {
					fmt.Println("wrong")
				}
			}
		}
	}()

	dq.Start()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := DelayItem{Data: i, T: time.Now().Add(time.Second * time.Duration(rand.Intn(3)))}
		dq.Push(d)
		recordMap.Store(i, d)
	}
}

func BenchmarkMutilWriter(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	dq := NewDelayQueue(1000, 100, time.Millisecond*10)

	go func() {
		for _ = range dq.readChan {
		}
	}()

	dq.Start()
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			d := DelayItem{Data: i, T: time.Now().Add(time.Second * time.Duration(rand.Intn(3)))}
			dq.Push(d)
		}
	}()
	go func() {
		for i := 0; i < b.N; i++ {
			d := DelayItem{Data: i, T: time.Now().Add(time.Second * time.Duration(rand.Intn(3)))}
			dq.Push(d)
		}
	}()
}

func BenchmarkMutilReader(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	dq := NewDelayQueue(100, 0, time.Millisecond*10)

	go func() {
		for _ = range dq.readChan {
			//fmt.Println(i)
		}
	}()
	go func() {
		for _ = range dq.readChan {
			//fmt.Println(i)
		}
	}()
	dq.Start()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := DelayItem{Data: i, T: time.Now().Add(time.Second * time.Duration(rand.Intn(3)))}
		dq.Push(d)
	}

}

func BenchmarkMutilWriterReader(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	dq := NewDelayQueue(100, 0, time.Millisecond*10)

	go func() {
		for _ = range dq.readChan {
		}
	}()
	go func() {
		for _ = range dq.readChan {
		}
	}()

	dq.Start()
	b.ResetTimer()
	b.ReportAllocs()
	go func() {
		for i := 0; i < b.N; i++ {
			d := DelayItem{Data: i, T: time.Now().Add(time.Second * time.Duration(rand.Intn(3)))}
			dq.Push(d)
		}
	}()

	for i := 0; i < b.N; i++ {
		d := DelayItem{Data: i, T: time.Now().Add(time.Second * time.Duration(rand.Intn(3)))}
		dq.Push(d)
	}
}
