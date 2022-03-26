package chapter18_group_opt

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-pkgz/syncs"
	"github.com/vardius/gollback"
	"golang.org/x/sync/errgroup"
)

// TestGroupOpt_1
// (1)分组编排(ErrGroup)：
//	设计：
//		https://github.com/golang/sync/tree/master/errgroup
//		https://golang.org/x/sync/errgroup
//	作用：将一个通用的父任务拆成几个小任务并发执行。
//	举例：ErrGroup就是用来应对这种场景的，它和WaitGroup有些类似，但是它提供功能更加丰富：
//		a.和Context集成。
//		b.error向上传播，可以把子任务的错误传递给Wait的调用者。
// (2)ErrGroup的基本使用
//	a.简单的使用：返回第一个错误
//	b.简单的使用：返回所有的错误
//	c.任务执行流水线：
//		代码示例：https://pkg.go.dev/golang.org/x/sync/errgroup#example-Group--Pipeline
// (3)分组编排的扩展库
//	a.bilibili/errgroup：可以使用一个固定数量的 goroutine 处理子任务。
//		https://pkg.go.dev/github.com/bilibili/kratos/pkg/sync/errgroup
//		缺点：
//			i)一旦你设置了并发数，超过并发数的子任务需要等到调用者调用Wait之后才会执行，而不是只要goroutine空闲下来，就去执行。
//			ii)在高并发的情况下，如果任务数大于设定的goroutine的数量，并且这些任务被集中加入到Group中，
//			   这个库的处理方式是把子任务加入到一个数组中，但是，这个数组不是线程安全的，有并发问题。
//	b.neilotoole/errgroup：增加了可以控制并发goroutine的功能。
//		https://github.com/neilotoole/errgroup
//	b.facebookgo/errgroup：对标准库WaitGroup的扩展，提供的Wait方法可以返回error。
//		https://github.com/facebookarchive/errgroup
// (4)其他分组并发原语
//	a.go-pkgz/syncs：提供了两个Group并发原语，分别是SizedGroup和ErrSizedGroup。
//	  其中SizedGroup内部是使用信号量和WaitGroup实现的，它通过信号量控制并发的goroutine数量，或者是不控制goroutine数量，只控制子任务并发执行时候的数量（通过）。
//		https://github.com/go-pkgz/syncs
//	b.vardius/gollback：用来处理一组子任务的执行的，不过它解决了ErrGroup收集子任务返回结果的痛点。
//		https://github.com/vardius/gollback
//	c.AaronJan/Hunch：提供的功能和gollback类似，不过它提供的方法更多，而且它提供的和gollback相应的方法，也有一些不同。
//		https://github.com/AaronJan/Hunch
//	d.mdlayher/schedgroup：一个worker pool，可以指定任务在某个时间或者某个时间之后执行。
//		https://github.com/mdlayher/schedgroup
func TestGroupOpt_1(t *testing.T) {
	// a.简单的使用
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

// TestGroupOpt_2
// 简单的使用：返回所有的错误
func TestGroupOpt_2(t *testing.T) {
	// b.简单的使用：返回所有的错误
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

// TestGroupOpt_3
// 一个多阶段的pipeline.使用有限的goroutine计算每个文件的md5值
func TestGroupOpt_3(t *testing.T) {
	// c.任务执行流水线
	m, err := MD5All(context.Background(), ".")
	if err != nil {
		log.Fatal(err)
	}

	for k, sum := range m {
		fmt.Printf("%s:\t%x\n", k, sum)
	}
}

type result struct {
	path string
	sum  [md5.Size]byte
}

// 遍历根目录下所有的文件和子文件夹,计算它们的md5的值.
func MD5All(ctx context.Context, root string) (map[string][md5.Size]byte, error) {
	g, ctx := errgroup.WithContext(ctx)
	paths := make(chan string) // 文件路径channel

	g.Go(func() error {
		defer close(paths) // 遍历完关闭paths chan
		return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			//将文件路径放入到paths
			// do something
			return nil
		})
	})

	// 启动20个goroutine执行计算md5的任务，计算的文件由上一阶段的文件遍历子任务生成.
	c := make(chan result)
	const numDigesters = 20
	for i := 0; i < numDigesters; i++ {
		g.Go(func() error {
			for path := range paths { // 遍历直到paths chan被关闭
				// 计算path的md5值，放入到c中
				// do something
				_ = path
			}
			return nil
		})
	}
	go func() {
		g.Wait() // 20个goroutine以及遍历文件的goroutine都执行完
		close(c) // 关闭收集结果的chan
	}()

	m := make(map[string][md5.Size]byte)
	for r := range c { // 将md5结果从chan中读取到map中,直到c被关闭才退出
		m[r.path] = r.sum
	}

	// 再次调用Wait，依然可以得到group的error信息
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return m, nil
}

// TestGroupOpt_4
// SizedGroup的使用
func TestGroupOpt_4(t *testing.T) {
	swg := syncs.NewSizedGroup(10) // 控制子任务的并发数量
	// swg := syncs.NewSizedGroup(10, syncs.Preemptive) // 控制goroutine的数量
	var c uint32

	// 执行1000个子任务，只会有10个goroutine去执行
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

// TestGroupOpt_5
// gollback的使用
//	a.All方法：等待所有的异步函数（AsyncFunc）都执行完才返回，而且返回结果的顺序和传入的函数的顺序保持一致。
//  	      第一个返回参数是子任务的执行结果，第二个参数是子任务执行时的错误信息。
//		func All(ctx context.Context, fns ...AsyncFunc) ([]interface{}, []error)
//	b.Race方法：跟All方法类似，只不过，在使用Race方法的时候，只要一个异步函数执行没有错误，就立马返回，而不会返回所有的子任务信息。
//	           如果所有的子任务都没有成功，就会返回最后一个error信息。
//		func Race(ctx context.Context, fns ...AsyncFunc) (interface{}, error)
//	c.Retry方法：Retry不是执行一组子任务，而是执行一个子任务。
// 		    	如果子任务执行失败，它会尝试一定的次数，如果一直不成功，就会返回失败错误，如果执行成功，它会立即返回；
//		    	如果retires等于0，它会永远尝试，直到成功。
//		func Retry(ctx context.Context, retires int, fn AsyncFunc) (interface{}, error)
func TestGroupOpt_5(t *testing.T) {
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

// TestGroupOpt_6
// gollback的使用
func TestGroupOpt_6(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试5次，或者超时返回
	res, err := gollback.Retry(ctx, 5, func(ctx context.Context) (interface{}, error) {
		return nil, errors.New("failed")
	})

	fmt.Println(res) // 输出结果
	fmt.Println(err) // 输出错误信息

}
