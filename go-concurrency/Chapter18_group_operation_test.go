package go_concurrency

import (
	"context"
	"errors"
	"fmt"
	fb_errgroup "github.com/facebookgo/errgroup"
	"github.com/go-pkgz/syncs"
	"github.com/vardius/gollback"
	"golang.org/x/sync/errgroup"
	"sync/atomic"
	"testing"
	"time"
)

func TestErrGroup(t *testing.T) {
	var g errgroup.Group

	// 启动第一个子任务,它执行成功
	g.Go(func() error {
		time.Sleep(5 * time.Second)
		fmt.Println("exec #1")
		return nil
	})
	// 启动第二个子任务，它执行失败
	g.Go(func() error {
		time.Sleep(10 * time.Second)
		fmt.Println("exec #2")
		return errors.New("failed to exec #2")
	})

	// 启动第三个子任务，它执行成功
	g.Go(func() error {
		time.Sleep(15 * time.Second)
		fmt.Println("exec #3")
		return nil
	})
	// 等待三个任务都完成
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully exec all")
	} else {
		fmt.Println("failed:", err)
	}
}

func TestErrGroup2(t *testing.T) {
	var g errgroup.Group
	var result = make([]error, 3)

	// 启动第一个子任务,它执行成功
	g.Go(func() error {
		time.Sleep(5 * time.Second)
		fmt.Println("exec #1")
		result[0] = nil // 保存成功或者失败的结果
		return nil
	})

	// 启动第二个子任务，它执行失败
	g.Go(func() error {
		time.Sleep(10 * time.Second)
		fmt.Println("exec #2")

		result[1] = errors.New("failed to exec #2") // 保存成功或者失败的结果
		return result[1]
	})

	// 启动第三个子任务，它执行成功
	g.Go(func() error {
		time.Sleep(15 * time.Second)
		fmt.Println("exec #3")
		result[2] = nil // 保存成功或者失败的结果
		return nil
	})

	if err := g.Wait(); err == nil {
		fmt.Printf("Successfully exec all. result: %v\n", result)
	} else {
		fmt.Printf("failed: %v\n", result)
	}
}

func TestErrGroup3(t *testing.T) {
	var g fb_errgroup.Group
	g.Add(3)

	// 启动第一个子任务,它执行成功
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("exec #1")
		g.Error(errors.New("failed to exec #1"))
		g.Done()
	}()

	// 启动第二个子任务，它执行失败
	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("exec #2")
		g.Error(errors.New("failed to exec #2"))
		g.Done()
	}()

	// 启动第三个子任务，它执行成功
	go func() {
		time.Sleep(15 * time.Second)
		fmt.Println("exec #3")
		g.Done()
	}()

	// 等待所有的goroutine完成，并检查error
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully exec all")
	} else {
		fmt.Println("failed:", err)
	}
}

func TestSyns(t *testing.T) {
	// 设置goroutine数是10
	swg := syncs.NewSizedGroup(10)
	//swg := syncs.NewSizedGroup(10, syncs.Preemptive)
	var c uint32

	// 执行1000个子任务，只会有10个goroutine去执行（实测还是会执行这个1000个子任务）
	for i := 0; i < 1000; i++ {
		swg.Go(func(ctx context.Context) {
			time.Sleep(5 * time.Millisecond)
			atomic.AddUint32(&c, 1)
		})
	}

	// 等待任务完成
	swg.Wait()
	// 输出结果
	fmt.Println(c)
}

func TestGollback(t *testing.T) {
	rs, errs := gollback.All( // 调用All方法
		context.Background(),
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(3 * time.Second)
			return 1, nil // 第一个任务没有错误，返回1
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("failed") // 第二个任务返回一个错误
		},
		func(ctx context.Context) (interface{}, error) {
			return 3, nil // 第三个任务没有错误，返回3
		},
	)

	fmt.Println(rs)   // 输出子任务的结果
	fmt.Println(errs) // 输出子任务的错误信息
}

func TestGollback2(t *testing.T) {
	rs, errs := gollback.Race( // 调用All方法
		context.Background(),
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(3 * time.Second)
			return 1, nil // 第一个任务没有错误，返回1
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("failed") // 第二个任务返回一个错误
		},
		func(ctx context.Context) (interface{}, error) {
			return 3, nil // 第三个任务没有错误，返回3
		},
	)

	fmt.Println(rs)   // 输出子任务的结果
	fmt.Println(errs) // 输出子任务的错误信息
}

func TestGollback3(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试5次，或者超时返回
	res, err := gollback.Retry(ctx, 5, func(ctx context.Context) (interface{}, error) {
		return nil, errors.New("failed")
	})

	fmt.Println(res) // 输出结果
	fmt.Println(err) // 输出错误信息
}
