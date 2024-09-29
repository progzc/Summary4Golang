package chinaunicom_20240929_consumer_producer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ProducerConsumerModel
// 面试题3: 生产者与消费者（goroutine与channel）
// 1.有多个生产者不停生产整数1，有多个消费者同时进行消费。
// 2.所有的消费者的消费行为是：从channel中读取数据进行消费，并用一个公用的计数器进行累加。
// 3.某个消费者在做累加过程中，当计数器达到某数值时，通知所有生产者停止生产， 同时也通知其它消费者退出，然后自己也退出。
// 4.生产者一旦收到退出通知后，立即停止生产数据，并退出。
// 5.最后主协程等所有子协程全部退出后，主协程再退出。
func ProducerConsumerModel() {
	ch := make(chan int, 10)
	stop := make(chan struct{})
	once := sync.Once{}
	target := 1
	var threshold uint32 = 100
	var count uint32
	var wg sync.WaitGroup
	// 生产者
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				select {
				case ch <- target:
					time.Sleep(time.Millisecond * 10)
					fmt.Printf("第 %d 个生产者发送: %d\n", i, target)
				case <-stop:
					fmt.Printf("producer %d stopped\n", i)
					return
				}
			}
		}(i)
	}

	// 消费者
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				select {
				case v := <-ch:
					if x := atomic.AddUint32(&count, 1); x > threshold {
						once.Do(func() {
							close(stop)
						})
						fmt.Printf("consumer %d stopped\n", i)
						return
					} else {
						fmt.Printf("第 %d 个消费者消费第 %v 个%d\n", i, x, v)
					}
				case <-stop:
					fmt.Printf("consumer %d stopped\n", i)
					return
				}
			}

		}(i)
	}
	wg.Wait()
	fmt.Printf("main goroutine exited\n")
}
