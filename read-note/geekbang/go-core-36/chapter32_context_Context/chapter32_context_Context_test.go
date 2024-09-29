package chapter32_context_Context

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// 推荐的两篇博客:
// https://blog.golang.org/pipelines
// https://blog.golang.org/context

// TestContext_1
// (1)总结常见的几种并发模式:
//
//	a.channel
//	b.sync.Mutex及sync.RWMutex
//	c.条件变量sync.Cond
//	d.原子操作
//	e.sync.WaitGroup
//	c.context.Context
//
// (2)Q:怎样使用context包中的程序实体，实现一对多的如下goroutine 协作流程？
//
//	A:使用context.Context+原子操作
//
// (3)Context类型实际上是一个接口类型，而context包中实现该接口的所有私有类型，都是基于某个数据类型的指针类型，
//
//	所以，如此传播并不会影响该类型值的功能和安全。
//
// (4)context包中还包含了四个用于繁衍Context值的函数: 具体用法详见https://pkg.go.dev/context
//
//	a.WithCancel：通过调用该函数，获得了一个衍生自上下文根节点的Context值，和一个用于触发撤销信号的函数。
//	b.WithDeadline：可以被用来产生一个会定时撤销的parent的子值
//	c.WithTimeout：都可以被用来产生一个会定时撤销的parent的子值
//	d.WithValue：可以通过调用它，产生一个会携带额外数据的parent的子值(需要注意的是：调用context.WithValue函数得到的Context值是不可撤销)
//
// (5)Q:“可撤销的”在context包中代表着什么？“撤销”一个Context值又意味着什么？
//
//	   A:有两个方法与撤销有关,一个是Done方法，另一个是Err方法
//			a.Done方法原本会返回一个元素类型为struct{}的接收通道（并不用户传递值,只是用于感知撤销信号）;撤销后,接收通道会关闭,Done()方法就会感知到,变为非阻塞状态
//			b.撤销后,Err方法会显示“撤销”的具体原因。
//
// (6)Q:撤销信号是如何在上下文树中传播的?
//
//	A:从上到下,树状传播。
//
// (7)Q:怎样通过Context值携带数据？怎样从中获取数据？
//
//	A:需要注意: WithValue函数在产生新的Context值（以下简称含数据的Context值）的时候需要三个参数，即：父值、键和值。
//	          与“字典对于键的约束”类似，这里键的类型必须是可判等的。
//	  查找数据: Context类型的Value方法就是被用来获取数据的。其查找方式为 从下到上（链式）,查到为止。
//
// (8)Q:Context值在传达撤销信号的时候是广度优先的，还是深度优先的？其优势和劣势都是什么?
//
//	A: 它是深度优先的。
//	   其优势和劣势都是：直接分支的产生时间越早，其中的所有子节点就会越先接收到信号。
//	   至于什么时候是优势、什么时候是劣势还要看具体的应用场景。例如，如果子节点的存续时间与资源的消耗是正相关的，那么这可能就是一个优势。
//	   但是，如果每个分支中的子节点都很多，而且各个分支中的子节点的产生顺序并不依从于分支的产生顺序，那么这种优势就很可能会变成劣势。
//	   最终的定论还是要看测试的结果。
func TestContext_1(t *testing.T) {
	coordinateWithContext()

	keys := []myKey{
		myKey(20),
		myKey(30),
		myKey(60),
		myKey(61),
	}
	values := []string{
		"value in node2",
		"value in node3",
		"value in node6",
		"value in node6Branch",
	}

	rootNode := context.Background()
	node1, cancelFunc1 := context.WithCancel(rootNode)
	defer cancelFunc1()

	// 示例1。
	node2 := context.WithValue(node1, keys[0], values[0])
	node3 := context.WithValue(node2, keys[1], values[1])
	fmt.Printf("The value of the key %v found in the node3: %v\n",
		keys[0], node3.Value(keys[0]))
	fmt.Printf("The value of the key %v found in the node3: %v\n",
		keys[1], node3.Value(keys[1]))
	fmt.Printf("The value of the key %v found in the node3: %v\n",
		keys[2], node3.Value(keys[2]))
	fmt.Println()

	// 示例2。
	node4, _ := context.WithCancel(node3)
	node5, _ := context.WithTimeout(node4, time.Hour)
	fmt.Printf("The value of the key %v found in the node5: %v\n",
		keys[0], node5.Value(keys[0]))
	fmt.Printf("The value of the key %v found in the node5: %v\n",
		keys[1], node5.Value(keys[1]))
	fmt.Println()

	// 示例3。
	node6 := context.WithValue(node5, keys[2], values[2])
	fmt.Printf("The value of the key %v found in the node6: %v\n",
		keys[0], node6.Value(keys[0]))
	fmt.Printf("The value of the key %v found in the node6: %v\n",
		keys[2], node6.Value(keys[2]))
	fmt.Println()

	// 示例4。
	node6Branch := context.WithValue(node5, keys[3], values[3])
	fmt.Printf("The value of the key %v found in the node6Branch: %v\n",
		keys[1], node6Branch.Value(keys[1]))
	fmt.Printf("The value of the key %v found in the node6Branch: %v\n",
		keys[2], node6Branch.Value(keys[2]))
	fmt.Printf("The value of the key %v found in the node6Branch: %v\n",
		keys[3], node6Branch.Value(keys[3]))
	fmt.Println()

	// 示例5。
	node7, _ := context.WithCancel(node6)
	node8, _ := context.WithTimeout(node7, time.Hour)
	fmt.Printf("The value of the key %v found in the node8: %v\n",
		keys[1], node8.Value(keys[1]))
	fmt.Printf("The value of the key %v found in the node8: %v\n",
		keys[2], node8.Value(keys[2]))
	fmt.Printf("The value of the key %v found in the node8: %v\n",
		keys[3], node8.Value(keys[3]))
}

type myKey int

func coordinateWithContext() {
	total := 12
	var num int32
	fmt.Printf("The number: %d [with context.Context]\n", num)
	cxt, cancelFunc := context.WithCancel(context.Background())
	for i := 1; i <= total; i++ {
		go addNum(&num, i, func() {
			if atomic.LoadInt32(&num) == int32(total) {
				cancelFunc()
			}
		})
	}
	<-cxt.Done()
	fmt.Println("End.")
}

// addNum 用于原子地增加一次numP所指的变量的值。
func addNum(numP *int32, id int, deferFunc func()) {
	defer func() {
		deferFunc()
	}()
	for i := 0; ; i++ {
		currNum := atomic.LoadInt32(numP)
		newNum := currNum + 1
		time.Sleep(time.Millisecond * 200)
		if atomic.CompareAndSwapInt32(numP, currNum, newNum) {
			fmt.Printf("The number: %d [%d-%d]\n", newNum, id, i)
			break
		} else {
			//fmt.Printf("The CAS operation failed. [%d-%d]\n", id, i)
		}
	}
}
