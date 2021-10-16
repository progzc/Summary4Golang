package sscape_analysis

// go tool compile -m 6_escape_2.go
// go tool compile -m=2 6_escape_2.go
func f() int {
	// 这里会进行逃逸分析,a不会被分配在堆上
	a := 1
	c := &a
	return *c
}
