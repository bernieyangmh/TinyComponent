package delayQueue

import (
	"container/heap"
	"fmt"
	"sync/atomic"
	"time"
)

type DelayQueue struct {
	ticker    *time.Ticker
	Queue     delayQueue
	state     int32
	readChan  chan interface{}
	closeChan chan int8
}

func (dq *DelayQueue) Start() {
	go func() {
		for {
			select {
			case <-dq.ticker.C:
				now := time.Now()
				top := dq.Peek()
				if now.Sub(top.T) > 0 {
					continue
				}
				for {
					di, ok := dq.Pop()
					if !ok {
						break
					}
					if now.Sub(di.T) < 0 {
						dq.readChan <- di.Data
					} else {
						break
					}
				}
			case <-dq.closeChan:
				for {
					di, ok := dq.Pop()
					if !ok {
						break
					}
					dq.readChan <- di.Data
				}
				close(dq.closeChan)
				break
				// already clear
			}
		}
	}()
}

func (dq *DelayQueue) Close() {
	ok := atomic.CompareAndSwapInt32(&dq.state, 0, 1)
	if !ok {
		// already closed
		return
	}
	dq.closeChan <- 1
	for {
		_, ok := <-dq.closeChan
		if !ok {
			break
		}
	}

	return
}

func (dq *DelayQueue) Push(x DelayItem) error {
	if atomic.LoadInt32(&dq.state) != 0 {
		return fmt.Errorf("queue already closed")
	}
	dq.Queue.Push(x)
	return nil
}

func (dq *DelayQueue) Pop() (x DelayItem, ok bool) {
	if dq.Queue.Len() == 0 {
		return DelayItem{}, false
	}
	return dq.Queue.Pop().(DelayItem), true
}

func (dq *DelayQueue) Peek() (x *DelayItem) {
	if dq.Queue.Len() > 0 {
		return &dq.Queue[0]
	}
	return nil
}

func NewDelayQueue(size int, interval time.Duration) {
	dq := new(DelayQueue)
	dq.ticker = time.NewTicker(interval)
	dq.Queue = make([]DelayItem, 0, size)
	heap.Init(&dq.Queue)
}

type DelayItem struct {
	Data interface{}
	T    time.Time
}

type delayQueue []DelayItem

func (q delayQueue) Len() int { return len(q) }

func (q delayQueue) Less(i, j int) bool {
	return q[i].T.UnixNano() < q[j].T.UnixNano()
}

func (q delayQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *delayQueue) Push(x interface{}) {
	*q = append(*q, x.(DelayItem))
}

func (q *delayQueue) Pop() interface{} {
	old := *q
	n := len(old)
	x := old[n-1]
	*q = old[0 : n-1]
	return x
}
