package go_knowledge

import (
	"encoding/json"
	"fmt"
	"testing"
	"unsafe"
)

func TestQuestion_1_1(t *testing.T) {
	var a1 []int
	a2 := []int{}
	fmt.Println(a1, a2) // [] []
	b1, _ := json.Marshal(a1)
	b2, _ := json.Marshal(a2)
	fmt.Println(string(b1), string(b2)) // null []
}

type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

func TestQuestion_1_1_2(t *testing.T) {
	var a1 []int
	a2 := make([]int, 0) // 其slice.array的初始值其实是一个全局变量的值 var zerobase uintptr（位于runtime/malloc.go）
	a3 := []string{}     // 其slice.array的初始值其实是一个全局变量的值 var zerobase uintptr（位于runtime/malloc.go）
	p1 := (*slice)(unsafe.Pointer(&a1))
	p2 := (*slice)(unsafe.Pointer(&a2))
	p3 := (*slice)(unsafe.Pointer(&a3))
	// a1, array=0, len=0, cap=0
	fmt.Printf("a1, array=%x, len=%d, cap=%d\n", uintptr(p1.array), uintptr(p1.len), uintptr(p1.cap))
	// a2, array=c000051e80, len=0, cap=0
	fmt.Printf("a2, array=%x, len=%d, cap=%d\n", uintptr(p2.array), uintptr(p2.len), uintptr(p2.cap))
	// a3, array=c000051e80, len=0, cap=0
	fmt.Printf("a3, array=%x, len=%d, cap=%d\n", uintptr(p3.array), uintptr(p3.len), uintptr(p3.cap))
}

func TestQuestion_1_2(t *testing.T) {
	var a1 map[int]int
	a2 := map[int]int{}
	//a1[2] = 2 // 这里会panic
	a2[3] = 3
	delete(a1, 1)
	delete(a2, 1)
	fmt.Println(a1, a2) // map[] map[]
	b1, _ := json.Marshal(a1)
	b2, _ := json.Marshal(a2)
	fmt.Println(string(b1), string(b2)) // null {}
}
