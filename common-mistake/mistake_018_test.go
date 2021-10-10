package common_mistake

import (
	"fmt"
	"testing"
)

// string 与索引操作符
// (1) 对字符串用索引访问返回的不是字符，而是一个 byte 值。
// (2) 如果需要使用 for range 迭代访问字符串中的字符（unicode code point / rune），
// 标准库中有 "unicode/utf8" 包来做 UTF8 的相关解码编码。
// 另外 utf8string 也有像 func (s *String) At(i int) rune 等很方便的库函数。
func TestMistake_018(t *testing.T) {
	right018()
}

func right018() {
	x := "ascii"
	fmt.Println(x[0])        // 97
	fmt.Printf("%T\n", x[0]) // uint8
}
