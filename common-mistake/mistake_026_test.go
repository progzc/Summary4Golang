package common_mistake

import (
	"fmt"
	"testing"
)

// switch 中的 fallthrough 语句
// (1) switch 语句中的 case 代码块会默认带上 break，
//     但可以使用 fallthrough 来强制执行下一个 case 代码块。
// (2) 也可以改写 case 为多条件判断
func TestMistake_026(t *testing.T) {
	wrong026()
	right026_1()
	right026_2()
}

func wrong026() {
}

func right026_1() {
	isSpace := func(char byte) bool {
		switch char {
		case ' ': // 空格符会直接 break，返回 false // 和其他语言不一样
		// fallthrough    // 返回 true
		case '\t':
			return true
		}
		return false
	}
	fmt.Println(isSpace('\t')) // true
	fmt.Println(isSpace(' '))  // false
}

func right026_2() {
	isSpace := func(char byte) bool {
		switch char {
		case ' ', '\t':
			return true
		}
		return false
	}
	fmt.Println(isSpace('\t')) // true
	fmt.Println(isSpace(' '))  // true
}
