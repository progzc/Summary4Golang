package go_explore

import (
	"context"
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync/atomic"
	"testing"
	"time"
)

// https://zhuanlan.zhihu.com/p/343761986 sync.singleflight 到底怎么用才对？

type Result string

func find(ctx context.Context, query string) (Result, error) {
	return Result(fmt.Sprintf("result for %q", query)), nil
}

func TestSingleFlight(t *testing.T) {
	var g singleflight.Group
	const n = 5
	waited := int32(n)
	done := make(chan struct{})
	key := "https://weibo.com/1227368500/H3GIgngon"
	for i := 0; i < n; i++ {
		go func(j int) {
			v, _, shared := g.Do(key, func() (interface{}, error) {
				ret, err := find(context.Background(), key)
				return ret, err
			})
			if atomic.AddInt32(&waited, -1) == 0 {
				close(done)
			}
			fmt.Printf("index: %d, val: %v, shared: %v\n", j, v, shared)
		}(i)
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		fmt.Println("Do hangs")
	}
}
