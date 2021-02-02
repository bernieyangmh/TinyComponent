package barrierChan

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

var total = int32(0)

func BenchmarkMutilChan(b *testing.B) {
	rand.Seed(time.Now().Unix())

	// 创建栅栏对象
	bc := NewBarrier(5, time.Millisecond*3000)
	// 达到的效果：前n个协程调用Wait()阻塞，第n个调用后n个协程全部唤醒
	for i := 0; i < 500; i++ {
		tmp := i
		time.Sleep(time.Millisecond * 5)
		if i == 50 {
			bc.Close()
		}
		go func() {
			bc.Wait(tmp)
			fmt.Println(fmt.Sprintf("finish (%d)", tmp))
			atomic.AddInt32(&total, 1)
		}()

	}

	time.Sleep(time.Second * 100)
}
