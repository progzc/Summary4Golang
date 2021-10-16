package inline_function

// 会自动进行函数内联
func max1(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 会取消指定函数的内联
//go:noinline
func max2(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 编译时取消所有的函数内联：go build -gcflags="l" xxxx.go
// 编译时禁止优化并取消所有的函数内联（-N参数代表禁止优化, -l参数代表禁止内联）：go build -gcflags="-N -l" xxxx.go
