package multi_producer_multi_consumer

import (
	"testing"
)

// 面试题1: 模拟生产者和消费者:
//	生产者1: 1,5,9...
//	生产者2: 2,6,10...
//	生产者3: 3,7,11...
//	生产者4: 4,8,12...
//	消费者1: 1,3...
//	消费者2: 2,4...
func Test_monitor(t *testing.T) {
	//方法一: 使用有缓冲通道
	//monitor(4, 2)
	//方法二: 使用无缓冲通道
	monitor_2(4, 2)
}

// 面试题2: 使用多个协程顺序打印1~10
func Test_SeqPrint(t *testing.T) {
	//方法一: 使用sync.Cond
	//seqPrint(3, 10)
	//方法二: 使用通道
	seqPrint_2(3, 10)
}
