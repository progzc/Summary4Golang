package chapter13_15_channel

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// 概要
//	(1)Channel在任务编排中的应用
//		a.Or-Done 模式
//		b.扇入 模式
//		c.扇出 模式
//		d.Stream 模式
//		e.map-reduce 模式

// TestTaskSchedule_1
// (1)复杂编排之 Or-Done 模式
// 	概念：有多个任务，只要有任意一个任务执行完，就结束。
//	实现方式：一种时递归、另一种是反射
func TestTaskSchedule_1(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(10*time.Second),
		sig(20*time.Second),
		sig(30*time.Second),
		sig(40*time.Second),
		sig(50*time.Second),
		sig(01*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}

// or 使用递归的方式实现 Or-Done 模式
func or(channels ...<-chan interface{}) <-chan interface{} {
	// 特殊情况，只有零个或者1个chan
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2: // 2个也是一种特殊情况
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default: //超过两个，二分法递归处理
			m := len(channels) / 2
			select {
			case <-or(channels[:m]...):
			case <-or(channels[m:]...):
			}
		}
	}()

	return orDone
}

// orReflect 使用反射的方式实现 Or-Done 模式
func orReflect(channels ...<-chan interface{}) <-chan interface{} {
	//特殊情况，只有0个或者1个
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		// 利用反射构建SelectCase
		var cases []reflect.SelectCase
		for _, c := range channels {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		// 随机选择一个可用的case
		reflect.Select(cases)
	}()

	return orDone
}

// TestTaskSchedule_2
// (2)复杂编排之 扇入 模式
//	概念：多个输入，一个输出。
//	实现方式：一种是递归、另一种是反射
func TestTaskSchedule_2(t *testing.T) {
	_ = fanInRec()
}

// fanInRec 使用递归方式实现 扇入 模式
func fanInRec(chans ...<-chan interface{}) <-chan interface{} {
	switch len(chans) {
	case 0:
		c := make(chan interface{})
		close(c)
		return c
	case 1:
		return chans[0]
	case 2:
		return mergeTwo(chans[0], chans[1])
	default:
		m := len(chans) / 2
		return mergeTwo(
			fanInRec(chans[:m]...),
			fanInRec(chans[m:]...))
	}
}

// mergeTwo 扇入的主要逻辑
func mergeTwo(a, b <-chan interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		for a != nil || b != nil { //只要还有可读的chan
			select {
			case v, ok := <-a:
				if !ok { // a 已关闭，设置为nil
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok { // b 已关闭，设置为nil
					b = nil
					continue
				}
				c <- v
			}
		}
	}()
	return c
}

// fanInReflect 使用反射方式实现 扇入 模式
func fanInReflect(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		// 构造SelectCase slice
		var cases []reflect.SelectCase
		for _, c := range chans {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		// 循环，从cases中选择一个可用的
		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok { // 此channel已经close
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			out <- v.Interface()
		}
	}()
	return out
}

// TestTaskSchedule_3
// (3)复杂编排之 扇出 模式
//	概念：一个输入，多个输出。（经常用在观察者模式中）
//	实现方式：一种是递归、另一种是反射
func TestTaskSchedule_3(t *testing.T) {
	ch := make(chan interface{}, 10)
	out := []chan interface{}{
		make(chan interface{}, 5),
		make(chan interface{}, 5),
		make(chan interface{}, 5),
	}
	fanOut(ch, out, true)
}

func fanOut(ch <-chan interface{}, out []chan interface{}, async bool) {
	go func() {
		defer func() { //退出时关闭所有的输出chan
			for i := 0; i < len(out); i++ {
				close(out[i])
			}
		}()

		for v := range ch { // 从输入chan中读取数据
			v := v
			for i := 0; i < len(out); i++ {
				i := i
				if async { //异步
					go func() {
						out[i] <- v // 放入到输出chan中,异步方式
					}()
				} else {
					out[i] <- v // 放入到输出chan中，同步方式
				}
			}
		}
	}()
}

// TestTaskSchedule_4
// (4)复杂编排之 Stream 模式
//	概念：把Channel当作流式管道使用，提供跳过几个元素，或者是只取其中的几个元素等方法。
//	实现方式：使用通道即可实现
func TestTaskSchedule_4(t *testing.T) {
	done := make(chan struct{})
	values := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	inStream := asStream(done, values...)
	outStream := takeN(done, inStream, 6)
	for v := range outStream {
		fmt.Println(v)
	}
	close(done)
}

// asStream 将 数据源(切片) 转换为 流
func asStream(done <-chan struct{}, values ...interface{}) <-chan interface{} {
	s := make(chan interface{}) //创建一个unbuffered的channel
	go func() {                 // 启动一个goroutine，往s中塞数据
		defer close(s)             // 退出时关闭chan
		for _, v := range values { // 遍历数组
			select {
			case <-done:
				return
			case s <- v: // 将数组元素塞入到chan中
			}
		}
	}()
	return s
}

// takeN 只取流中的前n个元素
func takeN(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{}) // 创建输出流
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ { // 只读取前num个元素
			select {
			case <-done:
				return
			case takeStream <- <-valueStream: //从输入流中读取元素
			}
		}
	}()
	return takeStream
}

// TestTaskSchedule_5
// (5)复杂编排之 map-reduce 模式
//	概念：一种面向大规模数据处理的并行计算方式。
//	实现方式：需要分两个步骤，
//		第一步是映射（map），处理队列中的数据。
//		第二步是约（reduce），把列表中的每一个元素按照一定的处理方式处理成结果，放入到结果队列中。
func TestTaskSchedule_5(t *testing.T) {
	in := asStream2(nil)

	// map操作: 乘以10
	mapFn := func(v interface{}) interface{} {
		return v.(int) * 10
	}

	// reduce操作: 对map的结果进行累加
	reduceFn := func(r, v interface{}) interface{} {
		return r.(int) + v.(int)
	}

	sum := reduce(mapChan(in, mapFn), reduceFn) //返回累加结果
	fmt.Println(sum)
}

// 生成一个数据流
func asStream2(done <-chan struct{}) <-chan interface{} {
	s := make(chan interface{})
	values := []int{1, 2, 3, 4, 5}
	go func() {
		defer close(s)
		for _, v := range values { // 从数组生成
			select {
			case <-done:
				return
			case s <- v:
			}
		}
	}()
	return s
}

func mapChan(in <-chan interface{}, fn func(interface{}) interface{}) <-chan interface{} {
	out := make(chan interface{}) //创建一个输出chan
	if in == nil {                // 异常检查
		close(out)
		return out
	}

	go func() { // 启动一个goroutine,实现map的主要逻辑
		defer close(out)
		for v := range in { // 从输入chan读取数据，执行业务操作，也就是map操作
			out <- fn(v)
		}
	}()

	return out
}

func reduce(in <-chan interface{}, fn func(r, v interface{}) interface{}) interface{} {
	if in == nil { // 异常检查
		return nil
	}

	out := <-in         // 先读取第一个元素
	for v := range in { // 实现reduce的主要逻辑
		out = fn(out, v)
	}

	return out
}
