package common_mistake

import (
	"testing"

	// 错误示例
	//"fmt"    // imported and not used: "fmt"
	//"log"    // imported and not used: "log"
	//"time"    // imported and not used: "time"

	// 正确示例
	_ "fmt"
	"log"
	"time"
)

// 未使用的 import
// 如果你 import 一个包，但包中的变量、函数、接口和结构体一个都没有用到的话，将编译失败。
// 可以使用 _ 下划线符号作为别名来忽略导入的包，从而避免编译错误，这只会执行 package 的 init()
func TestMistake_003(t *testing.T) {
	wrong()
	right()
}

func wrong() {
}

func right() {
	_ = log.Println
	_ = time.Now
}
