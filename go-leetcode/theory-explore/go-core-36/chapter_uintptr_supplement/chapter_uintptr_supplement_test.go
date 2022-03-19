package chapter_uintptr_supplement

import (
	"fmt"
	"math"
	"testing"
	"unsafe"
)

// https://docs.hacknode.org/gopl-zh/ch13/ch13-01.html
// 关于内存对齐: https://go101.org/optimizations/0.3-memory-allocations.html
// 关于内存分配: https://go101.org/article/memory-layout.html

func TestUnsafePrint1_0(t *testing.T) {
	type Number struct {
		L bool
		A int8
		B int16
	}
	num := Number{}
	fmt.Println(unsafe.Sizeof(num))  // 4
	fmt.Println(unsafe.Alignof(num)) // 2
}

func TestUnsafePrint1_1(t *testing.T) {
	type Number struct {
		L bool
		D int64
		A int8
		C int32
		B int16
	}
	num := Number{}
	fmt.Println(unsafe.Sizeof(num))     // 32
	fmt.Println(unsafe.Alignof(num))    // 8
	fmt.Println(unsafe.Offsetof(num.D)) // 8
}

func TestUnsafePrint1_2(t *testing.T) {
	type Number struct {
		L bool
		A int8
		B int16
		C int32
		D int64
	}
	num := Number{}
	fmt.Println(unsafe.Sizeof(num))     // 16
	fmt.Println(unsafe.Alignof(num))    // 8
	fmt.Println(unsafe.Offsetof(num.A)) // 1
}

func TestUnsafePrint2_1(t *testing.T) {
	type School struct {
		Name   string // 校名
		Grade  uint32 // 年级
		Class  uint8  // 班级
		Number uint16 // 学号
		Slogan []byte // 学校标语
	}
	sch := &School{}
	fmt.Println(unsafe.Sizeof(sch))  // 8
	fmt.Println(unsafe.Alignof(sch)) // 8
}

func TestUnsafePrint2_2(t *testing.T) {
	type School struct {
		Name   string // 校名
		Grade  uint32 // 年级
		Class  uint8  // 班级
		Number uint16 // 学号
		Slogan []byte // 学校标语
	}
	sch := School{}
	fmt.Println(unsafe.Sizeof(sch))  // 48
	fmt.Println(unsafe.Alignof(sch)) // 8
}

func TestUnsafePrint2(t *testing.T) {
	type School struct {
		Name   string // 校名
		Grade  uint32 // 年级
		Class  uint8  // 班级
		Number uint16 // 学号
		Slogan []byte // 学校标语
	}

	type Student struct {
		Name string // 姓名
		Age  int8   // 年龄
		Sch  School // 学校
	}
	stu := Student{}
	fmt.Println(unsafe.Sizeof(stu))  // 72
	fmt.Println(unsafe.Alignof(stu)) // 8
}

func TestUnsafePrint2_3(t *testing.T) {
	type School struct {
		Class uint8  // 班级
		Grade uint64 // 年级
	}

	type Student struct {
		Name string // 姓名
		Age  int8   // 年龄
		Sch  School // 学校
	}
	stu := Student{}
	fmt.Println(unsafe.Sizeof(stu))             // 40
	fmt.Println(unsafe.Alignof(stu))            // 8
	fmt.Println(unsafe.Offsetof(stu.Age))       // 16
	fmt.Println(unsafe.Offsetof(stu.Sch))       // 24
	fmt.Println(unsafe.Offsetof(stu.Sch.Class)) // 0
}

// TestUnsafePrint3_1 一道面试题
func TestUnsafePrint3_1(t *testing.T) {
	type S struct {
		A uint32
		B uint64
		C uint64
		D uint64
		E struct{}
	}
	s := S{}
	fmt.Println(unsafe.Alignof(s))    // 8
	fmt.Println(unsafe.Offsetof(s.E)) // 32
	fmt.Println(unsafe.Sizeof(s.E))   // 0
	fmt.Println(unsafe.Sizeof(s))     // 32
}

func TestUnsafePrint3_2(t *testing.T) {
	type S struct {
		A uint32
		B uint64
		C uint64
		E struct{}
		D uint64
	}
	s := S{}
	fmt.Println(unsafe.Alignof(s))    // 8
	fmt.Println(unsafe.Offsetof(s.E)) // 24
	fmt.Println(unsafe.Sizeof(s.E))   // 0
	fmt.Println(unsafe.Sizeof(s))     // 32
}

