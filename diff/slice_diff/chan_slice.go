package slice_diff

type ChanSlice struct {
	c   chan interface{}
	buf []interface{}
}

func NewChanSlice(size, cap int) *ChanSlice {
	return &ChanSlice{c: make(chan interface{}, size), buf: make([]interface{}, 0, cap)}
}

func (c ChanSlice) Append(s interface{}) {
	c.c <- s
}

func (c ChanSlice) Start() {
	go func() {
		for s := range c.c {
			c.buf = append(c.buf, s)
		}
	}()
}
