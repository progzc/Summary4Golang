package chapter10

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type user struct {
	name string
	age  int8
}

var u = user{name: "Ankur", age: 25}
var g = &u

func modifyUser(pu *user) {
	fmt.Println("modifyUser Received Value", pu) // modifyUser Received Value &{Ankur Anand 100}
	pu.name = "Anand"
}

func printUser(u <-chan *user) {
	time.Sleep(2 * time.Second)
	fmt.Println("printUser GoRoutine called", <-u) // printUser GoRoutine called &{Ankur 25}
}

func Test_chan_1(t *testing.T) {
	c := make(chan *user, 5)
	c <- g
	fmt.Println(g) // &{Ankur 25}
	g = &user{name: "Ankur Anand", age: 100}
	go printUser(c)
	go modifyUser(g)
	time.Sleep(5 * time.Second)
	fmt.Println(g) // &{Anand 100}
}

func Test_chan_2(t *testing.T) {
	go func() {
		ch := make(chan struct{}, 0)
		ch <- struct{}{}
		<-ch
		fmt.Println("goroutine over")
	}()

	for {
		select {
		default:
			fmt.Println("tick...")
		}
		time.Sleep(time.Second)
	}
}

// Test_chan_3
// 功能: 顺序输出1,2,3,4,5,6,7,8,9,10
func Test_chan_3(t *testing.T) {
	n := 10
	var count uint32
	ch := make(chan uint32)
	for i := 0; i < n; i++ {
		go func(i uint32) {
			for {
				if val := atomic.LoadUint32(&count); val == i {
					ch <- i + 1
					atomic.AddUint32(&count, 1)
					break
				}
			}
		}(uint32(i))
	}

	for i := 0; i < n; i++ {
		fmt.Println(<-ch)
	}
}

// Test_chan_4
// 功能: 交替输出1,2,3,4
func Test_chan_4(t *testing.T) {
	cycle := 4
	var count uint32
	ch := make(chan uint32)
	limit := make(chan struct{}, 10)
	go func() {
		for i := 0; ; i++ {
			limit <- struct{}{}
			go func(i uint32) {
				defer func() {
					<-limit
				}()
				for {
					if val := atomic.LoadUint32(&count); val == i {
						ch <- uint32((int(i))%cycle + 1)
						time.Sleep(time.Second)
						atomic.AddUint32(&count, 1)
						break
					}
				}
			}(uint32(i))
		}
	}()

	go func() {
		for {
			fmt.Println(<-ch)
		}
	}()
	select {}
}

// Test_chan_5
// 功能: 4个协程输出顺序输出
func Test_chan_5(t *testing.T) {
	SeqPrint(4, 10)
}

func SeqPrint(number, target int) {
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
