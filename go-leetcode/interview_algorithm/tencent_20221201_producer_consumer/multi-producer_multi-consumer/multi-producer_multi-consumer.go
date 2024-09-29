package multi_producer_multi_consumer

import (
	"fmt"
	"sync"
	"time"
)

// 面试题1: 模拟生产者和消费者:
//	生产者1: 1,5,9...
//	生产者2: 2,6,10...
//	生产者3: 3,7,11...
//	生产者4: 4,8,12...

//	消费者1: 1,3...
//	消费者2: 2,4...

// monitor 使用有缓冲通道
// m: 生产者数量
// n: 消费者数量
func monitor(m, n int) {
	ch := make(chan int, 10)

	ms := make([]chan struct{}, m)
	for i := 0; i < m; i++ {
		ms[i] = make(chan struct{}, 1)
	}

	ns := make([]chan struct{}, n)
	for i := 0; i < n; i++ {
		ns[i] = make(chan struct{}, 1)
	}

	ms[0] <- struct{}{}
	ns[0] <- struct{}{}
	for i := 0; i < m; i++ {
		go func(i int) {
			start := i + 1
			for {
				<-ms[i]
				ch <- start
				//fmt.Printf("producer: %d; value: %d\n", i+1, start)
				time.Sleep(time.Second * 1)
				start += m
				ms[(i+1)%m] <- struct{}{}
			}
		}(i)
	}

	for j := 0; j < n; j++ {
		go func(j int) {
			for {
				<-ns[j]
				x := <-ch
				fmt.Printf("consumer: %d; value: %d\n", j+1, x)
				time.Sleep(time.Second * 1)
				ns[(j+1)%n] <- struct{}{}
			}
		}(j)
	}
	select {}
}

// monitor_2 使用无缓冲通道
// m: 生产者数量
// n: 消费者数量
func monitor_2(m, n int) {
	ch := make(chan int)

	ms := make([]chan struct{}, m)
	for i := 0; i < m; i++ {
		ms[i] = make(chan struct{})
	}

	ns := make([]chan struct{}, n)
	for i := 0; i < n; i++ {
		ns[i] = make(chan struct{})
	}

	for i := 0; i < m; i++ {
		go func(i int) {
			start := i + 1
			for {
				<-ms[i]
				ch <- start
				//fmt.Printf("producer: %d; value: %d\n", i+1, start)
				time.Sleep(time.Second * 1)
				start += m
				ms[(i+1)%m] <- struct{}{}
			}
		}(i)
	}

	for j := 0; j < n; j++ {
		go func(j int) {
			for {
				<-ns[j]
				x := <-ch
				fmt.Printf("consumer: %d; value: %d\n", j+1, x)
				time.Sleep(time.Second * 1)
				ns[(j+1)%n] <- struct{}{}
			}
		}(j)
	}
	ms[0] <- struct{}{}
	ns[0] <- struct{}{}
	select {}
}

// ---------------------------------------------------------------------------
// 面试题2: 使用多个协程顺序打印1~10

// seqPrint 使用sync.Cond
func seqPrint(number, target int) {
	var wg sync.WaitGroup
	cond := sync.NewCond(&sync.Mutex{})
	printNum := 0
	for i := 0; i < number; i++ {
		index := i
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			for {
				cond.L.Lock()
				for printNum%number != index {
					cond.Wait()
				}
				printNum++
				if printNum > target {
					cond.L.Unlock()
					cond.Broadcast()
					return
				}
				fmt.Println("goroutine:", index+1, "打印", printNum)
				cond.L.Unlock()
				cond.Broadcast()
			}
		}(index)
	}
	wg.Wait()
}

// seqPrint_2 使用通道
func seqPrint_2(number, target int) {
	ch := make(chan int)
	ms := make([]chan struct{}, number)
	for i := 0; i < number; i++ {
		ms[i] = make(chan struct{})
	}

	for i := 0; i < number; i++ {
		go func(i int) {
			start := i + 1
			for {
				<-ms[i]
				if start > target {
					close(ch)
					break
				}
				ch <- start
				fmt.Printf("goroutine: %d; value: %d\n", i+1, start)
				start += number
				ms[(i+1)%number] <- struct{}{}
			}
		}(i)
	}

	ms[0] <- struct{}{}
	for v := range ch {
		// do nothing
		_ = v
	}
}
