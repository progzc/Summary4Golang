package sscape_analysis

var o *int

// go tool compile -m 6_escape_3.go
// go tool compile -m=2 6_escape_3.go
func f2() {
	// 这里会进行逃逸分析,new(int)会被分配在堆上
	l := new(int)
	*l = 42
	m := &l
	n := &m
	o = **n
}
