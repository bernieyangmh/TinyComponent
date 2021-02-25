package delayQueue

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func Test_TestAddSequence(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	dq := NewDelayQueue(100, 0, time.Second)
	for i := 0; i < 100; i++ {
		dq.Push(DelayItem{Data: i, T: time.Now().Add(time.Second * time.Duration(rand.Intn(101)))})
	}
	for {
		di, ok := dq.Pop()
		if !ok {
			break
		}
		t.Logf("%+v", di)
	}
}

func Test_SimpleExample(*testing.T) {
	dq := NewDelayQueue(100, 0, time.Second)
	go func() {
		for i := range dq.readChan {
			fmt.Println(i)
		}
	}()
	go func() {
		for i := range dq.readChan {
			fmt.Println(i)
		}
	}()
	dq.Start()

	for i := 0; i < 100; i++ {
		dq.Push(DelayItem{Data: strconv.Itoa(i), T: time.Now().Add(time.Second * time.Duration(i))})
	}
	time.Sleep(time.Second * 10)
	dq.Close()
}

func Test_CloseWait(t *testing.T)  {
	dq := NewDelayQueue(100, 0, time.Second)

	go func() {
		for i := range dq.readChan {
			if i.(int) == 90 {
				time.Sleep(time.Second*3)
			}
			fmt.Println(i)
		}
	}()
	dq.Start()

	for i := 0; i < 100; i++ {
		dq.Push(DelayItem{Data: i, T: time.Now().Add(time.Second * time.Duration(i))})
	}
	time.Sleep(time.Second * 5)
	dq.Close()

}