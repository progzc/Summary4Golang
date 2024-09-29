package common_mistake

import (
	"testing"
)

// string 类型的变量值不能为 nil
func TestMistake_011(t *testing.T) {
	wrong011()
	right011()
}

func wrong011() {
	//var s string = nil // cannot use nil as type string in assignment
	//if s == nil {      // invalid operation: s == nil (mismatched types string and nil)
	//	s = "default"
	//}
}

func right011() {
	var s string // 字符串类型的零值是空串 ""
	if s == "" {
		s = "default"
	}
}
