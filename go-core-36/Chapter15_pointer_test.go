package go_core_36

import (
	"testing"
	"unsafe"
)

type Dog2 struct {
	name string
}

func (dog *Dog2) SetName(name string) {
	dog.name = name
}

func New2(name string) Dog2 {
	return Dog2{name}
}

// TestPointer1 不可寻址的值在使用上的限制
// 满足以下三个条件之一即是不可寻址的：不可变的、临时结果、不安全的
func TestPointer(t *testing.T) {
	// 下面语句会编译报错：不能在New("little pig")的结果值上调用指针方法。
	// 本质是因为：不能取得New("little pig")的地址。
	//New2("little pig").SetName("monster")
}

// TestPointer2
func TestPointer2(t *testing.T) {
	dog := Dog{"little pig"}
	dogP := &dog
	dogPtr := uintptr(unsafe.Pointer(dogP))

	namePtr := dogPtr + unsafe.Offsetof(dogP.name)
	_ = (*string)(unsafe.Pointer(namePtr))
}
