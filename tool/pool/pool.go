package pool

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
)

type pool struct {
	isClose int32
	pool    chan struct{}
}

type Timeout struct {
	error
}

type Closed struct {
	error
}

func NewPool(max int) *pool {
	return &pool{
		pool: make(chan struct{}, max),
	}
}

func (p *pool) Close() error {
	if atomic.LoadInt32(&p.isClose) == 1 {
		return Closed{
			error: errors.New("池已经关闭了"),
		}
	}
	for atomic.CompareAndSwapInt32(&p.isClose, 0, 1) {
	}
	for len(p.pool) != 0 {
	}
	close(p.pool)
	return nil
}

func (p *pool) Submit(ctx context.Context, f func()) error {
	if atomic.LoadInt32(&p.isClose) == 1 {
		return Closed{
			error: errors.New("池已经关闭了"),
		}
	}
	select {
	case p.pool <- struct{}{}:
		go func() {
			defer func() {
				err := recover()
				if err != nil {
					fmt.Println(err)
				}
			}()
			f()
			<-p.pool
		}()
		return nil
	case <-ctx.Done():
		return Timeout{
			error: errors.New("等待超时"),
		}
	}
}
