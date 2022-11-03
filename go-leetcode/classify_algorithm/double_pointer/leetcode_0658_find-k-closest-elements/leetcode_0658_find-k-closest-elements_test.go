package leetcode_0658_find_k_closest_elements

import (
	"fmt"
	"sort"
	"testing"
)

func Test_findClosestElements(t *testing.T) {
	//arr, k, x := []int{1, 2, 3, 4, 5}, 4, 3 // 输出: [1,2,3,4]
	//arr, k, x := []int{1, 2, 3, 4, 5}, 4, -1 // 输出: [1,2,3,4]
	arr, k, x := []int{0, 1, 1, 1, 2, 3, 6, 7, 8, 9}, 9, 4 // 输出: [0 1 1 1 2 3 6 7 8]
	fmt.Println(findClosestElements_2(arr, k, x))
}

func Test_SearchInts(t *testing.T) {
	arr, x := []int{1, 2, 3, 3, 3}, 3
	fmt.Println(sort.SearchInts(arr, x)) // 2
}
