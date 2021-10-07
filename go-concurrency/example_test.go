package go_concurrency

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type goo int
type bar int

func TestMap(t *testing.T) {
	m := make(map[interface{}]int)
	m[goo(1)] = 1
	m[bar(1)] = 2
	fmt.Printf("%v\n", m) // map[1:2 1:1]
}

type MyError struct {
	s    string
	code int
}

func NewMyError(s string, code int) MyError {
	return MyError{s, code}
}

func (me MyError) Error() string {
	return me.s
}

func (me MyError) Code() int {
	return me.code
}

func TestError(t *testing.T) {
	myErr := NewMyError("not found", 404)
	wrapp1Err := fmt.Errorf("this is a wrapping1 error:%w", myErr)
	wrapp2Err := fmt.Errorf("this is a wrapping2 error:%w", wrapp1Err)
	fmt.Println(errors.Unwrap(wrapp2Err))
	fmt.Println(errors.Is(wrapp2Err, myErr))
	if target, ok := wrapp2Err.(MyError); ok {
		fmt.Println("type assert success:", target.Error())
	}
	var target MyError
	if errors.As(wrapp2Err, &target) {
		fmt.Println("as success:", target.Code())
	}
}

func TestString(t *testing.T) {
	var test1 = strings.NewReplacer("a", "A", "a", "B")
	s1 := test1.Replace("abc")
	fmt.Println(s1) //Abc
}

func TestJson(t *testing.T) {
	// 使用注意事项：
	// (1)json 包解析的是一个 JSON 数据，而 JSON 数据既可以是对象（object），也可以是数组（array），同时也可以是字符串（string）、
	// 数值（number）、布尔值（boolean）以及空值（null）。而上述的两个函数，其实也是支持对这些类型值的解析的。比如下面段代码也是能用的
	var s string
	// 注意字符串中的双引号不能缺，如果仅仅是 `Hello, world`，则这不是一个合法的 JSON 序列，会返回错误。
	err := json.Unmarshal([]byte(`"Hello, world!"`), &s)
	if err != nil {
		t.Error(err)
	}
	t.Log(s)

	// (2)json 在解析时，如果遇到大小写问题，会尽可能地进行大小写转换。即便是一个 key 与结构体中的定义不同，但如果忽略大小写后是相同的，
	// 那么依然能够为字段赋值。比如下面的例子可以说明：
	cert := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err = json.Unmarshal([]byte(`{"UserName":"root","passWord":"123456"}`), &cert)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("username =%s", cert.Username)
		t.Logf("password =%s", cert.Password)
	}
}
