package common_mistake

import (
	"fmt"
	"testing"
)

// range 迭代 string 得到的值
// (1) range 得到的索引是字符值（Unicode point / rune）第一个字节的位置,这个索引并不直接是字符在字符串中的位置
// (2) 注意一个字符可能占多个 rune，比如法文单词 café 中的 é
// (3) for range 迭代会尝试将 string 翻译为 UTF8 文本，对任何无效的码点都直接使用 0XFFFD rune（�）UNicode 替代字符来表示。
// (4) 如果 string 中有任何非 UTF8 的数据，应将 string 保存为 byte slice 再进行操作。
func TestMistake_024(t *testing.T) {
	wrong024()
	right024()
}

func wrong024() {

}

func right024() {
	data := "A\xfe\x02\xff\x04"
	for _, v := range data {
		fmt.Printf("%#x ", v) // 0x41 0xfffd 0x2 0xfffd 0x4    // 错误
	}
	fmt.Println()

	for _, v := range []byte(data) {
		fmt.Printf("%#x ", v) // 0x41 0xfe 0x2 0xff 0x4    // 正确
	}
	fmt.Println()
}
