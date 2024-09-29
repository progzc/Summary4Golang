package common_mistake

import (
	"fmt"
	"runtime"
	"testing"
)

// GOMAXPROCS、Concurrency（并发）and Parallelism（并行）
// (1) Go 1.4 及以下版本，程序只会使用 1 个执行上下文 / OS 线程，即任何时间都最多只有 1 个 goroutine 在执行。
// (2) Go 1.5 版本将可执行上下文的数量设置为 runtime.NumCPU() 返回的逻辑 CPU 核心数，这个数与系统实际总的 CPU 逻辑核心数是否一致，
//
//	取决于你的 CPU 分配给程序的核心数，可以使用 GOMAXPROCS 环境变量或者动态的使用 runtime.GOMAXPROCS() 来调整。
//	误区：GOMAXPROCS 表示执行 goroutine 的 CPU 核心数
//	事实：GOMAXPROCS 的值是可以超过 CPU 的实际数量的，在 1.5 中最大为 256
func TestMistake_056(t *testing.T) {
	wrong056()
	right056()
}

func wrong056() {
}

func right056() {
	fmt.Println(runtime.GOMAXPROCS(-1)) // 12
	fmt.Println(runtime.NumCPU())       // 12
	runtime.GOMAXPROCS(20)
	fmt.Println(runtime.GOMAXPROCS(-1)) // 20
	runtime.GOMAXPROCS(300)
	fmt.Println(runtime.GOMAXPROCS(-1)) // Go 1.14 // 300
}
