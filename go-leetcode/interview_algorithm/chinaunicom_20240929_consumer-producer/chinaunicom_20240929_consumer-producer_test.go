package chinaunicom_20240929_consumer_producer

import "testing"

// TestProducerConsumerModel
// 面试题3: 生产者与消费者（goroutine与channel）
// 1.有多个生产者不停生产整数1，有多个消费者同时进行消费。
// 2.所有的消费者的消费行为是：从channel中读取数据进行消费，并用一个公用的计数器进行累加。
// 3.某个消费者在做累加过程中，当计数器达到某数值时，通知所有生产者停止生产， 同时也通知其它消费者退出，然后自己也退出。
// 4.生产者一旦收到退出通知后，立即停止生产数据，并退出。
// 5.最后主协程等所有子协程全部退出后，主协程再退出。
func TestProducerConsumerModel(t *testing.T) {
	ProducerConsumerModel()
}
