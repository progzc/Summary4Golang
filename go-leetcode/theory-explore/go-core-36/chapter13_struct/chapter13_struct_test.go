package chapter13_struct

import (
	"fmt"
	"testing"
)

// TestStruct_1 结构体
// (1)在Go语言中，我们可以通过为一个类型编写名为String的方法，来自定义该类型的字符串表示形式，然后使用占位符%s就可以打印该字符串
// (2)方法隶属的类型其实并不局限于结构体类型，但必须是某个自定义的数据类型，并且不能是任何接口类型。
// (3)同一个方法集合中的方法不能出现重名;方法名称与该类型中任何字段的名称也不能重复
// (4)关于匿名字段
//	a.匿名字段的类型既是类型也是名称(在访问匿名字段的字段时可指定匿名字段名称或不指定)
//	b.嵌入字段的方法集合会被无条件地合并进被嵌入类型的方法集合
//	c.被嵌入类型的方法或字段都会“屏蔽”掉嵌入字段的同名方法或字段
//	d.可以通过链式的选择表达式来选择到被屏蔽的嵌入字段的字段或方法
//	e.针对多层嵌套,嵌入层级越深的字段或方法越可能被“屏蔽”
//	f.处于同一个层级的多个嵌入字段拥有同名的字段或方法，那么从被嵌入类型的值那里，选择此名称的时候就会引发一个编译错误
// (5)Q:Go语言是用嵌入字段实现了继承吗?
//	  A:Go语言中根本没有继承的概念，它所做的是通过嵌入字段的方式实现了类型之间的组合;
//	 	类型组合也是非侵入式的,它不会破坏类型的封装或加重类型之间的耦合;
//		组合要比继承更加简洁和清晰,Go语言可以轻而易举地通过嵌入多个字段来实现功能强大的类型,却不会有多重继承那样复杂的层次结构和可观的管理成本;
//		接口类型之间也可以组合;
// (6)Q:值方法和指针方法的区别?
//	  A: 详见TestStruct_2
//	 	a.值方法的接收者是该方法所属的那个类型值的一个副本,在该方法内对该副本的修改一般都不会体现在原值上，除非这个类型本身是某个引用类型（比如切片或字典）的别名类型;
//	      指针方法的接收者，是该方法所属的那个基本类型值的指针值的一个副本。我们在这样的方法内对该副本指向的值进行修改，却一定会体现在原值上;
//		b.一个自定义数据类型的方法集合中仅会包含它的所有值方法，而该类型的指针类型的方法集合却囊括了前者的所有方法，包括所有值方法和所有指针方法
//		c.一个指针类型实现了某某接口类型，但它的基本类型却不一定能够作为该接口的实现类型
// (7)Q:可以在结构体类型中嵌入某个类型的指针类型? 答案是可以
func TestStruct_1(t *testing.T) {
	// (1)在Go语言中，我们可以通过为一个类型编写名为String的方法，来自定义该类型的字符串表示形式。
	category := AnimalCategory{species: "cat"}
	fmt.Printf("The animal category: %s\n", category)

	//	b.嵌入字段的方法集合会被无条件地合并进被嵌入类型的方法集合
	animal := Animal{
		scientificName: "American Shorthair",
		// (7)Q:可以在结构体类型中嵌入某个类型的指针类型? 答案是可以
		AnimalCategory: &category,
	}
	fmt.Printf("The animal: %s\n", animal)

	//	c.被嵌入类型的方法或字段都会“屏蔽”掉嵌入字段的同名方法或字段
	cat := Cat{
		name:   "little pig",
		Animal: animal,
	}
	fmt.Printf("The cat: %s\n", cat)
}

// TestStruct_2
// 关于指针方法和值方法: https://blog.csdn.net/weixin_44676081/article/details/111309791
// (1) 值对象可以调用值方法和指针方法
// (2) 指针对象可以调用值方法和指针方法
// (3) 若调用一个接口里面的函数,结构体对象实现接口时的方法可能是指针方法也可以是值方法,那么需要注意:
//  a.值类型只能调用值方法
//  b.指针类型可以调用值方法和指针方法
// (4)尤其要注意: 一个指针类型实现了某某接口类型，但它的基本类型却不一定能够作为该接口的实现类型
func TestStruct_2(t *testing.T) {
	d1 := Dog{}
	d1.call()
	d1.phone() // 这里是语法糖, 相当于(&d1).phone()
	fmt.Println()

	d2 := &Dog{}
	d2.call() // 这里是语法糖, 相当于(*d1).call
	d2.phone()

	var jack Human
	jack = &Person{}
	jack.SayHello()
	jack.Cry()
	jack.Smile()

	//var tom Human
	//tom = Person{} // 值类型并没有实现Cry方法, 这里会编译报错
}

// --------------------------------------------------------------------------
type AnimalCategory struct {
	kingdom string // 界。
	phylum  string // 门。
	class   string // 纲。
	order   string // 目。
	family  string // 科。
	genus   string // 属。
	species string // 种。
}

func (ac AnimalCategory) String() string {
	return fmt.Sprintf("%s%s%s%s%s%s%s",
		ac.kingdom, ac.phylum, ac.class, ac.order,
		ac.family, ac.genus, ac.species)
}

// --------------------------------------------------------------------------
// (3)同一个方法集合中的方法不能出现重名;方法名称与该类型中任何字段的名称也不能重复
// 编译不通过
//func (ac AnimalCategory) class() string {
//	return ""
//}

// --------------------------------------------------------------------------
type Animal struct {
	scientificName  string // 学名。
	*AnimalCategory        // 动物基本分类。
}

func (a Animal) String() string {
	//	a.匿名字段的类型既是类型也是名称(在访问匿名字段的字段时可指定匿名字段名称或不指定)
	_ = a.AnimalCategory.order
	_ = a.order
	return fmt.Sprintf("%s (category: %s)",
		a.scientificName, a.AnimalCategory.order) // 匿名字段的类型既是类型也是名称
}

// --------------------------------------------------------------------------
type Cat struct {
	name string
	Animal
}

func (cat Cat) String() string {
	return fmt.Sprintf("%s (category: %s, name: %q)",
		cat.scientificName, cat.Animal.AnimalCategory, cat.name)
}

// ***************************************************************************
type Dog struct {
}

func (d Dog) call() {
	fmt.Println("call")
}

func (d *Dog) phone() {
	fmt.Println("phone")
}

// -----------------------------------------------------
type Human interface {
	SayHello()
	Cry()
	Smile()
}

type Person struct {
}

func (p Person) SayHello() {
	fmt.Println("SayHello")
}

func (p *Person) Cry() {
	fmt.Println("Cry")
}

func (p Person) Smile() {
	fmt.Println("Smile")
}
