package redis

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

func TestFuck(t *testing.T) {
	h := md5.New()
	err := filepath.Walk("C:\\Users\\Chen\\Documents\\Tencent Files\\1003975097\\FileRecv\\IV邮件-0319",
		func(path string, info fs.FileInfo, err error) error {
			name := string(h.Sum([]byte(info.Name())))
			fmt.Println(name)
			os.Rename(path, name)
			return nil
		})
	fmt.Println(err)
}
