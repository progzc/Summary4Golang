package common_mistake

import (
	"testing"
)

// 对内建数据结构的操作并不是同步的
// (1) 尽管 Go 本身有大量的特性来支持并发，但并不保证并发的数据安全，用户需自己保证变量等数据以原子操作更新。
// (2) goroutine 和 channel 是进行原子操作的好方法，或使用 "sync" 包中的锁。
func TestMistake_023(t *testing.T) {
	wrong023()
	right023()
}

func wrong023() {

}

func right023() {
}
