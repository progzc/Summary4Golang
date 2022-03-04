package chapter15_pointer

import (
	"fmt"
	"testing"
	"unsafe"
)

// TestFuncMethod_1
// (1)unsafe.Pointer可以表示任何指向可寻址的值的指针;unsafe.Pointer是指针值和uintptr值之间的桥梁
//	那么,怎么通过unsafe.Pointer操纵可寻址的值?
//	a.一个指针值（比如*Dog类型的值）可以被转换为一个unsafe.Pointer类型的值，反之亦然
//	b.一个uintptr类型的值也可以被转换为一个unsafe.Pointer类型的值，反之亦然
//	c.一个指针值无法被直接转换成一个uintptr类型的值，反过来也是如此
// 	所以，对于指针值和uintptr类型值之间的转换，必须使用unsafe.Pointer类型的值作为中转
//	常用的方法是:
//	a.unsafe.Offsetof: 用于获取 结构体值 和 其中某个字段 在内存中的起始存储地址之间的偏移量,以字节为单位

// (2)Go语言中的下列值是不可寻址的:
//	a.常量的值---->(1)不可变
//	b.基本类型值的字面量---->(1)不可变
//	c.算术操作的结果值---->(2)临时结果
//	d.对各种字面量的索引表达式和切片表达式的结果值【不过有一个例外，对切片字面量的索引结果值却是可寻址的】---->(2)临时结果
//	e.对字符串变量的索引表达式和切片表达式的结果值---->(1)不可变
//	f.对字典变量的索引表达式的结果值---->(3)不安全的
//	g.函数字面量和方法字面量，以及对它们的调用表达式的结果值---->(1)不可变+(2)临时结果
//	h.结构体字面量的字段值，也就是对结构体字面量的选择表达式的结果值---->(2)临时结果
//	i.类型转换表达式的结果值---->(2)临时结果
//	j.类型断言表达式的结果值---->(2)临时结果
//	l.接收表达式的结果值---->(2)临时结果
//	总结特点,分类: (1)不可变	(2)临时结果	(3)不安全的

// (3)不可寻址的值在使用上有哪些限制?
//	a.无法使用取址操作符&获取它们的指针(否则编译器会直接报错)；额外需要注意的是：++或--的左边表达式的值必须是可寻址的
//	b.在赋值语句中，赋值操作符左边的表达式的结果值必须可寻址的，但是对字典的索引结果值也是可以的
//	c.在带有range子句的for语句中，在range关键字左边的表达式的结果值也都必须是可寻址的，不过对字典的索引结果值同样可以被用在这里

