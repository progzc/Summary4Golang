package common_mistake

import (
	"fmt"
	"testing"
)

// 按位取反
// (1) 很多编程语言使用 ~ 作为一元按位取反（NOT）操作符，Go 重用 ^ XOR 操作符来按位取反
// (2) ^ 也是按位异或（XOR）操作符
// (3) 一个操作符能重用两次，是因为一元的 NOT 操作 NOT 0x02，与二元的 XOR 操作 0x22 XOR 0xff 是一致的。
// (4) Go 也有特殊的操作符 AND NOT &^ 操作符，不同位才取1
func TestMistake_028(t *testing.T) {
	wrong028()
	right028_1()
	right028_2()
}

func wrong028() {
	//fmt.Println(~2)        // bitwise complement operator is ^
}

func right028_1() {
	var d uint8 = 2
	fmt.Printf("%08b\n", d)  // 00000010
	fmt.Printf("%08b\n", ^d) // 11111101
}

func right028_2() {
	var a uint8 = 0x82
	var b uint8 = 0x02
	fmt.Printf("%08b [A]\n", a) // 10000010 [A]
	fmt.Printf("%08b [B]\n", b) // 00000010 [B]

	fmt.Printf("%08b (NOT B)\n", ^b)                                 // 11111101 (NOT B)
	fmt.Printf("%08b ^ %08b = %08b [B XOR 0xff]\n", b, 0xff, b^0xff) // 00000010 ^ 11111111 = 11111101 [B XOR 0xff]

	fmt.Printf("%08b ^ %08b = %08b [A XOR B]\n", a, b, a^b)          // 10000010 ^ 00000010 = 10000000 [A XOR B]
	fmt.Printf("%08b & %08b = %08b [A AND B]\n", a, b, a&b)          // 10000010 & 00000010 = 00000010 [A AND B]
	fmt.Printf("%08b &^%08b = %08b [A 'AND NOT' B]\n", a, b, a&^b)   // 10000010 &^00000010 = 10000000 [A 'AND NOT' B]
	fmt.Printf("%08b&(^%08b)= %08b [A AND (NOT B)]\n", a, b, a&(^b)) // 10000010&(^00000010)= 10000000 [A AND (NOT B)]
}
