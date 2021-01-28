package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type mutilChan struct {
	mutex *sync.Mutex

	curIndex uint64
	recvC    chan interface{}
	SendC    chan interface{}
	buf      []interface{}
	bufsize  uint64
	ticker   *time.Ticker
	close    uint32

	*sync.Pool
}

func NewMutilChan(ti time.Duration, bufsize uint64) *mutilChan {
	mc := new(mutilChan)
	mc.recvC = make(chan interface{})
	mc.SendC = make(chan interface{})
	mc.bufsize = bufsize
	mc.buf = make([]interface{}, bufsize, bufsize)
	mc.buf = make([]interface{}, bufsize, bufsize)
	mc.ticker = time.NewTicker(ti)
	mc.mutex = new(sync.Mutex)
	mc.Pool = &sync.Pool{
		New: func() interface{} {
			return make([]interface{}, bufsize)
		},
	}
	return mc
}

func (mc *mutilChan) Start() {
	go func() {
		for {

			select {
			case <-mc.ticker.C:
				mc.mutilSend()
				atomic.StoreUint64(&mc.curIndex, 0)
			default:

			}
			if mc.curIndex == mc.bufsize {
				mc.mutilSend()
				atomic.StoreUint64(&mc.curIndex, 0)
			}
			mc.buf[mc.curIndex] = <-mc.recvC
			mc.curIndex++

		}
	}()

}

func (mc *mutilChan) Close() {
	mc.close = 1

	mc.ticker.Stop()

	mc.mutilSend()
	for mc.close == 0 {
		close(mc.recvC)
		close(mc.SendC)
	}
}

func (mc *mutilChan) Send(item interface{}) error {
	if mc.close == 1 {
		return fmt.Errorf("chan alreay close")
	}
	mc.recvC <- item
	return nil
}

func (mc *mutilChan) mutilSend() {
	dirtyBuf := mc.buf
	mc.buf = mc.Pool.Get().([]interface{})

	go func(dirtyBuf []interface{}) {
		for i := 0; i < len(dirtyBuf); i++ {
			if dirtyBuf[i] != nil {
				mc.SendC <- dirtyBuf[i]
			}

		}

		dirtyBuf = nil
		atomic.CompareAndSwapUint32(&mc.close, 1, 0)
	}(dirtyBuf)
	return

}
