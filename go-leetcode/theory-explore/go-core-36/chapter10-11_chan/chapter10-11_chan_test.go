package chapter10

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// TestChan_1 通道的使用
// (1)通道容量为0时,称为非缓冲通道;通道容量大于0时,称为缓冲通道;通道类似于先进先出的队列
//  对于缓冲通道: 发送端-->数据-->副本-->接收端-->删除副本
//	对于非缓冲通道: 发送端-->数据-->数据直接复制到接收端
//	注意：这里复制是指浅拷贝(事实上go里面没有深拷贝?copy内建方法?)
// (2)基本同特性:
//	a.对于同一个通道，发送操作之间是互斥的，接收操作之间也是互斥的
//	b.发送操作和接收操作中对元素值的处理都是不可分割的
//	c.发送操作在完全完成之前会被阻塞。接收操作也是如此
// (3)发送或接收通道出现阻塞时,阻塞的goroutine会放置在等待队列,当阻塞接触时,通知的顺序是公平的,即满足先进先出
// (4)针对nil的通道,其发送和接收操作都会永远地处于阻塞状态
// (5)使用过程中出现panic的情况:
//	a.当堆关闭的通道执行发送操作,就会引发panic (针对已初始化但并未关闭的通道,收发操作一定不会引发panic)
//  b.当关闭一个已经关闭的通道时,也会出现panic
// (6)针对通道的接收操作:
//	a.当将表达式接收结果同时赋值两个变量时,第二个值为false说明通道已经关闭
//	b.当通道已经关闭但还有元素未取出时,第二个值为true(结合a知,通过第二个值判断通道是否关闭是可能有延时的,但为false说明通道一定已经关闭了)
//	c.千万不要让接收方关闭通道,而应该让发送方关闭通道(否则可能会出现panic)
func TestChan_1(t *testing.T) {
	//(1)通道容量为0时,称为非缓冲通道;通道容量大于0时,称为缓冲通道;通道类似于先进先出的队列
	ch1 := make(chan int, 3)
	ch1 <- 2
	ch1 <- 1
	ch1 <- 3
	elem1 := <-ch1
	fmt.Printf("The first element received from channel ch1: %v\n",
		elem1)

	// -----------------------------------------------------------------------
	//(4)针对nil的通道,其发送和接收操作都会永远地处于阻塞状态
	ch2 := make(chan int, 1)
	ch2 <- 1
	//ch2 <- 2 // 通道已满，因此这里会造成阻塞。

	ch3 := make(chan int, 1)
	//elem, ok := <-ch3 // 通道已空，因此这里会造成阻塞。
	//_, _ = elem, ok
	ch3 <- 1

	var ch4 chan int
	//ch4 <- 1 // 通道的值为nil，因此这里会造成永久的阻塞！
	//<-ch4 // 通道的值为nil，因此这里会造成永久的阻塞！
	_ = ch4

	// -----------------------------------------------------------------------
	// 通道的收发操作
	ch5 := make(chan int, 2)
	// 发送方。
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("Sender: sending element %v...\n", i)
			ch5 <- i
		}
		fmt.Println("Sender: close the channel...")
		close(ch5)
	}()

	// 接收方。
	for {
		elem, ok := <-ch5
		if !ok {
			fmt.Println("Receiver: closed channel")
			break
		}
		fmt.Printf("Receiver: received an element: %v\n", elem)
	}

	fmt.Println("End.")
}

// TestChan_2 通道里的数据复制
func TestChan_2(t *testing.T) {
	ch := make(chan []int, 1)
	s1 := []int{1, 2, 3}
	ch <- s1
	s2 := <-ch
	s2[0] = 100
	fmt.Println(s1, s2) //[100 2 3] [100 2 3]

	ch2 := make(chan [3]int, 1)
	s3 := [3]int{1, 2, 3}
	ch2 <- s3
	s4 := <-ch2
	s3[0] = 100
	fmt.Println(s3, s4) //[100 2 3] [1 2 3]
}

