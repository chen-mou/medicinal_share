package redis

import (
	"fmt"
	"sync"
	"testing"
)

func do(f func()) {
	f()
}

func TestLock(t *testing.T) {
	f := func(val int) {
		fmt.Println(val)
	}
	for i := 0; i < 10; i++ {
		do(func() {
			f(i)
		})
	}
}

func TestAntiShake(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			if !AntiShake("test") {
				fmt.Println("防抖成功")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
