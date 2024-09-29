package common_mistake

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

// 关闭 HTTP 的响应体
// (1) 使用 HTTP 标准库发起请求、获取响应时，即使你不从响应中读取任何数据或响应为空，都需要手动关闭响应体
//
//	常见错误：a. 忘记手动关闭  b. 或者写在了错误的位置
//
// (2) resp.Body.Close() 早先版本的实现是读取响应体的数据之后丢弃，保证了 keep-alive 的 HTTP 连接能重用处理不止一个请求。
//
//	但 Go 的最新版本将读取并丢弃数据的任务交给了用户，如果你不处理，HTTP 连接可能会直接关闭而非重用
//	两种处理方法：
//	a. 如果程序大量重用 HTTP 长连接，你可能要在处理响应的逻辑代码中加入
//	   _, err = io.Copy(ioutil.Discard, resp.Body)    // 手动丢弃读取完毕的数据
//	b. 如果你需要完整读取响应，上边的代码是需要写的。比如在解码 API 的 JSON 响应数据：
//	   json.NewDecoder(resp.Body).Decode(&data)
func TestMistake_036(t *testing.T) {
	wrong036_1()
	wrong036_2()
	right036_1()
}

// 下面代码能正确发起请求，但是一旦请求失败，变量 resp 值为 nil，造成 panic
// panic: runtime error: invalid memory address or nil pointer dereference
func wrong036_1() {
	resp, err := http.Get("https://api.ipify.org?format=json")
	defer resp.Body.Close() // resp 可能为 nil，不能读取 Body
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	fmt.Println(string(body))
}

// 下面的代码在绝大多数的情况下是正确的，但是在以下情况下是错误的：
// (1) 绝大多数请求失败的情况下，resp 的值为 nil 且 err 为 non-nil。
//
//	但如果你得到的是重定向错误，那它俩的值都是 non-nil，最后依旧可能发生内存泄露。
func wrong036_2() {
	resp, err := http.Get("https://api.ipify.org?format=json")
	checkError(err)

	defer resp.Body.Close() // 绝大多数情况下的正确关闭方式
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	fmt.Println(string(body))
}

// 正确的做法：可以直接在处理 HTTP 响应错误的代码块中，直接关闭非 nil 的响应体
func right036_1() {
	resp, err := http.Get("http://www.baidu.com")

	// 关闭 resp.Body 的正确姿势
	if resp != nil {
		defer resp.Body.Close()
	}
	checkError(err)

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	fmt.Println(string(body))
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
