package common_mistake

import (
	"fmt"
	"testing"
)

// 运算符的优先级
// 除了位清除（bit clear）操作符，Go 也有很多和其他语言一样的位操作符，但优先级另当别论。
// Precedence    Operator
//    5             *  /  %  <<  >>  &  &^
//    4             +  -  |  ^
//    3             ==  !=  <  <=  >  >=
//    2             &&
//    1             ||
func TestMistake_029(t *testing.T) {
	wrong029()
	right029()
}

func wrong029() {
}

func right029() {
	fmt.Printf("0x2 & 0x2 + 0x4 -> %#x\n", 0x2&0x2+0x4) // & 优先 +
	//prints: 0x2 & 0x2 + 0x4 -> 0x6
	//Go:    (0x2 & 0x2) + 0x4
	//C++:    0x2 & (0x2 + 0x4) -> 0x2

	fmt.Printf("0x2 + 0x2 << 0x1 -> %#x\n", 0x2+0x2<<0x1) // << 优先 +
	//prints: 0x2 + 0x2 << 0x1 -> 0x6
	//Go:     0x2 + (0x2 << 0x1)
	//C++:   (0x2 + 0x2) << 0x1 -> 0x8

	fmt.Printf("0xf | 0x2 ^ 0x2 -> %#x\n", 0xf|0x2^0x2) // | 优先 ^
	//prints: 0xf | 0x2 ^ 0x2 -> 0xd
	//Go:    (0xf | 0x2) ^ 0x2
	//C++:    0xf | (0x2 ^ 0x2) -> 0xf
}
