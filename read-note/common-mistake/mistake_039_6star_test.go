package common_mistake

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

// struct、array、slice 和 map 的值比较
// (1) 可以使用相等运算符 == 来比较结构体变量，前提是两个结构体的成员都是可比较的类型
// (2) 如果两个结构体中有任意成员是不可比较的，将会造成编译错误。注意数组成员只有在数组元素可比较时候才可比较。
//
//	无法比较的类型：切片/映射/函数/内嵌元素无法比较的类型
//
// (3) 可以使用 "reflect" 包的 DeepEqual()来比较那些无法使用 == 比较的变量
// (4) DeepEqual() 并不总适合于比较 slice
// (5) 如果要大小写不敏感来比较 byte 或 string 中的英文文本，可以使用 "bytes" 或 "strings" 包的 ToUpper() 和 ToLower() 函数。
//
//	比较其他语言的 byte 或 string，应使用 bytes.EqualFold() 和 strings.EqualFold()
//
// (6) 如果 byte slice 中含有验证用户身份的数据（密文哈希、token 等），不应再使用 reflect.DeepEqual()、bytes.Equal()、 bytes.Compare()。
//
//	这三个函数容易对程序造成 timing attacks，此时应使用 "crypto/subtle" 包中的 subtle.ConstantTimeCompare() 等函数
//
// (7) reflect.DeepEqual() 认为空 slice 与 nil slice 并不相等，但注意 byte.Equal() 会认为二者相等：
func TestMistake_039(t *testing.T) {
	wrong039()
	right039_1()
}

// (2) 如果两个结构体中有任意成员是不可比较的，将会造成编译错误。注意数组成员只有在数组元素可比较时候才可比较。
func wrong039() {
	type data struct {
		num    int
		checks [10]func() bool   // 无法比较
		doIt   func() bool       // 无法比较
		m      map[string]string // 无法比较
		bytes  []byte            // 无法比较
	}
	//v1 := data{}
	//v2 := data{}
	//
	//fmt.Println("v1 == v2: ", v1 == v2)
}

// (1) 可以使用相等运算符 == 来比较结构体变量，前提是两个结构体的成员都是可比较的类型
func right039_1() {
	type data struct {
		num     int
		fp      float32
		complex complex64
		str     string
		char    rune
		yes     bool
		events  <-chan string
		handler interface{}
		ref     *byte
		raw     [10]byte
	}
	v1 := data{}
	v2 := data{}
	fmt.Println("v1 == v2: ", v1 == v2) // true
}

// (3) 可以使用 "reflect" 包的 DeepEqual()来比较那些无法使用 == 比较的变量
func right039_2() {
	v1 := data{}
	v2 := data{}
	fmt.Println("v1 == v2: ", reflect.DeepEqual(v1, v2)) // true

	m1 := map[string]string{"one": "a", "two": "b"}
	m2 := map[string]string{"two": "b", "one": "a"}
	fmt.Println("v1 == v2: ", reflect.DeepEqual(m1, m2)) // true

	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	// 注意两个 slice 相等，值和顺序必须一致
	fmt.Println("v1 == v2: ", reflect.DeepEqual(s1, s2)) // true

	var b1 []byte = nil
	b2 := []byte{}
	fmt.Println("b1 == b2: ", reflect.DeepEqual(b1, b2)) // false
}

func right039_3() {
	var str = "one"
	var in interface{} = "one"
	fmt.Println("str == in: ", reflect.DeepEqual(str, in)) // true

	v1 := []string{"one", "two"}
	v2 := []string{"two", "one"}
	fmt.Println("v1 == v2: ", reflect.DeepEqual(v1, v2)) // false

	// (4) DeepEqual() 并不总适合于比较 slice
	data := map[string]interface{}{
		"code":  200,
		"value": []string{"one", "two"},
	}
	encoded, _ := json.Marshal(data)
	var decoded map[string]interface{}
	json.Unmarshal(encoded, &decoded)
	fmt.Println("data == decoded: ", reflect.DeepEqual(data, decoded)) // false
}

// (7) reflect.DeepEqual() 认为空 slice 与 nil slice 并不相等，但注意 byte.Equal() 会认为二者相等：
func right039_4() {
	var b1 []byte = nil
	b2 := []byte{}

	// b1 与 b2 长度相等、有相同的字节序
	// nil 与 slice 在字节上是相同的
	fmt.Println("b1 == b2: ", bytes.Equal(b1, b2)) // true
}
