package common_mistake

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

// 字符串的长度
// (1) Go 的内建函数 len() 返回的是字符串的 byte 数量
// (2) 如果要得到字符串的字符数，可使用 "unicode/utf8" 包中的 RuneCountInString(str string) (n int)
func TestMistake_020(t *testing.T) {
	right020_1()
	right020_2()
}

func right020_1() {
	char := "♥"
	fmt.Println(len(char)) // 3
}

func right020_2() {
	char := "é"
	fmt.Println(len(char)) // 3
	// 法文的é,实际上是两个 rune 的组合
	fmt.Println(utf8.RuneCountInString(char)) // 2
	fmt.Println("cafe\u0301")                 // café
}
