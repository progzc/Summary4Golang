package go_leetcode

import (
	"testing"
)

// 题目：给你一个整数数组 nums ，除某个元素仅出现 一次 外，其余每个元素都恰出现 三次 。
// 请你找出并返回那个只出现了一次的元素。

// 注意事项: go中异或和取反都是采用^
func Test_leetcode_0137_single_number_ii(t *testing.T) {
	type params struct {
		nums []int
	}
	tests := []struct {
		p    params
		want int
	}{
		{
			p: params{
				nums: []int{2, 2, 3, 2},
			},
			want: 3,
		},
		{
			p: params{
				nums: []int{0, 1, 0, 1, 0, 1, 99},
			},
			want: 99,
		},
	}
	for _, test := range tests {
		fact := leetcode_0137_single_number_ii_method4(test.p.nums)
		if fact != test.want {
			t.Errorf("params=%v,want=%v,fact=%v",
				test.p, test.want, fact)
		}
	}
}

// leetcode_0137_single_number_ii_method1 有限状态自动机
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func leetcode_0137_single_number_ii_method1(nums []int) int {
	ones, twos := 0, 0
	for _, num := range nums {
		ones = (ones ^ num) & (^twos)
		twos = (twos ^ num) & (^ones)
	}
	return ones
}

// leetcode_0137_single_number_ii_method2 统计法
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func leetcode_0137_single_number_ii_method2(nums []int) int {
	counts := make([]int, 32)
	for _, num := range nums {
		for j := 0; j < 32; j++ {
			counts[j] += num & 1
			num >>= 1
		}
	}
	res, m := 0, 3
	for i := 0; i < 32; i++ {
		res <<= 1
		res |= counts[31-i] % m
	}
	return res
}

// leetcode_0137_single_number_ii_method3 哈希表
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func leetcode_0137_single_number_ii_method3(nums []int) int {
	freqMap := map[int]int{}
	for _, num := range nums {
		freqMap[num]++
	}
	for num, count := range freqMap {
		if count == 1 {
			return num
		}
	}
	return 0
}

// leetcode_0137_single_number_ii_method4 依次确定每一个二进制位
// 时间复杂度: O(nlogC)
// 空间复杂度: O(1)
func leetcode_0137_single_number_ii_method4(nums []int) int {
	res := 0
	for i := 0; i < 32; i++ {
		total := 0
		for _, num := range nums {
			total += num >> i & 1
		}
		if total%3 > 0 {
			res |= 1 << i
		}
	}
	return res
}