// TestChan_3 通道的高级操作
// (1)单向收发通道: <-chan int	chan<- int
// (2)单向通道最主要的用途就是约束其他代码的行为
// (3)单向通道的使用
//	a.在接口中限制通道的行为
//	b.在函数间传递时,可以把一个元素类型匹配的双向通道传给单向通道（本质是Go会自动地把双向通道转换为函数所需的单向通道）
//	c.还可以在函数声明的结果列表中使用单向通道
// (4)循环迭代取出通道中的元素
//	a.即使通道已经关闭了,也可以取出通道中的元素;在通道关闭后,迭代取出所有元素后会结束执行
//	b.当通道intChan2中没有元素值时，这条for语句会被阻塞在有for关键字的那一行,直到有新的元素值可取
//	c.如果通道为nil,则迭代会永远阻塞
func TestChan_3(t *testing.T) {
	//(1)单向收发通道: <-chan int	chan<- int
	//(2)单向通道最主要的用途就是约束其他代码的行为
	var uselessChan = make(chan<- int, 1)        // 只能发不能收的通道
	var anotherUselessChan = make(<-chan int, 1) // 只能收不能发的通道
	fmt.Printf("The useless channels: %v, %v\n", uselessChan, anotherUselessChan)

	//(3)单向通道的使用
	//a.在接口中限制通道的行为
	//b.在函数间传递时,可以把一个元素类型匹配的双向通道传给单向通道（本质是Go会自动地把双向通道转换为函数所需的单向通道）
	//c.还可以在函数声明的结果列表中使用单向通道
	intChan1 := make(chan int, 3)
	SendInt(intChan1)
	intChan2 := getIntChan()

	//(4)循环迭代取出通道中的元素
	//a.即使通道已经关闭了,也可以取出通道中的元素;在通道关闭后,迭代取出所有元素后会结束执行
	//b.当通道intChan2中没有元素值时，这条for语句会被阻塞在有for关键字的那一行,直到有新的元素值可取
	//c.如果通道为nil,则迭代会永远阻塞
	for elem := range intChan2 {
		fmt.Printf("The element in intChan2: %v\n", elem)
	}
}

type Notifier interface {
	SendInt(ch chan<- int)
}

func SendInt(ch chan<- int) {
	ch <- rand.Intn(1000)
}

func getIntChan() <-chan int {
	num := 5
	ch := make(chan int, num)
	for i := 0; i < num; i++ {
		ch <- i
	}
	close(ch)
	return ch
}

