package go_group

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

const (
	k64 = 64 << 10
)

type Group struct {
	err error
	wg  sync.WaitGroup
	ctx context.Context

	errOnce  sync.Once
	workOnce sync.Once

	ch chan func(ctx context.Context) error

	cancel func()
}

func WithContext(ctx context.Context) *Group {
	return &Group{ctx: ctx}
}

func (g *Group) do(f func(ctx context.Context) error) {
	ctx := g.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, k64)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("Group panic recovered: %s\n%s", r, buf)
		}
		if err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
		g.wg.Done()
	}()
	err = f(ctx)
}
