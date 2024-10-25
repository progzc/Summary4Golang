package leetcode_1497_check_if_array_pairs_are_divisible_by_k

import (
	"fmt"
	"testing"
)

// 1497. 检查数组对是否可以被 k 整除
// https://leetcode.cn/problems/check-if-array-pairs-are-divisible-by-k

func TestCanArrange(t *testing.T) {
	arr, k := []int{1, 2, 3, 4, 5, 6}, 7
	fmt.Println(canArrange(arr, k))

	fmt.Println(3 + (-1 % 3)) // 2
	fmt.Println(3 + (-2 % 3)) // 1

	arr, k = []int{-1, 1, -2, 2, -3, 3, -4, 4}, 3
	fmt.Println(canArrange(arr, k))
}

// canArrange 哈希
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 注意: arr中有可能有负数
func canArrange(arr []int, k int) bool {
	m := make(map[int]int)
	for _, num := range arr {
		target := (k + num%k) % k
		m[target]++
	}
	for key, value := range m {
		if key == 0 {
			if value%2 != 0 {
				return false
			}
		} else {
			if m[key] != m[k-key] {
				return false
			}
		}
	}
	return true
}
