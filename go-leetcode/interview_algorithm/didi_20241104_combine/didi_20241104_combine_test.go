package didi_20241104_combine

import (
	"fmt"
	"testing"
)

// 滴滴（一面）
// 排列组合的方式

func TestCombine(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	fmt.Println(combine(nums)) // [[1 2 3 4] [1 2 4 3] [1 3 2 4] [1 3 4 2] [1 4 2 3] [1 4 3 2] [2 1 3 4] [2 1 4 3] [2 3 1 4] [2 3 4 1] [2 4 1 3] [2 4 3 1] [3 1 2 4] [3 1 4 2] [3 2 1 4] [3 2 4 1] [3 4 1 2] [3 4 2 1] [4 1 2 3] [4 1 3 2] [4 2 1 3] [4 2 3 1] [4 3 1 2] [4 3 2 1]]
}

func combine(nums []int) [][]int {
	var (
		ans [][]int
		dfs func(idx int, output []int)
	)
	n := len(nums)
	choose := make([]bool, n)
	dfs = func(idx int, output []int) {
		if idx == n {
			ans = append(ans, append([]int{}, output...))
			return
		}
		for i := 0; i < n; i++ {
			if choose[i] {
				continue
			}
			choose[i] = true
			output = append(output, nums[i])
			dfs(idx+1, output)
			output = output[:len(output)-1]
			choose[i] = false
		}
	}
	dfs(0, nil)
	return ans
}
