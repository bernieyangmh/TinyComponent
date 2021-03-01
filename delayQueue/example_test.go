package delayQueue

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func Test_TestAddSequence(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	dq := NewDelayQueue(100, 0, time.Second)
	dq.Start()
	recordMap := make(map[int]DelayItem)

	for i := 0; i < 100; i++ {
		d := DelayItem{Data: i, T: time.Now().Add(time.Second * time.Duration(rand.Intn(101)))}
		dq.Push(d)
		recordMap[i] = d
	}
	convey.Convey("delayTime Must before now", t, func(ctx convey.C) {
		for {
			for i := range dq.readChan {
				t.Logf("%+v, %+v", recordMap[i.(int)], time.Now())
				ctx.So(recordMap[i.(int)].T.UnixNano(), convey.ShouldBeLessThanOrEqualTo, time.Now().UnixNano())
			}
		}
	})

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

func Test_CloseWait(t *testing.T) {
	dq := NewDelayQueue(100, 0, time.Second)

	go func() {
		for i := range dq.readChan {
			if i.(int) == 90 {
				time.Sleep(time.Second * 3)
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
