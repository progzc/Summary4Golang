package demo2

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflect_1(t *testing.T) {
	var a float32 = 3.14
	reflectType(a) // type:float32 kind:float32
	var b int32 = 100
	reflectType(b) // type:int32 kind:int32
	var c = "小马"
	reflectType(c) // type:string kind:string
}

type myInt int64

//反射
func reflectType(x interface{}) {
	t := reflect.TypeOf(x)
	//打印变量名和底层类型
	fmt.Printf("type:%v kind:%v\n", t.Name(), t.Kind())
}

func TestReflect_2(t *testing.T) {
	var a *float32 // 指针
	var b myInt    // 自定义类型
	var c rune     // 类型别名
	reflectType(a) // type: kind:ptr
	reflectType(b) // type:myInt kind:int64
	reflectType(c) // type:int32 kind:int32

	type person struct {
		name string
		age  int
	}
	var d = person{
		name: "完满主任",
		age:  18,
	}
	type class struct {
		student int
	}
	var e = class{88}
	reflectType(d) // type:person kind:struct
	reflectType(e) // type:class kind:struct
}

//通过reflect.ValueOf反射获取值
func reflectValue(x interface{}) {
	v := reflect.ValueOf(x)
	k := v.Kind()
	//类型判断
	switch k {
	case reflect.Int64, reflect.Int:
		//从反射中获取整型的原始值，然后强转
		fmt.Printf("type is int64, value is %d\n", v.Int())
	case reflect.Float32:
		//从反射中获取浮点型的原始值，然后强转
		fmt.Printf("type is float32, value is %f\n", float32(v.Float()))
	case reflect.Float64:
		//从反射中获取浮点型的原始值，然后强转
		fmt.Printf("type is float64, value is %f\n", float64(v.Float()))
	case reflect.Bool:
		//从反射中获取布尔型的原始值，然后强转
		fmt.Printf("type is bool, value is %t\n", bool(v.Bool()))
	default:
		fmt.Println("nil")
	}
}

func TestReflect_3(t *testing.T) {
	var a float32 = 3.14
	var b int64 = 100
	var e = false
	reflectValue(a) // type is float32, value is 3.140000
	reflectValue(b) // type is int64, value is 100
	reflectValue(e) // type is bool, value is false
	// 将int类型的原始值转换为reflect.Value类型
	c := reflect.ValueOf(10)
	fmt.Printf("type c: %T\n", c) // type c: reflect.Value
	reflectValue(c)               // nil
}

func reflectSetValue1(x interface{}) {
	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Int64 {
		v.SetInt(250) //修改副本，会抛出panic异常
	}
}

func reflectSetValue2(x interface{}) {
	v := reflect.ValueOf(x)
	//使用Elem()方法修改变量的值
	if v.Elem().Kind() == reflect.Int64 {
		v.Elem().SetInt(251)
	}
}

// TestReflect_4 通过反射修改变量值
// 想要通过反射修改变量的值，就要知道参数在函数之间传递是值类型还是引用类型，必须是传递变量的内存地址才能修改变量值。
// 反射包reflect提供了一个Elem()方法返回指针对应的值。
func TestReflect_4(t *testing.T) {
	var a int64
	//reflectSetValue1(a)
	//panic: reflect: reflect.flag.mustBeAssignable using unaddressable value
	reflectSetValue2(&a)
	fmt.Printf("type is int64, value is %d\n", a)
}

type student struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

// TestReflect_5 结构体反射
// 当我们使用反射得到一个结构体之后可以通过索引依次获取其字段信息，也可以通过字段名去获取指定的字段信息。
func TestReflect_5(t *testing.T) {
	//初始化
	stu1 := student{
		Name:  "小马",
		Score: 99,
	}
	//返回类型
	t1 := reflect.TypeOf(stu1)
	fmt.Printf("name:%s kind:%v\n", t1.Name(), t1.Kind()) // name:student kind:struct
	// 通过for循环遍历结构体的所有字段信息
	for i := 0; i < t1.NumField(); i++ {
		field := t1.Field(i)
		fmt.Printf("name:%s index:%d type:%v json tag:%v\n", field.Name, field.Index, field.Type, field.Tag.Get("json"))
	}
	//通过字段名返回指定结构体的信息
	if scoreField, ok := t1.FieldByName("Score"); ok {
		fmt.Printf("name:%s index:%d type:%v json tag:%v\n", scoreField.Name, scoreField.Index, scoreField.Type, scoreField.Tag.Get("json"))
	}
}
