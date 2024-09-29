package lianying_20221202_quchong

import (
	"fmt"
	"sort"
)

func removeDuplicateInplace(nums []int) int {
	n := len(nums)
	if n < 2 {
		return n
	}

	sort.Ints(nums)
	// [1,2,2,3,4,4,4,5]

	i := 0
	next := 0
	for i < n {
		start, end := i, i
		for end < n && nums[end] == nums[start] {
			end++
		}
		i = end
		nums[next] = nums[start]
		next++
	}
	fmt.Println(nums, next)
	return next
}
