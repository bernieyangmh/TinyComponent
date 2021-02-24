package timingWheel

import (
	"container/heap"
	"time"
)

type TimingWheel struct {
	taskTiming *taskHeap
	funcMap    map[int64][]func()
	unit       time.Duration
}

func (tw *TimingWheel) Add(t time.Duration, f func()) {
	expire := time.Now().Add(t)
	heap.Push(tw.taskTiming, expire)
	if arr, ok := tw.funcMap[expire.UnixNano()]; ok {
		arr = append(arr, f)
	} else {
		arr = []func(){f}
		tw.funcMap[expire.UnixNano()] = arr
	}
}

func (tw *TimingWheel) Start() {
	for {
		deadline := tw.taskTiming.Pop().(time.Time)
		time.Sleep(time.Now().Sub(deadline))
		for _, f := range tw.funcMap[deadline.UnixNano()] {
			go f()
		}
		delete(tw.funcMap, deadline.UnixNano())
	}
}

type taskHeap []time.Time

func (th taskHeap) Len() int { return len(th) }

func (th taskHeap) Less(i, j int) bool {
	return th[i].Nanosecond() < th[j].Nanosecond()
}

func (th taskHeap) Swap(i, j int) {
	th[i], th[j] = th[j], th[i]

}

func (th *taskHeap) Push(x interface{}) {
	*th = append(*th, x.(time.Time))
}

func (th *taskHeap) Pop() interface{} {
	old := *th
	n := len(old)
	x := old[n-1]
	*th = old[0 : n-1]
	return x
}
