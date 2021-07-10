package go_core_36

import (
	"fmt"
	"testing"
)

// TestChannel 通道的使用
func TestChannel(t *testing.T) {
	ch1 := make(chan int, 3)
	// 发送表达式
	ch1 <- 2
	ch1 <- 1
	ch1 <- 3
	// 接收表达式
	elem1 := <-ch1
	fmt.Printf("The first element received from channel ch1: %v\n", elem1)
}

// TestChannel2 通道的阻塞情况
func TestChannel2(t *testing.T) {
	// 示例1。
	ch1 := make(chan int, 1)
	ch1 <- 1
	//ch1 <- 2 // 通道已满，因此这里会造成阻塞。

	// 示例2。
	ch2 := make(chan int, 1)
	//elem, ok := <-ch2 // 通道已空，因此这里会造成阻塞。
	//_, _ = elem, ok
	ch2 <- 1

	// 示例3。
	var ch3 chan int
	//ch3 <- 1 // 通道的值为nil，因此这里会造成永久的阻塞！
	//<-ch3 // 通道的值为nil，因此这里会造成永久的阻塞！
	_ = ch3
}

// TestChannel3 通道的注意事项
func TestChannel3(t *testing.T) {
	// 1. 通道关闭后还进行收发操作时会引发panic
	// 2. 关闭一个已关闭的通道会引发panic
	// 3. 通道关闭后,不能发送数据,但可以接收数据,直至接收操作完毕后才会安全退出
	// 4. 基于第3点可知: 通过接收表达式的第二个结果值,来判断通道是否关闭是可能有延时的
	// 5. 基于以上四点,不要让接收方关闭通道,应该让发送方关闭通道
	ch1 := make(chan int, 2)
	// 发送方。
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("Sender: sending element %v...\n", i)
			ch1 <- i
		}
		fmt.Println("Sender: close the channel...")
		close(ch1)
	}()

	// 接收方。
	for {
		elem, ok := <-ch1
		if !ok {
			fmt.Println("Receiver: closed channel")
			break
		}
		fmt.Printf("Receiver: received an element: %v\n", elem)
	}

	fmt.Println("End.")
}
