package chapter_uintptr_supplement

import (
	"fmt"
	"testing"
	"unsafe"
)

// 基本概念
// 普通指针：*类型代表普通指针，用户传递对象的地址，但是不不能进行指针的运算。
// unsafe.Pointer：通用的指针类型，可以转换成任意的指针类型。不能进行指针的运算，也不能读取存储的值。
//				   如果要读取内存存储的值，需要转化为普通指针，再取值。
// uintptr：可用于指针运算，并不是指针，可以将unsafe.Pointer指针转化为uintptr类型（是和当前指针相同的一个数字值），然后进行指针指针运算。
// 总结：unsafe.Pointer 是 uintptr 和 普通指针 之间的桥梁。下面这四句话来自于官网：
//		1) 任意类型的指针可以转换为一个Pointer类型值
//		2) 一个Pointer类型值可以转换为任意类型的指针
//		3) 一个uintptr类型值可以转换为一个Pointer类型值
//		4) 一个Pointer类型值可以转换为一个uintptr类型值

// TestUintptrPointer_1
// (1)使用unsafe.Pointer用于指针类型转换
//
//	Go是一门强类型的静态语言，类型一旦定义就不能改变了，类型检查会再运行前进行。
//	为了安全考虑，Go不允许两个指针类型进行转换。
func TestUintptrPointer_1(t *testing.T) {
	h := 1
	p := &h

	// 编译错误：Cannot convert an expression of the type '*int' to the type '*int64'
	//var l *int64 = (*int64)(p)

	// 使用unsafe.Pointer用于指针类型转换
	var l *int64 = (*int64)(unsafe.Pointer(p))
	_ = l
}

// TestUintptrPointer_2
// (2)使用uintptr进行指针类型计算
//
//	a.unsafe.Pointer可以用于指针类型转换，但是不能进行指针运算。
//	b.如果像访问特定的内存，可以使用uintptr：
//	  将指针转换位unsafe.Pointer,再将unsafe.Pointer转化为uintptr，然后进行偏移量计算，这样就可以访问到特定的内存。
func TestUintptrPointer_2(t *testing.T) {
	var tA = &A{}
	var a = (*int8)(unsafe.Pointer(tA))
	*a = 1
	// Offsetof(x.y) y字段相对于x起始地址的偏移量，包括可能的空洞。
	var b = (*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(tA)) + unsafe.Offsetof(tA.i))) // 这里有语法糖，相当于unsafe.Offsetof((*tA).i)
	*b = 2
	var c = (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(tA)) + unsafe.Offsetof(tA.j))) // 这里有语法糖，相当于unsafe.Offsetof((*tA).j)
	*c = 3
	fmt.Println(tA) // &{1 2 3}
}

// TestUintptrPointer_3
// (3)注意事项
//
//	a.不要试图引入一个uintptr的临时变量，这可能会破环代码的安全性。
//	  uintptr类型的临时变量只是一个数字，并不持有对象。垃圾回收器无法识别这是个指向变量t的指针。
//	  当变量t被移动后，临时变量的值就已经不是变量t的地址了。再使用临时变量就会出现不可预测的错误。
//	b.unsafe是不安全的，尽可能少用。
func TestUintptrPointer_3(t *testing.T) {
	var tA = &A{}
	// uintptr用在临时变量里，可能会遇到不可预测的错误
	tmp := uintptr(unsafe.Pointer(tA)) + unsafe.Offsetof(tA.i)
	var b = (*int32)(unsafe.Pointer(tmp))
	*b = 2
}

// TestUintptrPointer_4
// (4)普通指针、unsafe.Pointer、uintptr 打印处理都是同一个数。
func TestUintptrPointer_4(t *testing.T) {
	var tA = &A{}
	// tA addr: 0xc000094300; tA unsafe.Pointer: 824634327808; tA uintptr: 824634327808
	// 其中：十六进制0xc000094300 等价于 十进制824634327808
	fmt.Printf("tA addr: %p; tA unsafe.Pointer: %d; tA uintptr: %d\n", tA, unsafe.Pointer(tA), uintptr(unsafe.Pointer(tA)))
}

// TestUintptrPointer_5
// (5)三个方法：unsafe.Sizeof、unsafe.Alignof、unsafe.OffsetOf
func TestUintptrPointer_5(t *testing.T) {
	// 下面这种是误用：
	// 注意事项：
	// 	a.unsafe.Sizeof接受指针类型的话，返回永远是8，这是一种误用。要知道占用内存大小应该接收值类型。
	//	b.unsafe.Offsetof(tA.i)这里相当于语法糖：unsafe.Offsetof((*tA).i)。
	//	c.unsafe.Pointer必须接受指针类型。
	var tA = &A{}
	fmt.Println(unsafe.Sizeof(tA), unsafe.Alignof(tA), unsafe.Offsetof(tA.i)) // 8 8 4

	// unsafe.Sizeof：返回占用内存大小
	// unsafe.Alignof：返回对齐方大小
	// unsafe.OffsetOf：返回y字段相对于x类型起始地址的偏移量
	var tB = A{}
	fmt.Println(unsafe.Sizeof(tB), unsafe.Alignof(tB), unsafe.Offsetof(tB.i)) // 16 8 4
}

type A struct {
	k int8  //1字节
	i int32 //4字节
	j int64 //8字节
}
