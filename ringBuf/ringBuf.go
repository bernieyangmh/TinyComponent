package ringBuf

// k8s.io/utils/buffer/ring_growing.go

type RingBuf struct {
	data     []interface{}
	size     uint64
	index    uint64
	avaCount uint64
}

func NewRingBuf(size uint64) *RingBuf {
	return &RingBuf{size: size, data: make([]interface{}, size)}
}

func (r *RingBuf) Read() (data interface{}, ok bool) {
	if r.avaCount == 0 {
		return nil, false
	}
	r.avaCount--
	item := r.data[r.index]
	r.data[r.index] = nil
	if r.index == r.size-1 {
		r.index = 0
	} else {
		r.index++
	}
	return item, true
}

func (r *RingBuf) Write(data interface{}) {
	if r.avaCount == r.size {
		var newSize uint64
		if r.size < 1024 {
			newSize = r.size * 2
		} else if r.size < 102400 {
			newSize = r.size + 1<<10
		} else {
			newSize = r.size + 1<<12
		}

		newData := make([]interface{}, newSize)

		cur := r.index + r.avaCount
		if cur < r.size {
			copy(newData, r.data[r.index:cur])
		} else {
			used := copy(newData, r.data[r.index:])
			copy(newData[used:], r.data[:(cur%r.size)])
		}

		r.index = 0
		r.data = newData
		r.size = newSize
	}
	r.data[(r.avaCount+r.index)%r.size] = data
	r.avaCount++
}