// TestChan_4 select妙用
// (1)select的基本使用
//	a.如果加入了默认分支,那么无论涉及通道操作的表达式是否有阻塞,select语句都不会被阻塞
//	b.如果没有加入默认分支那一旦,所有的case表达式都没有满足求值条件,那么select语句就会被阻塞
//	c.在很多时候，我们需要通过接收表达式的第二个结果值来判断通道是否已经关闭(这对于程序逻辑和程序性能都有好处)
//	d.select语句只能对其中的每一个case表达式各求值一次,如果需要连续监听通道,一般采用for+select;
//但是要注意:简单地在select语句的分支中使用break语句，只能结束当前的select语句的执行，而并不会对外层的for语句产生作用。
// (2)select的进阶总结
//	a.select语句包含的候选分支中的case表达式都会在该语句执行开始时先被求值，并且求值的顺序是依从代码编写的顺序从上到下的
//	b.对于每一个case表达式，如果其中的发送表达式或者接收表达式在被求值时，相应的操作正处于阻塞状态，那么对该case表达式的求值就是不成功的。
//在这种情况下，我们可以说，这个case表达式所在的候选分支是不满足选择条件的。
//	c.仅当select语句中的所有case表达式都被求值完毕后，它才会开始选择候选分支。这时候，它只会挑选满足选择条件的候选分支执行。
//	如果所有的候选分支都不满足选择条件，那么默认分支就会被执行。如果这时没有默认分支，那么select语句就会立即进入阻塞状态，直到至少有一个候选分支满足选择条件为止。
//	一旦有一个候选分支满足选择条件，select语句（或者说它所在的 goroutine）就会被唤醒，这个候选分支就会被执行。
//	d.如果select语句发现同时有多个候选分支满足选择条件，那么它就会用一种伪随机的算法在这些分支中选择一个并执行。
//	e.一条select语句中只能够有一个默认分支。并且，默认分支只在无候选分支可选时才会被执行，这与它的编写位置无关。
//	f.select语句的每次执行，包括case表达式求值和分支选择，都是独立的。不过，至于它的执行是否是并发安全的，就要看其中的case表达式以及分支中，是否包含并发不安全的代码了。
//	(使用go run race可以检测是否并发安全)
func TestChan_4(t *testing.T) {
	//(1)select的基本使用
	intChannels := [3]chan int{
		make(chan int, 1),
		make(chan int, 1),
		make(chan int, 1),
	}
	index := rand.Intn(3)
	fmt.Printf("The index: %d\n", index)
	intChannels[index] <- index
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

	// --------------------------------------------------------------------------
	//但是要注意:简单地在select语句的分支中使用break语句，只能结束当前的select语句的执行，而并不会对外层的for语句产生作用。
	intChan := make(chan int, 1)
	// 一秒后关闭通道。
	time.AfterFunc(time.Second, func() {
		close(intChan)
	})
	select {
	case _, ok := <-intChan:
		if !ok {
			fmt.Println("The candidate case is closed.") // 1s后会输出这句
			break
		}
		fmt.Println("The candidate case is selected.")
	}

	// --------------------------------------------------------------------------
	// (2)select的进阶总结
	//	a.select语句包含的候选分支中的case表达式都会在该语句执行开始时先被求值，并且求值的顺序是依从代码编写的顺序从上到下的
	//	b.对于每一个case表达式，如果其中的发送表达式或者接收表达式在被求值时，相应的操作正处于阻塞状态，那么对该case表达式的求值就是不成功的。
	//在这种情况下，我们可以说，这个case表达式所在的候选分支是不满足选择条件的。
	//	c.仅当select语句中的所有case表达式都被求值完毕后，它才会开始选择候选分支。这时候，它只会挑选满足选择条件的候选分支执行。
	//	如果所有的候选分支都不满足选择条件，那么默认分支就会被执行。如果这时没有默认分支，那么select语句就会立即进入阻塞状态，直到至少有一个候选分支满足选择条件为止。
	//	一旦有一个候选分支满足选择条件，select语句（或者说它所在的 goroutine）就会被唤醒，这个候选分支就会被执行。
	//	d.如果select语句发现同时有多个候选分支满足选择条件，那么它就会用一种伪随机的算法在这些分支中选择一个并执行。
	//	e.一条select语句中只能够有一个默认分支。并且，默认分支只在无候选分支可选时才会被执行，这与它的编写位置无关。
	//	f.select语句的每次执行，包括case表达式求值和分支选择，都是独立的。不过，至于它的执行是否是并发安全的，就要看其中的case表达式以及分支中，是否包含并发不安全的代码了。
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

var channels = [3]chan int{
	nil,
	make(chan int),
	nil,
}

var numbers = []int{1, 2, 3}

func getNumber(i int) int {
	fmt.Printf("numbers[%d]\n", i)
	return numbers[i]
}

func getChan(i int) chan int {
	fmt.Printf("channels[%d]\n", i)
	return channels[i]
}

// TestChan_6 思考题1
// Q:如果在select语句中发现某个通道已关闭，那么应该怎样屏蔽掉它所在的分支？
// A:发现某个channel被关闭后，为了防止再次进入这个分支，可以把这个channel重新赋值成为一个长度为0的非缓冲通道，这样这个case就一直被阻塞了
func TestChan_6(t *testing.T) {
	ch1 := make(chan int, 1)
	for {
		select {
		case _, ok := <-ch1:
			if !ok {
				ch1 = make(chan int)
			}
		//case ...:
		////
		default:
			////
		}
	}
}

// TestChan_7 思考题2
// Q:在select语句与for语句联用时，怎样直接退出外层的for语句？
// A:可以用 break和标签配合使用，直接break出指定的循环体，或者goto语句直接跳转到指定标签执行
func TestChan_7(t *testing.T) {
	// 方法一:break配合标签：
	ch1 := make(chan int, 1)
	time.AfterFunc(time.Second, func() { close(ch1) })
loop1:
	for {
		select {
		case _, ok := <-ch1:
			if !ok {
				break loop1
			}
			fmt.Println("ch1")
		}
	}
	fmt.Println("END")

	// 方法二:goto配合标签
	ch2 := make(chan int, 1)
	time.AfterFunc(time.Second, func() { close(ch2) })
	for {
		select {
		case _, ok := <-ch2:
			if !ok {
				goto loop2
			}
			fmt.Println("ch1")
		}
	}
loop2:
	fmt.Println("END")

	// 方法三:直接return
}
