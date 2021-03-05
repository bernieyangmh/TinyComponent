package slice_diff

import (
	"fmt"
	"sync"
)

type LockSlice struct {
	mutx *sync.Mutex
	buf  []interface{}
}

func NewLockSlice(cap int) *LockSlice {
	return &LockSlice{mutx: &sync.Mutex{}, buf: make([]interface{}, 0, cap)}
}

func (l *LockSlice) Append(s interface{}) {
	l.mutx.Lock()
	l.buf = append(l.buf, s)
	l.mutx.Unlock()
}

func (l *LockSlice) Read(index int) (interface{}, error) {
	if index >= len(l.buf) || index < 0 {
		return "", fmt.Errorf("index out of range")
	}
	return l.buf[index], nil
}
