package chapter36_unicode

import (
	"fmt"
	"testing"
)

// TestUnicode_1
// (1)UTF-8:
//    a.“-”右边的整数的含义是，以多少个比特位作为一个编码单元。以 UTF-8 为例，它会以 8 个比特，也就是一个字节，作为一个编码单元。
//    b.并且，它与标准的 ASCII 编码是完全兼容的。也就是说，在[0x00, 0x7F]的范围内，这两种编码表示的字符都是相同的。这也是 UTF-8 编码格式的一个巨大优势。
//    c.Go语言的所有源代码，都必须按照Unicode编码规范中的UTF-8编码格式进行编码。
// (2)Q:一个string类型的值在底层是怎样被表达的?
//	  A:在底层，一个string类型的值是由一系列相对应的Unicode代码点的UTF-8编码值来表达的。
//      具体而言,在底层都会被转换为UTF-8编码值，而这些UTF-8编码值又会以字节序列的形式表达和存储。
//      因此，一个string类型的值在底层就是一个能够表达若干个UTF-8编码值的字节序列。
// (3)Q:使用带有range子句的for语句遍历字符串值的时候应该注意什么？
//	  A:带有range子句的for语句会先把被遍历的字符串值拆成一个字节序列，然后再试图找出这个字节序列中包含的每一个 UTF-8 编码值，或者说每一个 Unicode 字符。
//	    需要注意的是:
//	    	如果存在两个迭代变量，那么赋给第一个变量的值，就将会是当前字节序列中的某个 UTF-8 编码值的第一个字节所对应的那个索引值。
//			赋给第二个变量的值，则是这个 UTF-8 编码值代表的那个 Unicode 字符，其类型会是rune。
// (4)Q:判断一个 Unicode 字符是否为单字节字符通常有几种方式?
//	  A:unicode/utf8代码包中有几个可以做此判断的函数，比如：RuneLen函数、EncodeRune函数等。
//	    我们需要根据输入的不同来选择和使用它们。具体可以查看该代码包的文档。
func TestUnicode_1(t *testing.T) {
	// rune其实是int32的别名  type rune = int32
	str := "Go爱好者"
	fmt.Printf("The string: %q\n", str)                 // "Go爱好者"
	fmt.Printf("  => runes(char): %q\n", []rune(str))   // ['G' 'o' '爱' '好' '者']
	fmt.Printf("  => runes(hex): %x\n", []rune(str))    // [47 6f 7231 597d 8005]
	fmt.Printf("  => bytes(hex): [% x]\n", []byte(str)) // [47 6f e7 88 b1 e5 a5 bd e8 80 85]

	//0: 'G' [47]
	//1: 'o' [6f]
	//2: '爱' [e7 88 b1]
	//5: '好' [e5 a5 bd]
	//8: '者' [e8 80 85]
	for i, c := range str {
		fmt.Printf("%d: %q [% x]\n", i, c, []byte(string(c)))
	}
}
