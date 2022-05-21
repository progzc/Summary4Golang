package sscape_analysis

// go_knowledge tool compile -m 6_escape_2.go_knowledge
// go_knowledge tool compile -m=2 6_escape_2.go_knowledge
func f() int {
	// 这里会进行逃逸分析,a不会被分配在堆上
	a := 1
	c := &a
	return *c
}
