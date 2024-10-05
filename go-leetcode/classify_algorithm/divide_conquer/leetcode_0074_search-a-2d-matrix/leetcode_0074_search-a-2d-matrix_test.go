package leetcode_0074_search_a_2d_matrix

import (
	"fmt"
	"sort"
	"testing"
)

func TestSearch(t *testing.T) {
	nums := []int{1, 3, 3, 3, 4, 4, 5, 6}
	fmt.Println(sort.SearchInts(nums, 2))  // 1
	fmt.Println(sort.SearchInts(nums, -1)) // 0
	fmt.Println(sort.SearchInts(nums, 3))  // 3
	fmt.Println(sort.Search(len(nums), func(i int) bool {
		return nums[i] > 3
	})) // 4
}
