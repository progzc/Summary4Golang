package go_core_36

import (
	"fmt"
	"testing"
)

// TestMap 字典的基本使用
func TestMap(t *testing.T) {
	aMap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	k := "two"
	v, ok := aMap[k]
	if ok {
		fmt.Printf("The element of key %q: %d\n", k, v)
	} else {
		fmt.Println("Not found!")
	}
}

// TestMap2 字典的键类型
func TestMap2(t *testing.T) {
	// 1. 字典的键不能是函数类型、字典类型和切片类型
	// 2. 如果键的类型是接口类型的,那么键值的实际类型也不能是上述三种类型(编译通过,但是运行时会报错),最好不要把字典的键类型设定为任何接口类型
	// 3. 如果键的类型是数组类型,那么还要确保该类型的元素类型不是函数类型、字典类型或切片类型
	var badMap2 = map[interface{}]int{
		"1":      1,
		[]int{2}: 2, // 这里会引发panic。
		3:        3,
	}
	fmt.Println(badMap2)
}

// TestMap3 字典的键类型
func TestMap3(t *testing.T) {
	// 示例1: 键类型不能为切片
	//var badMap1 = map[[]int]int{} // 这里会引发编译错误。
	//_ = badMap1

	// 示例2： 键类型不建议为接口类型,使用不当会引发panic
	var badMap2 = map[interface{}]int{
		"1":      1,
		[]int{2}: 2, // 这里会引发panic。
		3:        3,
	}
	_ = badMap2

	// 示例3: 键类型为数组的话,数组的值不能是切片、字典、函数类型
	//var badMap3 map[[1][]string]int // 这里会引发编译错误。
	//_ = badMap3

	// 示例4： 键类型为结构体的话,结构体中不能包含切片、字典、函数类型（无论藏的多深）
	//type BadKey1 struct {
	//	slice []string
	//}
	//var badMap4 map[BadKey1]int // 这里会引发编译错误。
	//_ = badMap4

	// 示例5：键类型为结构体的话,结构体中不能包含切片、字典、函数类型（无论藏的多深）
	//var badMap5 map[[1][2][3][]string]int // 这里会引发编译错误。
	//_ = badMap5

	// 示例6：键类型为结构体的话,结构体中不能包含切片、字典、函数类型（无论藏的多深）
	//type BadKey2Field1 struct {
	//	slice []string
	//}
	//type BadKey2 struct {
	//	field BadKey2Field1
	//}
	//var badMap6 map[BadKey2]int // 这里会引发编译错误。
	//_ = badMap6
}

// TestMap4 map的注意事项
func TestMap4(t *testing.T) {
	// map的声明: 此时map的值为nil
	var map1 map[string]int
	fmt.Printf("%v\n", map1)
	// map的声明: 此时map的值不为nil
	map2 := make(map[string]int)
	fmt.Printf("%v\n", map2)

	// 除了添加键-元素对，我们在一个值为nil的字典上做任何操作都不会引起错误（即: 在值为nil的map上添加键值对时会引发panic）
	var m map[string]int
	//m := make(map[string]int) // 将上面声明换成这句就不会出现panic
	key := "two"
	elem, ok := m["two"]
	fmt.Printf("The element paired with key %q in nil map: %d (%v)\n",
		key, elem, ok)
	fmt.Printf("The length of nil map: %d\n",
		len(m))
	fmt.Printf("Delete the key-element pair by key %q...\n",
		key)
	delete(m, key)
	elem = 2
	fmt.Println("Add a key-element pair to a nil map...")
	m["two"] = elem // 这里会引发panic。
}
