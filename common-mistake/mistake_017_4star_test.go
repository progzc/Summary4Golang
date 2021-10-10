package common_mistake

import (
	"fmt"
	"testing"
)

// string 与 byte slice 之间的转换
// (1) 当进行 string 和 byte slice 相互转换时，参与转换的是拷贝的原始值。
// (2) Go 在 string 与 byte slice 相互转换上优化了两点，避免了额外的内存分配:
//     a. 在 map[string] 中查找 key 时，使用了对应的 []byte，避免做 m[string(key)] 的内存分配
//	   b. 使用 for range 迭代 string 转换为 []byte 的迭代：for i,v := range []byte(str) {...}
func TestMistake_017(t *testing.T) {
	right017()
}

func right017() {
	str := "我爱中国"

	// 以字节数组遍历
	for i := 0; i < len(str); i++ {
		fmt.Print(str[i], " ") // 230 136 145 231 136 177 228 184 173 229 155 189
	}
	fmt.Println()

	// 以字节数组遍历
	for _, s := range []byte(str) {
		fmt.Print(s, " ") // 230 136 145 231 136 177 228 184 173 229 155 189
	}
	fmt.Println()

	// 以Unicode字符遍历
	for _, s := range str {
		fmt.Print(s, " ") // 25105 29233 20013 22269
	}
	fmt.Println()

	// 以Unicode字符遍历
	for _, s := range []rune(str) {
		fmt.Print(s, " ") // 25105 29233 20013 22269
	}
	fmt.Println()
}
