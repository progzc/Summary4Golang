package chapter17_SingleFlight_CyclicBarrier

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/marusama/cyclicbarrier"
	"golang.org/x/sync/semaphore"
	"golang.org/x/sync/singleflight"
)

// TestSingleFlight_1
// (1)SingleFlight(读请求合并)
//	SingleFlight：将并发请求合并成一个请求，以减少对下层服务的压力。
//	a.含义：在处理多个goroutine同时调用同一个函数的时候，只让一个goroutine去调用这个函数，等到这个goroutine返回结果的时候，
//		   再把结果返回给这几个同时调用的goroutine，这样可以减少并发调用的数量。
//	b.比较：SingleFlight和sync.Once有什么区别?
//	  sync.Once主要是用在单次初始化场景中，而SingleFlight主要用在合并并发请求的场景中，尤其是缓存场景。
//	c.适用场景：在面对秒杀等大并发请求的场景，而且这些请求都是读请求时，你就可以把这些请求合并为一个请求，这样，你就可以将后端服务的压力从n降到1。
//			  尤其是在面对后端是数据库这样的服务的时候，采用SingleFlight可以极大地提高性能。
//	d.实现原理：用互斥锁Mutex(提供并发时的读写保护) + Map(用来保存同一个key的正在处理（in flight）的请求)
//		https://github.com/golang/go/blob/b1b67841d1e229b483b0c9dd50ddcd1795b0f90f/src/net/lookup.go
//		https://pkg.go.dev/golang.org/x/sync/singleflight
//	e.使用SingleFlight的知名项目：
//		https://github.com/golang/groupcache
//		https://github.com/cockroachdb/cockroach
//		https://github.com/coredns/coredns
func TestSingleFlight_1(t *testing.T) {
	var (
		g                singleflight.Group
		errorNotExist    = errors.New("not exist")
		wg               sync.WaitGroup
		getData          func(key string) (string, error) // 获取数据
		getDataFromCache func(key string) (string, error) // 模拟从cache中获取值，cache中无该值
		getDataFromDB    func(key string) (string, error) // 模拟从数据库中获取值
	)
	getDataFromCache = func(key string) (string, error) {
		return "", errorNotExist
	}
	getDataFromDB = func(key string) (string, error) {
		log.Printf("get %s from database", key)
		return "data", nil
	}
	getData = func(key string) (string, error) {
		data, err := getDataFromCache(key)
		if err == errorNotExist {
			//模拟从db中获取数据
			v, err, _ := g.Do(key, func() (interface{}, error) {
				return getDataFromDB(key)
				// set cache
			})
			if err != nil {
				log.Println(err)
				return "", err
			}
			// set cache
			data = v.(string)
		} else if err != nil {
			return "", err
		}
		return data, nil
	}

	wg.Add(10)

	//模拟10个并发
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			data, err := getData("key")
			if err != nil {
				log.Print(err)
				return
			}
			log.Println(data)
		}()
	}
	wg.Wait()
}