func TestMintModify1_1(t *testing.T) {
	type Mint struct {
		a uint8
		b uint8
		c uint8
		d uint8
		e uint8
	}
	m := Mint{}
	*(*uint8)(unsafe.Pointer(&m)) = 100
	fmt.Printf("m.a = %v\n", m.a) // m.a = 100
	fmt.Printf("m.b = %v\n", m.b) // m.b = 0
	fmt.Printf("m.c = %v\n", m.c) // m.c = 0
	fmt.Printf("m.d = %v\n", m.d) // m.d = 0
	fmt.Printf("m.e = %v\n", m.e) // m.e = 0
}

func TestMintModify1_2(t *testing.T) {
	type Mint struct {
		a uint8
		b uint8
		c uint8
		d uint8
		e uint8
	}
	m := Mint{}
	*(*uint8)(unsafe.Pointer(&m)) = 100
	*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&m)) + unsafe.Offsetof(m.b))) = 101
	*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&m)) + unsafe.Offsetof(m.c))) = 102
	*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&m)) + 3)) = 103
	*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&m)) + 4)) = 104
	fmt.Printf("m.a = %v\n", m.a) // m.a = 100
	fmt.Printf("m.b = %v\n", m.b) // m.b = 101
	fmt.Printf("m.c = %v\n", m.c) // m.c = 102
	fmt.Printf("m.d = %v\n", m.d) // m.d = 103
	fmt.Printf("m.e = %v\n", m.e) // m.e = 104
}

func TestMintModify1_3(t *testing.T) {
	type Mint struct {
		a int8
		b int8
		c int8
		d int8
		e int8
	}
	m := Mint{}
	*(*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&m)) + 3)) = 128
	fmt.Printf("m.a = %v\n", m.a) // m.a = 0
	fmt.Printf("m.b = %v\n", m.b) // m.b = 0
	fmt.Printf("m.c = %v\n", m.c) // m.c = 0
	fmt.Printf("m.d = %v\n", m.d) // m.d = -128
	fmt.Printf("m.e = %v\n", m.e) // m.e = 0
}

func TestMintModify1_4(t *testing.T) {
	type Mint struct {
		a int8
		b int8
		c int8
		d int8
		e int8
	}
	m := Mint{}
	*(*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&m)) + 3)) = 129
	fmt.Printf("m.a = %v\n", m.a) // m.a = 0
	fmt.Printf("m.b = %v\n", m.b) // m.b = 0
	fmt.Printf("m.c = %v\n", m.c) // m.c = 0
	fmt.Printf("m.d = %v\n", m.d) // m.d = -127
	fmt.Printf("m.e = %v\n", m.e) // m.e = 0
}

func TestMintModify1_5(t *testing.T) {
	type Mint struct {
		a int8
		b int8
		c int8
		d int8
		e int8
	}
	m := Mint{}
	*(*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&m)) + 3)) = 256
	fmt.Printf("m.a = %v\n", m.a) // m.a = 0
	fmt.Printf("m.b = %v\n", m.b) // m.b = 0
	fmt.Printf("m.c = %v\n", m.c) // m.c = 0
	fmt.Printf("m.d = %v\n", m.d) // m.d = 0
	fmt.Printf("m.e = %v\n", m.e) // m.e = 1
}

func TestMintModify1_6(t *testing.T) {
	type Mint struct {
		a int8
		b int8
		c int8
		d int8
		e int8
	}
	m := Mint{}
	*(*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&m)) + 3)) = -1
	fmt.Printf("m.a = %v\n", m.a) // m.a = 0
	fmt.Printf("m.b = %v\n", m.b) // m.b = 0
	fmt.Printf("m.c = %v\n", m.c) // m.c = 0
	fmt.Printf("m.d = %v\n", m.d) // m.d = -1
	fmt.Printf("m.e = %v\n", m.e) // m.e = -1
}

func TestUnsafeString(t *testing.T) {
	var a string
	bytes := make([]byte, 0)
	bytes = append(bytes, 0x23, 0x24)
	x := (*[3]uintptr)(unsafe.Pointer(&bytes))
	*(*uintptr)(unsafe.Pointer(&a)) = x[0]
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&a)) + 8)) = x[1]
	fmt.Println(a)      // #$
	fmt.Println(len(a)) // 2
}

func TestFloat32(t *testing.T) {
	bits := math.Float32bits(6.9)
	fmt.Printf("%b\n", bits)
	// 0  10000001  10111001100110011001100 （人工）
	// 0  10000001  10111001100110011001101 （计算机）
}
