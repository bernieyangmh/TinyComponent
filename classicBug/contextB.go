package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	hctx, hcancel := context.WithCancel(ctx)
	go func() {
		//fmt.Println(foo(hctx, 1))
	}()

	hctx, hcancel = context.WithTimeout(ctx, time.Second*10)
	go func() {
		fmt.Println(foo(hctx, 2))
	}()
	time.Sleep(time.Second * 3)
	fmt.Println("start cancel")
	hcancel()
	select {}

}

func foo(ctx context.Context, i int) (err error) {
	// err = 2, nil
	a1 := &AAA{Name: "aaa", C: &CCC{Age: 18}}
	a2 := AAA{Name: "aaa", C: &CCC{Age: 18}}
	defer fmt.Println(i, err, fmt.Sprintf("%+v", a1), fmt.Sprintf("%+v", a2))
	i = 222
	a1.Name = "a"
	a2.Name = "a"
	a1.C.Age = 16
	a2.C.Age = 16

	time.Sleep(time.Second * 15)
	select {
	default:
	case <-ctx.Done():
		err = ctx.Err()
		return err
	}
	return nil
}

type AAA struct {
	Name string
	C    *CCC
}

type CCC struct {
	Age int
}
