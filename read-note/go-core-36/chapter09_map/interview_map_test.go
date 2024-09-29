package chapter09_map

import (
	"fmt"
	"testing"
)

func Test_interview_map_1(t *testing.T) {
	m := map[int]int{}
	m[1] = 1
	b := m[1]
	b += 1
	fmt.Println(m[1]) // 1

	m2 := make(map[int]*int)
	x := 1
	m2[1] = &x
	fmt.Println(*m2[1], m2[1]) // 1 0xc00000a3c0
	c := m2[1]
	*c = *c + 1
	fmt.Println(*m2[1], m2[1]) // 2 0xc00000a3c0
}
