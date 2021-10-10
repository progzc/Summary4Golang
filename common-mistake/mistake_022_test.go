package common_mistake

import (
	"log"
	"testing"
)

// log.Fatal 和 log.Panic 不只是 log
// (1) log 标准库提供了不同的日志记录等级，与其他语言的日志库不同，
// Go 的 log 包在调用 Fatal*()、Panic*() 时能做更多日志外的事，如中断程序的执行等：
func TestMistake_022(t *testing.T) {
	right022()
}

func right022() {
	log.Fatal("Fatal level log: log entry") // 输出信息后，程序终止执行
	log.Println("Nomal level log: log entry")
}
