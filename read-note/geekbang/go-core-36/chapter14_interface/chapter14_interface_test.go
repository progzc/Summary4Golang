package chapter14_interface

import (
	"fmt"
	"reflect"
	"testing"
)

// TestInterface_1 第三点是常考的面试题
// (1)判定一个数据类型的某一个方法实现的就是某个接口类型中的某个方法的充要条件如下:
//	a.两个方法的签名需要完全一致
//	b.两个方法的名称要一模一样

// (2)在实现接口时,指针类型可以包含值方法,反之则不能

// (3)接口的底层是iface的实例,会包含两个指针:一个是指向类型信息的指针，另一个是指向动态值的指针
//	总之，接口变量被赋予动态值的时候，存储的是包含了这个动态值的副本的一个结构更加复杂的值。
//	基于此，即使我们把一个值为nil的某个实现类型的变量赋给了接口变量，后者的值也不可能是真正的nil。虽然这时它的动态值会为nil，但它的动态类型确是存在的
// Q:怎样才能让一个接口变量的值真正为nil呢?
// A:要么只声明它但不做初始化;要么直接把字面量nil赋给它
// Q:如果我们把一个值为nil的某个实现类型的变量赋给了接口变量，那么在这个接口变量上仍然可以调用该接口的方法吗?
// A:可以调用。但是请注意，这个被调用的方法在此时所持有的接收者的值是nil。因此，如果该方法引用了其接收者的某个字段,那么就会引发panic

// (4)接口间的组合: 推荐声明体量较小的接口,然后通过这种接口间的组合来扩展程序、增加程序的灵活性
//
//	例如:Go语言标准库代码包中,以io.ReadWriteCloser接口为例,它是由io.Reader、io.Writer和io.Closer这三个接口组成的。
func TestInterface_1(t *testing.T) {
	//(2)在实现接口时,指针类型可以包含值方法,反之则不能
	dog := Dog{"little pig"}
	_, ok := interface{}(dog).(Pet)
	fmt.Printf("Dog implements interface Pet: %v\n", ok) // false
	_, ok = interface{}(&dog).(Pet)
	fmt.Printf("*Dog implements interface Pet: %v\n", ok) // true
	fmt.Println()

	// pet的实际值叫 动态值, pet的实际值的实际类型叫 动态类型, pet的 静态类型 是Pet
	var pet Pet = &dog
	fmt.Printf("This pet is a %s, the name is %q.\n",
		pet.Category(), pet.Name())

	// --------------------------------------------------------------------------------
	// (3)接口的底层是iface的实例,会包含两个指针:一个是指向类型信息的指针，另一个是指向动态值的指针。
	// 总之，接口变量被赋予动态值的时候，存储的是包含了这个动态值的副本的一个结构更加复杂的值。
	// Q:怎样才能让一个接口变量的值真正为nil呢?
	// A:要么只声明它但不做初始化;要么直接把字面量nil赋给它
	var dog1 *Dog
	fmt.Println("The first dog is nil.")
	dog2 := dog1
	fmt.Println("The second dog is nil.")
	var pet2 Pet2 = dog2
	if pet2 == nil {
		fmt.Println("The pet is nil.")
	} else {
		fmt.Println("The pet is not nil.") // not nil
	}
	fmt.Printf("The type of pet is %T.\n", pet2)                          // *Dog
	fmt.Printf("The type of pet is %s.\n", reflect.TypeOf(pet2).String()) // *Dog
	fmt.Printf("The type of second dog is %T.\n", dog2)                   // *Dog
	fmt.Println()

	wrap := func(dog *Dog) Pet2 {
		if dog == nil {
			return nil
		}
		return dog
	}
	pet2 = wrap(dog2)
	if pet2 == nil {
		fmt.Println("The pet is nil.")
	} else {
		fmt.Println("The pet is not nil.")
	}

	// --------------------------------------------------------------------------------
	// (4)接口间的组合: 推荐声明体量较小的接口,然后通过这种接口间的组合来扩展程序、增加程序的灵活性
	type Animal interface {
		ScientificName() string
		Category() string
	}
	type Named interface {
		Name() string
	}
	type Pet interface {
		Animal
		Named
	}
}

// TestInterface_2 几个思考题
// (1)如果我们使用一个变量给另外一个变量赋值,那么真正赋给后者的,并不是前者持有的那个值,而是该值的一个副本
func TestInterface_2(t *testing.T) {
	// 思考题1
	// (1)如果我们使用一个变量给另外一个变量赋值,那么真正赋给后者的,并不是前者持有的那个值,而是该值的一个副本
	dog := Dog{"little pig"}
	fmt.Printf("The dog's name is %q.\n", dog.Name()) // little pig
	var pet Pet2 = dog                                // 赋值的是副本而已
	dog.SetName("monster")
	fmt.Printf("The dog's name is %q.\n", dog.Name())                             // monster
	fmt.Printf("This pet is a %s, the name is %q.\n", pet.Category(), pet.Name()) // little pig
	fmt.Println()

	// 思考题2
	// (1)如果我们使用一个变量给另外一个变量赋值,那么真正赋给后者的,并不是前者持有的那个值,而是该值的一个副本
	dog1 := Dog{"little pig"}
	fmt.Printf("The name of first dog is %q.\n", dog1.Name()) // little pig
	dog2 := dog1
	fmt.Printf("The name of second dog is %q.\n", dog2.Name()) // little pig
	dog1.name = "monster"
	fmt.Printf("The name of first dog is %q.\n", dog1.Name())  // monster
	fmt.Printf("The name of second dog is %q.\n", dog2.Name()) // little pig
	fmt.Println()

	// 思考题3
	dog = Dog{"little pig"}
	fmt.Printf("The dog's name is %q.\n", dog.Name())
	pet = &dog
	dog.SetName("monster")
	fmt.Printf("The dog's name is %q.\n", dog.Name())
	fmt.Printf("This pet is a %s, the name is %q.\n",
		pet.Category(), pet.Name())
}

type Pet interface {
	SetName(name string)
	Name() string
	Category() string
}

type Pet2 interface {
	Name() string
	Category() string
}

type Dog struct {
	name string // 名字
}

func (dog *Dog) SetName(name string) {
	dog.name = name
}

func (dog Dog) Name() string {
	return dog.name
}

func (dog Dog) Category() string {
	return "dog"
}
