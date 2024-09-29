package common_mistake

import (
	"fmt"
	"testing"
)

// 旧 slice
// (1) 当你从一个已存在的 slice 创建新 slice 时，二者的数据指向相同的底层数组。
//
//	如果你的程序使用这个特性，那需要注意 "旧"（stale） slice 问题。
//
// (2) 某些情况下，向一个 slice 中追加元素而它指向的底层数组容量不足时，将会重新分配一个新数组来存储数据。
//
//	而其他 slice 还指向原来的旧底层数组。
func TestMistake_044(t *testing.T) {
	wrong044()
	right044()
}

func wrong044() {
}

func right044() {
	s1 := []int{1, 2, 3}
	fmt.Println(len(s1), cap(s1), s1) // 3 3 [1 2 3]

	s2 := s1[1:]
	fmt.Println(len(s2), cap(s2), s2) // 2 2 [2 3]

	for i := range s2 {
		s2[i] += 20
	}
	// 此时的 s1 与 s2 是指向同一个底层数组的
	fmt.Println(s1) // [1 22 23]
	fmt.Println(s2) // [22 23]

	s2 = append(s2, 4) // 向容量为 2 的 s2 中再追加元素，此时将分配新数组来存
	for i := range s2 {
		s2[i] += 10
	}
	fmt.Println(s1) // [1 22 23]    // 此时的 s1 不再更新，为旧数据
	fmt.Println(s2) // [32 33 14]
}
