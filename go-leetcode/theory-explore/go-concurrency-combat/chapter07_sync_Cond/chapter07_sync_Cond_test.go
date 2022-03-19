package chapter07_sync_Cond

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestCond_1 通知与唤醒
// (1)从开发实践上，我们真正使用 Cond 的场景比较少，因为一旦遇到需要使用 Cond 的场景，我们更多地会使用 Channel 的方式，
//    这才是更地道的 Go 语言的写法。
// (2)sync.Cond的基本使用
//	a.Signal方法：允许调用者 Caller 唤醒一个等待此 Cond 的 goroutine；如果 Cond 等待队列中有一个或者多个等待的 goroutine，
//  	         则需要从等待队列中移除第一个 goroutine 并把它唤醒。相当于Java种的notify方法。
//				 调用 Signal 方法时，不强求你一定要持有 c.L 的锁。
//	b.Broadcast方法：允许调用者 Caller 唤醒所有等待此 Cond 的 goroutine；如果 Cond 等待队列中有一个或者多个等待的 goroutine，
//	                则清空所有等待的 goroutine，并全部唤醒。相当于Java种的notifyAll方法。
//					调用 Broadcast 方法时，也不强求你一定持有 c.L 的锁。
//	c.Wait方法：会把调用者 Caller 放入 Cond 的等待队列中并阻塞，直到被 Signal 或者 Broadcast 的方法从等待队列中移除并唤醒。
//			   相当于Java种的wait方法。调用 Wait 方法时必须要持有 c.L 的锁。
// (3)字段组成
//	a.notify：notifyList类型，即一个等待/通知队列。
//	b.L：Locker类型，互斥锁。
//	c.checker：copyChecker类型，可以在运行时检查 Cond 是否被复制使用。
//	d.noCopy：使用go vet检查是否被复制。
// (4)实现原理：互斥锁+等待队列(主要由runtime实现)
//	a.Signal 和 Broadcast 只涉及到 notifyList 数据结构，不涉及到锁。
//	b.Wait 把调用者加入到等待队列时需要加锁，加入之后会释放锁，在被唤醒之后还会请求锁。
//	  调用者加入到等待队列后之所以要释放锁，是为了让其他 Wait 的调用者有机会加入到 notify 队列中。
//	  在阻塞休眠期间，调用者是不持有锁的，这样能让其他 goroutine 有机会检查或者更新等待变量。
// (5)常见错误
//	a.调用 Wait 的时候没有加锁。
//	b.只调用了一次 Wait，没有检查等待条件是否满足，结果条件没满足，程序就继续执行了。
//	  一定要记住，waiter goroutine 被唤醒不等于等待条件被满足，只是有 goroutine 把它唤醒了而已，等待条件有可能已经满足了，
//	  也有可能不满足，我们需要进一步检查。你也可以理解为，等待者被唤醒，只是得到了一次检查的机会而已。
// (6)Cond 在实际项目中被使用的机会比较少, 其原因如下：
//	a.同样的场景我们会使用其他的并发原语来替代。Go 特有的 Channel 类型，有一个应用很广泛的模式就是通知机制，这个模式使用起来也特别简单。
//	  所以很多情况下，我们会使用 Channel 而不是 Cond 实现 wait/notify 机制。
//	b.对于简单的 wait/notify 场景，比如等待一组 goroutine 完成之后继续执行余下的代码，我们会使用 WaitGroup 来实现。
//	  因为 WaitGroup 的使用方法更简单，而且不容易出错。
// (7)sync.Cond和Channel的区别：sync.Cond有三点特性是Channel 无法替代的。
//	a.Cond 和一个 Locker 关联，可以利用这个 Locker 对相关的依赖条件更改提供保护。
//	c.Cond 的 Broadcast 方法可以被重复调用。等待条件再次变成不满足的状态后，我们又可以调用 Broadcast 再次唤醒等待的 goroutine。
//	  这也是 Channel 不能支持的，Channel 被 close 掉了之后不支持再 open。
//	b.Cond 可以同时支持 Signal 和 Broadcast 方法，而 Channel 只能同时支持其中一种。
// (8)sync.Cond和sync.WaitGroup的区别：
//	a.WaitGroup 是主 goroutine 等待确定数量的子 goroutine 完成任务；
//	  而 Cond 是等待某个条件满足，这个条件的修改可以被任意多的 goroutine 更新。
//	b.Cond 的 Wait 不关心也不知道其他 goroutine 的数量，只关心等待条件。
//	c.Cond 还有单个通知的机制，也就是 Signal 方法。
func TestCond_1(t *testing.T) {
	c := sync.NewCond(&sync.Mutex{})
	var ready int

	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
			c.L.Lock()
			ready++ // 注意事项：条件变量的更改，需要原子操作或者互斥锁保护
			c.L.Unlock()
			log.Printf("运动员#%d 已准备就绪\n", i)
			// 广播唤醒所有的等待者
			c.Broadcast() // 注意事项：由于等待者只有一个，这里也可以换成c.Signal()
		}(i)
	}

	// 易错点：a.调用 Wait 的时候没有加锁。
	c.L.Lock()
	// 易错点：b.只调用了一次 Wait，没有检查等待条件是否满足，结果条件没满足，程序就继续执行了。
	for ready != 10 {
		c.Wait()
		log.Println("裁判员被唤醒一次")
	}
	c.L.Unlock()

	//所有的运动员是否就绪
	log.Println("所有运动员都准备就绪。比赛开始，3，2，1, ......")
}

