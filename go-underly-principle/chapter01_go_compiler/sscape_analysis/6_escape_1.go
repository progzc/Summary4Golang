package sscape_analysis

var z *int

// go tool compile -m 6_escape_1.go
// go tool compile -m=2 6_escape_1.go
func escape() {
	// 这里会进行逃逸分析,a会被分配在堆上
	a := 1
	z = &a
}
