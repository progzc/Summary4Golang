package common_mistake

import (
	"fmt"
	"testing"
)

// range 迭代 map
// Go 的运行时是有意打乱迭代顺序的，所以你得到的迭代结果可能不一致
func TestMistake_025(t *testing.T) {
	wrong025()
	right025()
}

func wrong025() {

}

func right025() {
	m := map[string]int{"one": 1, "two": 2, "three": 3, "four": 4}
	for k, v := range m {
		fmt.Println(k, v)
	}
}
