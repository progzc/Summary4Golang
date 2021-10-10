package common_mistake

import (
	"fmt"
	"testing"
)

// 按位取反
// 很多编程语言使用 ~ 作为一元按位取反（NOT）操作符，Go 重用 ^ XOR 操作符来按位取反
func TestMistake_028(t *testing.T) {
	wrong028()
	right028()
}

func wrong028() {
	//fmt.Println(~2)        // bitwise complement operator is ^
}

func right028() {
	var d uint8 = 2
	fmt.Printf("%08b\n", d)  // 00000010
	fmt.Printf("%08b\n", ^d) // 11111101
}
