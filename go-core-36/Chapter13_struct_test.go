package go_core_36

import (
	"fmt"
	"testing"
)

// AnimalCategory 代表动物分类学中的基本分类法。
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

func TestStruct(t *testing.T) {
	category := AnimalCategory{species: "cat"}
	// 我们可以通过为一个类型编写名为String的方法，来自定义该类型的字符串表示形式。
	// 这个String方法不需要任何参数声明，但需要有一个string类型的结果声明。
	fmt.Printf("The animal category: %s\n", category)
}

type Animal struct {
	scientificName string // 学名。
	// 1. 如果一个字段的声明中只有字段的类型名而没有字段的名称，那么它就是一个嵌入字段，也可以被称为匿名字段。嵌入字段的类型既是类型也是名称。
	AnimalCategory // 动物基本分类。
}

func (a Animal) Category() string {
	return a.AnimalCategory.String()
}

// TestStruct2 嵌入字段的屏蔽现象（将Animal的String方法先注释）
func TestStruct2(t *testing.T) {
	category := AnimalCategory{species: "cat"}
	animal := Animal{
		scientificName: "American Shorthair",
		AnimalCategory: category,
	}
	// 2. 嵌入字段的方法集合会被无条件地合并进被嵌入类型的方法集合中。
	fmt.Printf("The animal: %s\n", animal)
}

func (a Animal) String() string {
	return fmt.Sprintf("%s (category: %s)",
		a.scientificName, a.AnimalCategory)
}

// TestStruct3 嵌入字段的屏蔽现象（将Animal的String方法取消注释）
func TestStruct3(t *testing.T) {
	category := AnimalCategory{species: "cat"}
	animal := Animal{
		scientificName: "American Shorthair",
		AnimalCategory: category,
	}
	// 2. 嵌入字段的方法集合会被无条件地合并进被嵌入类型的方法集合中。
	fmt.Printf("The animal: %s\n", animal)
}

type Cat struct {
	name           string
	scientificName string // 学名。
	category       string // 动物学基本分类。
	Animal
}

// TestStruct4 嵌入字段的嵌套现象
func TestStruct4(t *testing.T) {
	// 示例1。
	category := AnimalCategory{species: "cat"}
	fmt.Printf("The animal category: %s\n", category)

	// 示例2。
	animal := Animal{
		scientificName: "American Shorthair",
		AnimalCategory: category,
	}
	fmt.Printf("The animal: %s\n", animal)

	// 示例3。
	cat := Cat{
		name:   "little pig",
		Animal: animal,
	}
	fmt.Printf("The cat: %s\n", cat)
}

func New(name, scientificName, category string) Cat {
	return Cat{
		name:           name,
		scientificName: scientificName,
		category:       category,
	}
}

func (cat *Cat) SetName(name string) {
	cat.name = name
}

func (cat Cat) SetNameOfCopy(name string) {
	cat.name = name
}

func (cat Cat) Name() string {
	return cat.name
}

func (cat Cat) ScientificName() string {
	return cat.scientificName
}

func (cat Cat) Category() string {
	return cat.category
}

func (cat Cat) String() string {
	return fmt.Sprintf("%s (category: %s, name: %q)",
		cat.scientificName, cat.category, cat.name)
}

// TestStruct5 值方法与指针方法
func TestStruct5(t *testing.T) {
	cat := New("little pig", "American Shorthair", "cat")
	cat.SetName("monster") // (&cat).SetName("monster")
	fmt.Printf("The cat: %s\n", cat)

	cat.SetNameOfCopy("little pig")
	fmt.Printf("The cat: %s\n", cat)

	type Pet interface {
		SetName(name string)
		Name() string
		Category() string
		ScientificName() string
	}

	_, ok := interface{}(cat).(Pet)
	fmt.Printf("Cat implements interface Pet: %v\n", ok)
	_, ok = interface{}(&cat).(Pet)
	fmt.Printf("*Cat implements interface Pet: %v\n", ok)
}
