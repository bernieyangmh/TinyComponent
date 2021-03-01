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
	dq.Start()
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

	for i := 0; i < 100; i++ {
		dq.Push(DelayItem{Data: strconv.Itoa(i), T: time.Now().Add(time.Second * time.Duration(i))})
	}
	time.Sleep(time.Second * 5)
	dq.Close()
}

func Test_CloseWait(t *testing.T) {
	dq := NewDelayQueue(100, 0, time.Second)

	go func() {
		for i := range dq.readChan {
			if i.(int) == 90 {
				//会等全部readChan均处理完才发送fin
				time.Sleep(time.Second * 3)
			}
			fmt.Println(i)
		}
	}()
	dq.Start()
	for i := 0; i < 100; i++ {
		dq.Push(DelayItem{Data: i, T: time.Now().Add(time.Second * time.Duration(i))})
	}
	// 10s内只打印10个过期的item， 之后dq Close， 剩余的全部写入readChan
	time.Sleep(time.Second * 10)
	dq.Close()

	convey.Convey("chan should close", t, func(ctx convey.C) {
		for i := 0; i < 2; i++ {
			ctx.So(dq.Push(DelayItem{Data: i, T: time.Now().Add(time.Second * time.Duration(i))}),
				convey.ShouldBeError,
				fmt.Errorf("queue already closed"))
		}

		_, ok1 := <-dq.finChan
		ctx.So(ok1, convey.ShouldBeFalse)

		_, ok2 := <-dq.pushChan
		ctx.So(ok2, convey.ShouldBeFalse)

		_, ok3 := <-dq.popChan
		ctx.So(ok3, convey.ShouldBeFalse)

		_, ok4 := <-dq.pushChan
		ctx.So(ok4, convey.ShouldBeFalse)

		_, ok5 := <-dq.readChan
		ctx.So(ok5, convey.ShouldBeFalse)

		_, ok6 := <-dq.closeChan
		ctx.So(ok6, convey.ShouldBeFalse)

		_, ok7 := <-dq.popResChan
		ctx.So(ok7, convey.ShouldBeFalse)
	})

	fmt.Println(dq.Queue)
	t.Log("end")
}

func Test_TestLongLive(t *testing.T) {
	dq := NewDelayQueue(100, 0, time.Millisecond)
	dq.Start()
	go func() {
		for i := range dq.readChan {
			time.Sleep(time.Millisecond * 200)
			fmt.Println(i)
		}
	}()
	go func() {
		for i := range dq.readChan {
			time.Sleep(time.Millisecond * 200)
			fmt.Println(i)
		}
	}()

	go func() {
		for {
			dq.Push(DelayItem{Data: time.Now().String(), T: time.Now().Add(time.Second * time.Duration(2))})
			time.Sleep(time.Millisecond * 100)
		}
	}()
	go func() {
		for {
			dq.Push(DelayItem{Data: time.Now().String(), T: time.Now().Add(time.Second * time.Duration(2))})
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * 5)
			dq.Close()
		}
	}()

	select {}

}