// TestCond_2
// (9)使用sync.Cond实现生产者消费者问题
func TestCond_2(t *testing.T) {
	rand.Seed(time.Now().UnixNano()) // 设置随机数种子
	quit := make(chan bool)          // 创建用于结束通信的 channel

	product := make(chan int, 3)        // 产品区（公共区）使用channel 模拟
	cond := sync.NewCond(&sync.Mutex{}) // 创建互斥锁和条件变量

	for i := 0; i < 5; i++ { // 5个消费者
		go producer(product, i+1, cond)
	}
	for i := 0; i < 3; i++ { // 3个生产者
		go consumer(product, i+1, cond)
	}
	<-quit // 主go程阻塞 不结束
}

// producer 生产者
func producer(out chan<- int, idx int, cond *sync.Cond) {
	for {
		cond.L.Lock()       // 条件变量对应互斥锁加锁
		for len(out) == 3 { // 产品区满 等待消费者消费
			cond.Wait() // 挂起当前go程， 等待条件变量满足，被消费者唤醒
		}
		num := rand.Intn(1000) // 产生一个随机数
		out <- num             // 写入到 channel 中 （生产）
		fmt.Printf("%dth 生产者，产生数据 %3d, 公共区剩余%d个数据\n", idx, num, len(out))
		cond.L.Unlock()         // 生产结束，解锁互斥锁
		cond.Signal()           // 唤醒 阻塞的 消费者
		time.Sleep(time.Second) // 生产完休息一会，给其他go程执行机会
	}
}

// consumer 消费者
func consumer(in <-chan int, idx int, cond *sync.Cond) {
	for {
		cond.L.Lock()      // 条件变量对应互斥锁加锁（与生产者是同一个）
		for len(in) == 0 { // 产品区为空 等待生产者生产
			cond.Wait() // 挂起当前go程， 等待条件变量满足，被生产者唤醒
		}
		num := <-in // 将 channel 中的数据读走 （消费）
		fmt.Printf("%dth 消费者, 消费数据 %3d,公共区剩余%d个数据\n", idx, num, len(in))
		cond.L.Unlock()                    // 消费结束，解锁互斥锁
		cond.Signal()                      // 唤醒 阻塞的 生产者
		time.Sleep(time.Millisecond * 500) //消费完 休息一会，给其他go程执行机会
	}
}

// TestCond_3
// (10)使用sync.Cond实现一个容量有限的 queue
func TestCond_3(t *testing.T) {
	_ = NewQueue(10)
}

type Queue struct {
	cond *sync.Cond
	data []interface{}
	cap  int
	logs []string
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		cond: &sync.Cond{L: &sync.Mutex{}},
		data: make([]interface{}, 0),
		cap:  capacity,
		logs: make([]string, 0)}
}

func (q *Queue) Enqueue(d interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for len(q.data) == q.cap {
		q.cond.Wait()
	}
	// FIFO入队
	q.data = append(q.data, d)
	// 记录操作日志
	q.logs = append(q.logs, fmt.Sprintf("En %v\n", d))
	// 通知其他waiter进行Dequeue或Enqueue操作
	q.cond.Broadcast()
}

func (q *Queue) Dequeue() (d interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for len(q.data) == 0 {
		q.cond.Wait()
	}
	// FIFO出队
	d = q.data[0]
	q.data = q.data[1:]
	// 记录操作日志
	q.logs = append(q.logs, fmt.Sprintf("De %v\n", d))
	// 通知其他waiter进行Dequeue或Enqueue操作
	q.cond.Broadcast()
	return
}

func (q *Queue) Len() int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return len(q.data)
}

func (q *Queue) String() string {
	var b strings.Builder
	for _, l := range q.logs {
		b.WriteString(l)
	}
	return b.String()
}