// TestCyclicBarrier_1
// (2)CyclicBarrier(循环栅栏)
//	a.含义：可重用的栅栏，用来控制一组请求同时执行。
//		允许一组goroutine彼此等待，到达一个共同的执行点。同时，因为它可以被重复使用，所以叫循环栅栏。
//		具体的机制是，大家都在栅栏前等待，等全部都到齐了，就抬起栅栏放行。
//	b.比较：WaitGroup和CyclicBarrier的区别？
//		不同点：
//			i)功能不同：
//				WaitGroup(减法计数器)：一个goroutine等待一组goroutine都执行完到达同一个执行点。
//				CyclicBarrier(加法计数器)：允许一组goroutine彼此等待，到达一个共同的执行点，然后再放行。
//			ii)WaitGroup零值可用，CyclicBarrier零值不可用(需要New创建)。
//		相同点：都是可重用的（但是复用的方式不一样，WaitGroup需要重新调用Add函数进行设置）。
//	c.使用场景：重复进行一组goroutine同时执行的场景中。
//	d.经典案例：一氧化二氢制造工厂
func TestCyclicBarrier_1(t *testing.T) {
	//用来存放水分子结果的channel
	var ch chan string
	releaseHydrogen := func() {
		ch <- "H"
	}
	releaseOxygen := func() {
		ch <- "O"
	}

	// 300个原子，300个goroutine,每个goroutine并发的产生一个原子
	var N = 100
	ch = make(chan string, N*3)

	h2o := New()

	// 用来等待所有的goroutine完成
	var wg sync.WaitGroup
	wg.Add(N * 3)

	// 200个氢原子goroutine
	for i := 0; i < 2*N; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2o.hydrogen(releaseHydrogen)
			wg.Done()
		}()
	}
	// 100个氧原子goroutine
	for i := 0; i < N; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2o.oxygen(releaseOxygen)
			wg.Done()
		}()
	}

	//等待所有的goroutine执行完
	wg.Wait()

	// 结果中肯定是300个原子
	if len(ch) != N*3 {
		t.Fatalf("expect %d atom but got %d", N*3, len(ch))
	}

	// 每三个原子一组，分别进行检查。要求这一组原子中必须包含两个氢原子和一个氧原子，这样才能正确组成一个水分子。
	var s = make([]string, 3)
	for i := 0; i < N; i++ {
		s[0] = <-ch
		s[1] = <-ch
		s[2] = <-ch
		sort.Strings(s)

		water := s[0] + s[1] + s[2]
		if water != "HHO" {
			t.Fatalf("expect a water molecule but got %s", water)
		}
	}
}

// 定义水分子合成的辅助数据结构
type H2O struct {
	semaH *semaphore.Weighted         // 氢原子的信号量
	semaO *semaphore.Weighted         // 氧原子的信号量
	b     cyclicbarrier.CyclicBarrier // 循环栅栏，用来控制合成
}

func New() *H2O {
	return &H2O{
		semaH: semaphore.NewWeighted(2), // 氢原子需要两个
		semaO: semaphore.NewWeighted(1), // 氧原子需要一个
		b:     cyclicbarrier.New(3),     // 需要三个原子才能合成
	}
}

func (h2o *H2O) hydrogen(releaseHydrogen func()) {
	h2o.semaH.Acquire(context.Background(), 1)

	releaseHydrogen()                 // 输出H
	h2o.b.Await(context.Background()) // 等待栅栏放行
	h2o.semaH.Release(1)              // 释放氢原子空槽
}

func (h2o *H2O) oxygen(releaseOxygen func()) {
	h2o.semaO.Acquire(context.Background(), 1)

	releaseOxygen()                   // 输出O
	h2o.b.Await(context.Background()) // 等待栅栏放行
	h2o.semaO.Release(1)              // 释放氢原子空槽
}

// ------------------------------------使用WaitGroup来实现这个水分子制造工厂------------------------------------------------
type H2O_2 struct {
	semaH *semaphore.Weighted
	semaO *semaphore.Weighted
	wg    sync.WaitGroup //将循环栅栏替换成WaitGroup
}

func New2() *H2O_2 {
	var wg sync.WaitGroup
	wg.Add(3)

	return &H2O_2{
		semaH: semaphore.NewWeighted(2),
		semaO: semaphore.NewWeighted(1),
		wg:    wg,
	}
}

func (h2o *H2O_2) hydrogen(releaseHydrogen func()) {
	h2o.semaH.Acquire(context.Background(), 1)
	releaseHydrogen()

	// 标记自己已达到，等待其它goroutine到达
	h2o.wg.Done()
	h2o.wg.Wait()

	h2o.semaH.Release(1)
}

func (h2o *H2O_2) oxygen(releaseOxygen func()) {
	h2o.semaO.Acquire(context.Background(), 1)
	releaseOxygen()

	// 标记自己已达到，等待其它goroutine到达
	h2o.wg.Done()
	h2o.wg.Wait()
	//都到达后重置wg
	h2o.wg.Add(3)

	h2o.semaO.Release(1)
}