// (4) Q:引用类型的值的指针值是有意义的吗？如果没有意义，为什么？如果有意义，意义在哪里？
//	A:从存储和传递的角度看，没有意义。因为引用类型的值已经相当于指向某个底层数据结构的指针了。当然，引用类型的值不只是指针那么简单。
func TestFuncMethod_1(t *testing.T) {
	// (2)Go语言中的下列值是不可寻址的:
	const num = 123
	//_ = &num // a.常量不可寻址
	//_ = &(123) // b.基本类型值的字面量不可寻址

	var str = "abc"
	_ = str
	//_ = &(str[0]) // e.对字符串变量的索引结果值不可寻址
	//_ = &(str[0:2]) // e.对字符串变量的切片结果值不可寻址
	str2 := str[0]
	_ = &str2 // 但这样的寻址就是合法的

	//_ = &(123 + 456) // c.算术操作的结果值不可寻址
	num2 := 456
	_ = num2
	//_ = &(num + num2) // c.算术操作的结果值不可寻址

	//_ = &([3]int{1, 2, 3}[0]) // d.对数组字面量的索引结果值不可寻址
	//_ = &([3]int{1, 2, 3}[0:2]) // d.对数组字面量的切片结果值不可寻址
	_ = &([]int{1, 2, 3}[0]) // d.对切片字面量的索引结果值却是可寻址的
	//_ = &([]int{1, 2, 3}[0:2]) // d.对切片字面量的切片结果值不可寻址

	//_ = &(map[int]string{1: "a"}[0]) // f.对字典字面量的索引结果值不可寻址

	var map1 = map[int]string{1: "a", 2: "b", 3: "c"}
	_ = map1
	//_ = &(map1[2]) // f.对字典变量的索引结果值不可寻址

	//_ = &(func(x, y int) int {
	//	return x + y
	//}) // g.字面量代表的函数不可寻址
	//_ = &(fmt.Sprintf) // g.标识符代表的函数不可寻址
	//_ = &(fmt.Sprintln("abc")) // g.对函数的调用结果值不可寻址

	dog := Dog{"little pig"}
	_ = dog
	//_ = &(dog.Name) // g.标识符代表的函数不可寻址
	//_ = &(dog.Name()) // g.对方法的调用结果值不可寻址

	//_ = &(Dog{"little pig"}.name) // h.结构体字面量的字段不可寻址

	//_ = &(interface{}(dog)) // i.类型转换表达式的结果值不可寻址
	dog1 := interface{}(dog)
	_ = dog1
	//_ = &(dogI.(Named)) // j.类型断言表达式的结果值不可寻址
	named := dog1.(Named)
	_ = named
	//_ = &(named.(Dog)) // j.类型断言表达式的结果值不可寻址

	var chan1 = make(chan int, 1)
	chan1 <- 1
	//_ = &(<-chan1) // l.接收表达式的结果值不可寻址

	// --------------------------------------------------------------
	// (3)不可寻址的值在使用上有哪些限制?
	//	a.无法使用取址操作符&获取它们的指针(否则编译器会直接报错)；额外需要注意的是：++或--的左边表达式的值必须是可寻址的
	//	b.在赋值语句中，赋值操作符左边的表达式的结果值必须可寻址的，但是对字典的索引结果值也是可以的
	//	c.在带有range子句的for语句中，在range关键字左边的表达式的结果值也都必须是可寻址的，不过对字典的索引结果值同样可以被用在这里
	//不能调用不可寻址的值的指针方法(根本原因是链式调用产生的New("little pig")是不可寻址的,不能采用语法糖施加&)
	//下边这行链式调用会让编译器报告两个错误，一个是果，即：不能在New("little pig")的结果值上调用指针方法。一个是因，即：不能取得New("little pig")的地址
	//New("little pig").SetName("monster")
	d := New("little pig")
	d.SetName("monster") // 这样实际会有语法糖 (&d).SetName("monster")

	map[string]int{"the": 0, "word": 0, "counter": 0}["word"]++
	map2 := map[string]int{"the": 0, "word": 0, "counter": 0}
	map2["word"]++

	// --------------------------------------------------------------
	// (1)unsafe.Pointer可以表示任何指向可寻址的值的指针;unsafe.Pointer是指针值和uintptr值之间的桥梁
	//	那么,怎么通过unsafe.Pointer操纵可寻址的值?
	//	a.一个指针值（比如*Dog类型的值）可以被转换为一个unsafe.Pointer类型的值，反之亦然
	//	b.一个uintptr类型的值也可以被转换为一个unsafe.Pointer类型的值，反之亦然
	//	c.一个指针值无法被直接转换成一个uintptr类型的值，反过来也是如此
	// 	所以，对于指针值和uintptr类型值之间的转换，必须使用unsafe.Pointer类型的值作为中转
	//	常用的方法是:
	//	a.unsafe.Offsetof: 用于获取 结构体值 和 其中某个字段 在内存中的起始存储地址之间的偏移量,以字节为单位
	dog2 := Dog{"little pig"}
	dogP := &dog2
	dogPtr := uintptr(unsafe.Pointer(dogP))

	nameP := (*string)(unsafe.Pointer(dogPtr + unsafe.Offsetof(dogP.name)))
	fmt.Printf("nameP == &(dogP.name)? %v\n", nameP == &(dogP.name)) // true
	fmt.Printf("The name of dog is %q.\n", *nameP)                   // little pig

	*nameP = "monster"
	fmt.Printf("The name of dog is %q.\n", dogP.name)
	fmt.Println()
}

type Named interface {
	// Name 用于获取名字。
	Name() string
}

type Dog struct {
	name string
}

func New(name string) Dog {
	return Dog{name}
}

func (dog *Dog) SetName(name string) {
	dog.name = name
}

func (dog Dog) Name() string {
	return dog.name
}
