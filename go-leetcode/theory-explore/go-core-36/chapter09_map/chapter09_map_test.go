package chapter09_map

import (
	"fmt"
	"testing"
)

type f int
type g int

// TestMap_1
// 相关知识: https://stackoverflow.com/questions/60341627/how-gos-map-hash-function-workes-so-different-type-with-same-value-result-in
//    比如写一行代码是new map代码时，编译器是知道你这个是啥类型的，把相应的type给你带上到map元信息里。 同时interface里具体类型的type信息也是编译器带进去的
//    这里需要重点掌握interface的底层原理：https://mp.weixin.qq.com/s/EbxkBokYBajkCR-MazL0ZA
func TestMap_1(t *testing.T) {
	f := f(1)
	g := g(1)
	m := make(map[interface{}]string)
	m[f] = "aaa"
	m[g] = "bbb"
	fmt.Println(m)    // map[1:aaa 1:bbb]
	fmt.Println(m[1]) // ""
	fmt.Println(m[f]) // aaa
	fmt.Println(m[g]) // bbb
}

// TestMap_2 map的使用
// (1)键类型的约束: 键要求支持判等操作,不能是 函数/字典/切片, 但可以是 通道
// (2)若键类型为interface{},则键值的实际类型也不能是 函数/字典/切片(否则编译时不会报错,但运行时会panic),所以尽量不要设置键类型为interface{}
// (3)如果键的类型是数组类型，那么还要确保该类型的元素类型不是函数类型、字典类型或切片类型 (无论嵌套多深)
// (4)如果键的类型是结构体类型，那么还要保证其中字段均可以进行判等操作 (无论嵌套多深)
// (5)求哈希和判等操作的速度越快，对应的类型就越适合作为键类型。优先级：基本类型 > 指针类型 > 数组类型 > 结构体类型 > 接口类型;
// 同种类型 类型的宽度（类型宽度即类型所占的字节数）越短,速度越快
// (6)当试图在一个值为nil的字典中添加 键-元素 对的时候，Go 语言的运行时系统就会立即抛出一个panic，读则不会
func TestMap_2(t *testing.T) {
	//(1)键类型的约束: 键要求支持判等操作,不能是 函数/字典/切片/通道
	//var badMap1_1 = map[func()]int{} // 键不能为函数(这里会引发编译错误)
	//_ = badMap1_1

	//var badMap1_2 = map[map[string]string]int{} // 键不能为字典(这里会引发编译错误)
	//_ = badMap1_2

	//var badMap1_3 = map[[]int]int{} // 键不能为切片(这里会引发编译错误)
	//_ = badMap1_3

	//var goodMap1 = map[chan int]int{} // 键可以为通道
	//_ = goodMap1

	// -------------------------------------------------------------------------------
	//(2)若键类型为interface{},则键值的实际类型也不能是 函数/字典/切片(否则编译时不会报错,但运行时会panic),所以尽量不要设置键类型为interface{}
	var badMap2 = map[interface{}]int{
		"1":      1,
		[]int{2}: 2, // 编译时不会报错,这里会引发panic。
		3:        3,
	}
	_ = badMap2

	// -------------------------------------------------------------------------------
	//(3)如果键的类型是数组类型，那么还要确保该类型的元素类型不是函数类型、字典类型或切片类型 (无论嵌套多深)
	//var badMap3_1 map[[1][]string]int // 这里会引发编译错误。
	//_ = badMap3_1
	//var badMap3_2 map[[1][2][3][]string]int // 这里会引发编译错误。
	//_ = badMap3_2

	// -------------------------------------------------------------------------------
	//(4)如果键的类型是结构体类型，那么还要保证其中字段均可以进行判等操作 (无论嵌套多深)
	//type BadKey1 struct {
	//	slice []string
	//}
	//var badMap4_1 map[BadKey1]int // 这里会引发编译错误。
	//_ = badMap4_1

	//type BadKey2 struct {
	//	field BadKey2Field1
	//}
	//type BadKey2Field1 struct {
	//	slice []string
	//}
	//var badMap4_2 map[BadKey2]int // 这里会引发编译错误。
	//_ = badMap4_2
}
