package chapter11_context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestContext_1
// (1)context.Context的使用场景
//
//	a.上下文信息传递（request-scoped），比如处理http请求、在请求处理链路上传递信息。
//	b.控制子 goroutine 的运行。
//	c.超时控制的方法调用。
//	d.可以取消的方法调用。
//
// (2)context.Context的基本使用
//
//		type Context interface {
//			Deadline() (deadline time.Time, ok bool)
//			Done() <-chan struct{}
//			Err() error
//			Value(key interface{}) interface{}
//		}
//	a.Deadline方法：返回这个Context被取消的截止日期。如果没有设置截止日期，ok的值是false。
//				   后续每次调用这个对象的Deadline方法时，都会返回和第一次调用相同的结果。
//	b.Done方法：方法返回一个Channel对象。如果Done没有被close，Err方法返回nil；如果Done被close，Err方法会返回Done被close的原因。
//			   var Canceled = errors.New("context canceled")
//
//	c.Value方法：返回此ctx中和指定的key相关联的value。
//	d.两个常用的生成顶层Context的方法：context.Background() 和 context.TODO()：
//	  都会返回一个非nil的、空的Context，没有任何值，不会被cancel，不会超时，没有截止日期。
//	e.WithValue：链式查找
//	f.WithCancel：记住，不是只有你想中途放弃，才去调用cancel；只要你的任务正常完成了，就需要调用cancel（切记不要忘记调用cancel方法，可以放在defer中）。
//				  cancel时会向下传递，而不会向上传递。
//				  当这个cancelCtx的cancel函数被调用的时候，或者parent的Done被close的时候，这个cancelCtx的Done才会被close。
//	g.WithTimeout：可以设置超时时间。
//	h.WithDeadline：可以设置截止时间。WithDeadline（WithTimeout）返回的cancel一定要调用，并且要尽可能早地被调用，这样才能尽早释放资源，
//	                不要单纯地依赖截止时间被动取消。
//
// (3)context.Context的注意事项：
//
//	a.一般函数使用Context的时候，会把这个参数放在第一个参数的位置。
//	b.从来不把nil当做Context类型的参数值，可以使用context.Background()创建一个空的上下文对象，也不要使用nil。
//	c.Context只用来临时做函数之间的上下文透传，不能持久化Context或者把Context长久保存。
//	  把Context持久化到数据库、本地文件或者全局变量、缓存中都是错误的用法。
//	d.key的类型不应该是字符串类型或者其它内建类型，否则容易在包之间使用Context时候产生冲突。使用WithValue时，key的类型应该是自己定义的类型。
//	e.常常使用struct{}作为底层类型定义key的类型。对于exported key的静态类型，常常是接口或者指针。这样可以尽量减少内存分配。
//
// (4)timerCtx的Done被Close掉，主要是由下列时间触发的：
//
//	a.截止时间到了。
//	b.超时了。
//	c.cancel函数被调用。
//	d.parent的Done被close。
func TestContext_1(t *testing.T) {
	// e.context.WithValue的使用
	ctx := context.TODO()
	ctx = context.WithValue(ctx, "key1", "0001")
	ctx = context.WithValue(ctx, "key2", "0001")
	ctx = context.WithValue(ctx, "key3", "0001")
	ctx = context.WithValue(ctx, "key4", "0004")
	fmt.Println(ctx.Value("key1")) // 优先从子集从自己的存储中检查这个key，不存在的话会从parent中进行链式查找

	// g.WithTimeout：可以设置截止时间。WithDeadline（WithTimeout）返回的cancel一定要调用，并且要尽可能早地被调用，这样才能尽早释放资源，
	// 不要单纯地依赖截止时间被动取消。
	_, _ = slowOperationWithTimeout(ctx)
}

// TestContext_2
// (5)context.Context的基本使用
//
//	详细使用可参见：https://github.com/kat-co/concurrency-in-go-src
func TestContext_2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer func() {
			fmt.Println("goroutine exit")
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}()
	time.Sleep(time.Second)
	cancel()
	time.Sleep(2 * time.Second)
}

type Result struct{}

func slowOperationWithTimeout(ctx context.Context) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel() // 一旦慢操作完成就立马调用cancel
	return slowOperation(ctx)
}

func slowOperation(ctx context.Context) (Result, error) {
	// do something
	return Result{}, nil
}
