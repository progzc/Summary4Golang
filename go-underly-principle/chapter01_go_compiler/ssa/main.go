package main

var d uint8

// 在编译时指定GOSSAFUNC=main可以查看SSA初始及其后续优化阶段生成的代码片段
// GOSSAFUNC=main GOOS=linux GOARCH=amd64 go tool compile main.go
func main() {
	var a uint8 = 1
	a = 2
	if true {
		a = 3
	}
	d = a
}
