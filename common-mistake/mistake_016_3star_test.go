package common_mistake

import (
	"fmt"
	"testing"
)

// string 类型的值是常量，不可更改
// (1) 尝试使用索引遍历字符串，来更新字符串中的个别字符，是不允许的。
// (2) string 类型的值是只读的二进制 byte slice，如果真要修改字符串中的字符，
//     将 string 转为 []byte 修改后，再转为 string 即可
func TestMistake_016(t *testing.T) {
	wrong016()
	right016_1()
	right016_2()
}

func wrong016() {
	//x := "text"
	//x[0] = "T" // error: cannot assign to x[0]
	//fmt.Println(x)
}

// 下面的示例并不是更新字符串的正确姿势，因为一个 UTF8 编码的字符可能会占多个字节，
// 比如汉字就需要 3~4 个字节来存储，此时更新其中的一个字节是错误的。
func right016_1() {
	x := "text"
	xBytes := []byte(x)
	xBytes[0] = 'T' // 注意此时的 T 是 rune 类型
	x = string(xBytes)
	fmt.Println(x) // Text
}

// 更新字串的正确姿势：
// 将 string 转为 rune slice（此时 1 个 rune 可能占多个 byte），直接更新 rune 中的字符
func right016_2() {
	x := "text"
	xRunes := []rune(x)
	xRunes[0] = '我'
	x = string(xRunes)
	fmt.Println(x) // 我ext
}
