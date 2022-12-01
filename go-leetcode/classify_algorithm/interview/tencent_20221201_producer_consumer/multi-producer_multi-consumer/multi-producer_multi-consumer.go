package multi_producer_multi_consumer

import (
	"fmt"
	"time"
)

// 模拟生产者和消费者:
//	生产者1: 1,5,9...
//	生产者2: 2,6,10...
//	生产者3: 3,7,11...
//	生产者4: 4,8,12...

//	消费者1: 1,3...
//	消费者2: 2,4...

// monitor
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
				if i+1 < m {
					ms[i+1] <- struct{}{}
				} else {
					ms[0] <- struct{}{}
				}
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
				if j+1 < n {
					ns[j+1] <- struct{}{}
				} else {
					ns[0] <- struct{}{}
				}
			}
		}(j)
	}
	select {}
}
