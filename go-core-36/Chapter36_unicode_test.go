package go_core_36

import (
	"fmt"
	"testing"
)

func TestUnicode(t *testing.T) {
	str := "Go爱好者"
	fmt.Printf("The string: %q\n", str) // "Go爱好者"

	fmt.Printf("  => runes(char): %q\n", []rune(str)) // ['G' 'o' '爱' '好' '者']
	fmt.Printf("  => runes(char): %v\n", []rune(str)) // [71 111 29233 22909 32773]

	fmt.Printf("  => runes(hex): %x\n", []rune(str)) // [47 6f 7231 597d 8005]

	fmt.Printf("  => bytes(hex): [% x]\n", []byte(str)) // [47 6f e7 88 b1 e5 a5 bd e8 80 85]
	fmt.Printf("  => bytes(hex): [%x]\n", []byte(str))  // [476fe788b1e5a5bde88085]
}

func TestUnicode2(t *testing.T) {
	str := "Go爱好者"
	fmt.Println("len:", len(str)) // len: 11
	// 0: 'G' [47]
	// 1: 'o' [6f]
	// 2: '爱' [e7 88 b1]
	// 5: '好' [e5 a5 bd]
	// 8: '者' [e8 80 85]
	for i, c := range str {
		fmt.Printf("%d: %q [% x]\n", i, c, []byte(string(c)))
	}
}
