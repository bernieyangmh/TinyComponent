package delayQueue

import (
	"container/heap"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

//1.队列的初始cap, 避免频繁分配
//2.readChan的buf, 0是阻塞Chan，需考虑消费速度
//3.每次判断到达定时的间隔，按需设置
func NewDelayQueue(size, buff int, interval time.Duration) *DelayQueue {
	dq := new(DelayQueue)
	dq.ticker = time.NewTicker(interval)
	dq.Queue = make([]DelayItem, 0, size)
	dq.readChan = make(chan interface{}, buff)
	dq.closeChan = make(chan int8)
	dq.popChan = make(chan struct{})
	dq.pushChan = make(chan DelayItem)
	dq.popResChan = make(chan DelayItem)
	dq.mutx = new(sync.Mutex)
	heap.Init(&dq.Queue)
	return dq
}

type DelayQueue struct {
	ticker     *time.Ticker
	Queue      delayQueue
	state      int32
	readChan   chan interface{}
	closeChan  chan int8
	popChan    chan struct{}
	pushChan   chan DelayItem
	popResChan chan DelayItem
	mutx       *sync.Mutex
}

// Start之前必须有goroutine消费readChan
func (dq *DelayQueue) Start() {
	go func() {
		for {
			select {
			case <-dq.ticker.C:
				now := time.Now()
				top, ok := dq.Peek()
				if !ok {
					continue
				}
				if now.Sub(top.T) < 0 {
					continue
				}
				for dq.state == 0 {
					top, ok := dq.Peek()
					if !ok {
						break
					}
					if now.Sub(top.T) < 0 {
						break
					}
					di, ok := dq.Pop()
					if !ok {
						break
					}
					if now.Sub(di.T) >= 0 {
						dq.readChan <- di.Data
						continue
					}
				}
			case _, ok := <-dq.closeChan:
				if !ok {
					continue
				}
				for {
					di, ok := dq.Pop()
					if !ok {
						break
					}
					dq.readChan <- di.Data
					continue
				}
				close(dq.closeChan)
				break
				// already clear
			}
		}
	}()

	go func() {
		for {
			select {
			case <-dq.closeChan:
				// TODO: update to quit gracefully
				// TODO: maybe need to dump state somewhere?
				return
			case <-dq.popChan:
				popD, ok := dq.pop().(DelayItem)
				if ok {
					dq.popResChan <- popD
				}
			case pushD := <-dq.pushChan:
				dq.push(pushD)
			}
		}
	}()
}

// 关闭之后，剩余的数据会被立刻写入chan，直到消费完成前Close方法都会阻塞
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
	dq.pushChan <- x
	return nil
}

func (dq *DelayQueue) push(x DelayItem) error {
	heap.Push(&dq.Queue, x)
	return nil
}

func (dq *DelayQueue) Pop() (x DelayItem, ok bool) {
	if dq.Queue.Len() == 0 {
		return DelayItem{}, false
	}
	dq.popChan <- struct{}{}
	return <-dq.popResChan, true
}

func (dq *DelayQueue) pop() (x interface{}) {
	return heap.Pop(&dq.Queue)
}

func (dq *DelayQueue) Peek() (x DelayItem, ok bool) {
	if dq.Queue.Len() > 0 {
		return dq.Queue[0], true
	}
	return DelayItem{}, false
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

//https://husobee.github.io/heaps/golang/safe/2016/09/01/safe-heaps-golang.html
