package common_mistake

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

// 字符串并不都是 UTF8 文本
// (1) string 的值不必是 UTF8 文本，可以包含任意的值。只有字符串是文字字面值时才是 UTF8 文本，字串可以通过转义来包含其他数据。
// (2) 判断字符串是否是 UTF8 文本，可使用 "unicode/utf8" 包中的 ValidString() 函数
func TestMistake_019(t *testing.T) {
	right019()
}

func right019() {
	str1 := "ABC"
	fmt.Println(utf8.ValidString(str1)) // true

	str2 := "A\xfeC"
	fmt.Println(utf8.ValidString(str2)) // false

	str3 := "A\\xfeC"
	fmt.Println(utf8.ValidString(str3)) // true    // 把转义字符转义成字面值
}
