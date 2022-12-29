package redis

import (
	"fmt"
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
