package common_mistake

import (
	"testing"
)

// 使用简短声明来重复声明变量
// 不能用简短声明方式来单独为一个变量重复声明， := 左侧至少有一个新变量，才允许多变量的重复声明。
func TestMistake_005(t *testing.T) {
	wrong005()
	right005()
}

func wrong005() {
	//one := 0
	//one := 1 // error: no new variables on left side of :=
}

func right005() {
	one := 0
	one, two := 1, 2    // two 是新变量，允许 one 的重复声明。比如 error 处理经常用同名变量 err
	one, two = two, one // 交换两个变量值的简写
}
