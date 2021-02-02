package barrierChan

// 拦截多个goroutine直到到达阈值 或到一定间隔

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	stateClose = uint32(1)
)

type Barrier struct {
	chCount chan struct{}
	chSync  chan struct{} // 增加一个chan
	count   uint32
	Ticker  *time.Ticker

	state uint32
}

func NewBarrier(n uint32, duration time.Duration) *Barrier {
	b := &Barrier{count: n, chCount: make(chan struct{}), chSync: make(chan struct{})}
	b.Ticker = time.NewTicker(duration)
	go b.Sync()
	return b
}

func (b *Barrier) Wait(i int) {
	fmt.Println(fmt.Sprintf("wait (%d)", i))

	if atomic.LoadUint32(&b.state) == stateClose {
		return
	}

	b.chCount <- struct{}{}
	<-b.chSync // 再次阻塞
}

// close之后，不再阻塞
func (b *Barrier) Close() {
	b.state |= stateClose
	b.Ticker.Stop()
	close(b.chSync)
	close(b.chCount)

}

func (b *Barrier) Sync() {
	count := uint32(1)
	for {
		if atomic.LoadUint32(&b.state) == stateClose {
			return
		}
		select {
		case <-b.chCount:
			count++
			if count >= b.count {
				close(b.chSync) // close这个chan所有阻塞协程都会被激活
				b.chSync = make(chan struct{})
				atomic.StoreUint32(&count, 0)
			}
		case tt := <-b.Ticker.C:
			fmt.Println(tt.String())
			close(b.chSync) // close这个chan所有阻塞协程都会被激活
			b.chSync = make(chan struct{})
			atomic.StoreUint32(&count, 0)
		default:
		}
	}
}
