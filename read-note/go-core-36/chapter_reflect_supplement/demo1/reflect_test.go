package demo1

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// https://segmentfault.com/a/1190000040639681
// TestReflect_1 获取类型名以及 kind
func TestReflect_1(t *testing.T) {
	// 声明一个空结构体
	type cat struct {
	}
	// 获取结构体实例的反射类型对象
	typeOfCat := reflect.TypeOf(cat{})
	// 显示反射类型对象的名称和种类
	fmt.Println(typeOfCat.Name(), typeOfCat.Kind()) // cat struct

	v := reflect.ValueOf(interface{}(cat{}))
	fmt.Println(v.Kind()) // struct
}

// TestReflect_2 获取成员反射信息
func TestReflect_2(t *testing.T) {
	// 声明一个空结构体
	type cat struct {
		Name string
		// 带有结构体tag的字段
		Type int `json:"type" id:"100"`
	}
	// 创建cat的实例
	ins := cat{Name: "mimi", Type: 1}
	// 获取结构体实例的反射类型对象
	typeOfCat := reflect.TypeOf(ins)
	// 遍历结构体所有成员
	for i := 0; i < typeOfCat.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldType := typeOfCat.Field(i)
		// 输出成员名和tag
		fmt.Printf("name: %v  tag: '%v'\n", fieldType.Name, fieldType.Tag)
	}
	// 通过字段名, 找到字段类型信息
	if catType, ok := typeOfCat.FieldByName("Type"); ok {
		// 从tag中取出需要的tag
		fmt.Println(catType.Tag.Get("json"), catType.Tag.Get("id"))
	}
}

// TestReflect_3 通过类型信息创建实例
func TestReflect_3(t *testing.T) {
	var a int
	// 取变量a的反射类型对象
	typeOfA := reflect.TypeOf(a)
	// 根据反射类型对象创建类型实例
	aIns := reflect.New(typeOfA)
	aIns.Elem().SetInt(14)
	// 输出：*int ptr 14
	fmt.Println(aIns.Type(), aIns.Kind(), aIns.Elem().Int())
}

// TestReflect_4 生成原始类型的对象
func TestReflect_4(t *testing.T) {
	// 声明整型变量a并赋初值
	var a = 1024
	// 获取变量a的反射值对象
	valueOfA := reflect.ValueOf(a)
	// Interface方法将值以interface{}类型返回，可以通过类型断言转换为指定类型
	var getA = valueOfA.Interface().(int)
	// 获取64位的值, 强制类型转换为int类型
	var getA2 = int(valueOfA.Int())
	fmt.Println(getA, getA2)
}

// TestReflect_5 操作结构体成员的值
func TestReflect_5(t *testing.T) {
	type dog struct {
		LegCount int
		age      int
	}
	// 获取dog实例地址的反射值对象
	valueOfDog := reflect.ValueOf(&dog{})
	// 取出dog实例地址的元素
	valueOfDog = valueOfDog.Elem()
	// 获取legCount字段的值
	vLegCount := valueOfDog.FieldByName("LegCount")
	vAge := valueOfDog.FieldByName("age")
	// 尝试设置legCount的值
	vLegCount.SetInt(4)
	// 这里会报错
	// 如果该对象不可寻址或者成员是私有的，则无法修改对象值
	vAge.SetInt(4)
	fmt.Println(vLegCount.Int())
}

// 普通函数
func add(a, b int) int {
	return a + b
}

// TestReflect_6 通过反射调用函数
func TestReflect_6(t *testing.T) {
	// 将函数包装为反射值对象
	funcValue := reflect.ValueOf(add)
	// 构造函数参数, 传入两个整型值
	paramList := []reflect.Value{reflect.ValueOf(10), reflect.ValueOf(20)}
	// 反射调用函数
	retList := funcValue.Call(paramList)
	// 获取第一个返回值, 取整数值
	fmt.Println(retList[0].Int())
}

type Cat struct {
	Name string
}

func (c Cat) Sleep() {
	fmt.Println("Sleep...")
}

func (c *Cat) Cry() {
	fmt.Println("Cry...")
}

// TestReflect_7 通过反射调用对象中的方法
func TestReflect_7(t *testing.T) {
	cat := Cat{}
	valueOf := reflect.ValueOf(&cat)
	showMethod1 := valueOf.MethodByName("Sleep")
	showMethod1.Call(nil)

	// 下面这句会报错
	// showMethod2 := valueOf.Elem().MethodByName("Cry")
	showMethod2 := valueOf.MethodByName("Cry")
	showMethod2.Call(nil)
}

// Map2Struct map转struct
func Map2Struct(m map[string]interface{}, obj interface{}) {
	value := reflect.ValueOf(obj)

	// obj 必须是指针且指针指向的必须是 struct
	if value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct {
		value = value.Elem()
		getMapName := func(key string) interface{} {
			for k, v := range m {
				if strings.EqualFold(k, key) {
					return v
				}
			}
			return nil
		}
		// 循环赋值
		for i := 0; i < value.NumField(); i++ {
			// 获取字段 type 对象
			field := value.Field(i)
			if !field.CanSet() {
				continue
			}
			// 获取字段名称
			fieldName := value.Type().Field(i).Name
			fmt.Println("fieldName -> ", fieldName)
			// 获取 map 中的对应的值
			fieldVal := getMapName(fieldName)
			if fieldVal != nil {
				field.Set(reflect.ValueOf(fieldVal))
			}
		}
	} else {
		panic("must prt")
	}
}

// Struct2Map struct 转 map
func Struct2Map(obj interface{}) map[string]interface{} {
	value := reflect.ValueOf(obj)

	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		panic("must prt")
	}
	value = value.Elem()
	t := value.Type()

	// 创建 map
	resultMap := make(map[string]interface{})

	// 循环获取字段名称以及对应的值
	for i := 0; i < value.NumField(); i++ {
		val := value.Field(i)
		typeName := t.Field(i)
		if !val.CanSet() {
			resultMap[typeName.Name] = reflect.New(typeName.Type).Elem().Interface()
			continue
		}
		resultMap[typeName.Name] = val.Interface()
	}

	return resultMap
}
