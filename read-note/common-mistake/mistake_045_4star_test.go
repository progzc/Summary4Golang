package common_mistake

import (
	"sync"
	"testing"
)

// 类型声明与方法
// (1) 从一个现有的非 interface 类型创建新类型时，并不会继承原有的方法
// (2) 如果你需要使用原类型的方法，可将原类型以匿名字段的形式嵌到你定义的新 struct 中
// (3) interface 类型声明也保留它的方法集
func TestMistake_045(t *testing.T) {
	wrong045()
	right045_1()
}

// (1) 从一个现有的非 interface 类型创建新类型时，并不会继承原有的方法
func wrong045() {
	// 定义 Mutex 的自定义类型
	type myMutex sync.Mutex

	//var mtx myMutex
	//mtx.Lock()
	//mtx.UnLock()
}

// (2) 如果你需要使用原类型的方法，可将原类型以匿名字段的形式嵌到你定义的新 struct 中
func right045_1() {
	// 类型以字段形式直接嵌入
	type myLocker struct {
		sync.Mutex
	}

	var locker myLocker
	locker.Lock()
	locker.Unlock()
}

// (3) interface 类型声明也保留它的方法集
func right045_2() {
	type myLocker sync.Locker

	var locker myLocker
	locker.Lock()
	locker.Unlock()
}
