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
	right016()
}

func wrong016() {
	//x := "text"
	//x[0] = "T" // error: cannot assign to x[0]
	//fmt.Println(x)
}

func right016() {
	x := "text"
	xBytes := []byte(x)
	xBytes[0] = 'T' // 注意此时的 T 是 rune 类型
	x = string(xBytes)
	fmt.Println(x) // Text
}
