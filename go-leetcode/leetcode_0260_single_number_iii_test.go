package go_leetcode

import (
	"reflect"
	"testing"
)

// 题目：给定一个整数数组 nums，其中恰好有两个元素只出现一次，其余所有元素均出现两次。
// 找出只出现一次的那两个元素。你可以按 任意顺序 返回答案。

func Test_leetcode_0260_single_number_iii(t *testing.T) {
	type params struct {
		nums []int
	}
	tests := []struct {
		p    params
		want []int
	}{
		{
			p: params{
				nums: []int{1, 2, 1, 3, 2, 5},
			},
			want: []int{3, 5},
		},
		{
			p: params{
				nums: []int{-1, 0},
			},
			want: []int{-1, 0},
		},
		{
			p: params{
				nums: []int{0, 1},
			},
			want: []int{1, 0},
		},
	}
	for _, test := range tests {
		fact := leetcode_0260_single_number_iii(test.p.nums)
		if !reflect.DeepEqual(fact, test.want) {
			t.Errorf("params=%v,want=%v,fact=%v",
				test.p, test.want, fact)
		}
	}
}

// leetcode_0260_single_number_iii 分组异或
// 时间复杂度：O(n)
// 空间复杂度：O(1)
func leetcode_0260_single_number_iii(nums []int) []int {
	ret := 0
	for _, num := range nums {
		ret ^= num
	}
	div := 1
	for div&ret == 0 {
		div <<= 1
	}
	a, b := 0, 0
	for _, num := range nums {
		if div&num != 0 {
			a ^= num
		} else {
			b ^= num
		}
	}
	return []int{a, b}
}
