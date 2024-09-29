package common_mistake

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// 关闭 HTTP 连接
// (1) 一些支持 HTTP1.1 或 HTTP1.0 配置了 connection: keep-alive 选项的服务器会保持一段时间的长连接
// (2) 但标准库 "net/http" 的连接默认只在服务器主动要求关闭时才断开，所以你的程序可能会消耗完 socket 描述符。关闭长连接的办法：
//
//	a. 直接设置请求变量的 Close 字段值为 true，每次请求结束后就会主动关闭连接。
//	b. 设置 Header 请求头部选项 Connection: close，然后服务器返回的响应头部也会有这个选项，此时 HTTP 标准库会主动断开连接。
//	c. 创建一个自定义配置的 HTTP transport 客户端，用来取消 HTTP 全局的复用连接
//
// (3) 根据场景决定是否选用长连接：
//
//	a. 若你的程序要向同一服务器发大量请求，使用默认的保持长连接。
//	b. 若你的程序要连接大量的服务器，且每台服务器只请求一两次，那收到请求后直接关闭连接。或增加最大文件打开数 fs.file-max 的值。
func TestMistake_037(t *testing.T) {
	wrong037()
	right037_1()
	right037_2()
}

func wrong037() {
}

// 主动关闭连接
func right037_1() {
	req, err := http.NewRequest("GET", "http://golang.org", nil)
	checkError(err)

	// a. 直接设置请求变量的 Close 字段值为 true，每次请求结束后就会主动关闭连接
	req.Close = true
	// b. 设置 Header 请求头部选项 Connection: close，然后服务器返回的响应头部也会有这个选项，此时 HTTP 标准库会主动断开连接。
	//req.Header.Add("Connection", "close")

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	checkError(err)

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	fmt.Println(string(body))
}

func right037_2() {
	// c. 创建一个自定义配置的 HTTP transport 客户端，用来取消 HTTP 全局的复用连接
	tr := http.Transport{
		DisableKeepAlives: true,
	}
	client := http.Client{
		Transport: &tr,
	}

	resp, err := client.Get("https://golang.google.cn/")
	if resp != nil {
		defer resp.Body.Close()
	}
	checkError(err)

	fmt.Println(resp.StatusCode) // 200

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	fmt.Println(len(string(body)))
}
