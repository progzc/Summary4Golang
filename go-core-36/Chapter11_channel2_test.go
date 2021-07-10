package go_core_36

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// TestChannelImprove 单向通道
func TestChannelImprove(t *testing.T) {
	intChan1 := make(chan int, 3)
	// go会自动将双向通道转换为单向通道
	SendInt(intChan1)
}

// TestChannelImprove2 在函数类型中使用单向通道
func TestChannelImprove2(t *testing.T) {
	intChan2 := getIntChan()
	// 1.for语句会不断地尝试从通道intChan2中取出元素值。即使intChan2已经被关闭了，它也会在取出所有剩余的元素值之后再结束执行。
	// 2.通常，当通道intChan2中没有元素值时，这条for语句会被阻塞在有for关键字的那一行，直到有新的元素值可取。
	//   不过，由于这里的getIntChan函数会事先将intChan2关闭，所以它在取出intChan2中的所有元素值之后会直接结束执行。
	// 3.倘若通道intChan2的值为nil，那么这条for语句就会被永远地阻塞在有for关键字的那一行。
	// 另一种专门为了操作通道而存在的语句: select语句
	for elem := range intChan2 {
		fmt.Printf("The element in intChan2: %v\n", elem)
	}
}

// TestChannelImprove3 select语句的使用
func TestChannelImprove3(t *testing.T) {
	// 准备好几个通道
	intChannels := [3]chan int{
		make(chan int, 1),
		make(chan int, 1),
		make(chan int, 1),
	}
	// 随机选择一个通道,并向它发送元素值
	index := rand.Intn(3)
	fmt.Printf("The index: %d\n", index)
	intChannels[index] <- index
	// 哪一个通道中有可取的元素值,哪个对应的分支就会被执行
	select {
	case <-intChannels[0]:
		fmt.Println("The first candidate case is selected.")
	case <-intChannels[1]:
		fmt.Println("The second candidate case is selected.")
	case elem := <-intChannels[2]:
		fmt.Printf("The third candidate case is selected, the element is %d.\n", elem)
	default:
		fmt.Println("No candidate case is selected!")
	}
}

// TestChannelImprove4 select语句的使用2
func TestChannelImprove4(t *testing.T) {
	intChan := make(chan int, 1)
	// 一秒后关闭通道
	time.AfterFunc(time.Second, func() {
		close(intChan)
	})
	select {
	case _, ok := <-intChan:
		if !ok {
			fmt.Println("The candidate case is closed.")
			// break语句立即结束当前select语句的执行
			break
		}
		fmt.Println("The candidate case is selected.")
	}
}

var channels = [3]chan int{
	nil,
	//make(chan int, 1), // 改成这句会选择第2个case
	make(chan int), // 改成这句会选择默认的case
	nil,
}

var numbers = []int{1, 2, 3}

// TestChannelImprove5 select语句的使用3
func TestChannelImprove5(t *testing.T) {
	select {
	case getChan(0) <- getNumber(0):
		fmt.Println("The first candidate case is selected.")
	case getChan(1) <- getNumber(1):
		fmt.Println("The second candidate case is selected.")
	case getChan(2) <- getNumber(2):
		fmt.Println("The third candidate case is selected")
	default:
		fmt.Println("No candidate case is selected!")
	}
}

func getNumber(i int) int {
	fmt.Printf("numbers[%d]\n", i)
	return numbers[i]
}

func getChan(i int) chan int {
	fmt.Printf("channels[%d]\n", i)
	return channels[i]
}

// SendInt 该函数只能使用该通道进行发送操作
func SendInt(ch chan<- int) {
	ch <- rand.Intn(1000)
}

// getIntChan 该函数得到的通道只能进行接收操作
func getIntChan() <-chan int {
	num := 5
	ch := make(chan int, num)
	for i := 0; i < num; i++ {
		ch <- i
	}
	close(ch)
	return ch
}
