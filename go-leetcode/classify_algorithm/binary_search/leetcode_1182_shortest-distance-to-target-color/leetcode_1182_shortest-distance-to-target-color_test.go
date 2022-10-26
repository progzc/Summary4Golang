package leetcode_1182_shortest_distance_to_target_color

import (
	"fmt"
	"sort"
	"testing"
)

func Test_SearchInts(t *testing.T) {
	nums := []int{4, 7, 8}
	idx1 := sort.SearchInts(nums, 1)
	idx2 := sort.SearchInts(nums, 5)
	idx3 := sort.SearchInts(nums, 7)
	idx4 := sort.SearchInts(nums, 9)
	fmt.Println(idx1, idx2, idx3, idx4) // 0 1 1 3
}

func Test_SearchInts_2(t *testing.T) {
	nums := []int{4, 7, 7, 8}
	idx1 := sort.SearchInts(nums, 1)
	idx2 := sort.SearchInts(nums, 7)
	idx3 := sort.SearchInts(nums, 9)
	fmt.Println(idx1, idx2, idx3) // 0 1 4
}
