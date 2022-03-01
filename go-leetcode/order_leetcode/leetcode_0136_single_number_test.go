package order_leetcode

import "testing"

// 题目：给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。
// 找出那个只出现了一次的元素。

func Test_leetcode_0136_singleNumber(t *testing.T) {
	type params struct {
		nums []int
	}
	tests := []struct {
		p    params
		want int
	}{
		{
			p: params{
				nums: []int{2, 2, 1},
			},
			want: 1,
		},
		{
			p: params{
				nums: []int{4, 1, 2, 1, 2},
			},
			want: 4,
		},
	}
	for _, test := range tests {
		fact := leetcode_0136_singleNumber(test.p.nums)
		if fact != test.want {
			t.Errorf("params=%v,want=%v,fact=%v",
				test.p, test.want, fact)
		}
	}
}

func leetcode_0136_singleNumber(nums []int) int {
	res := 0
	for _, num := range nums {
		res ^= num
	}
	return res
}
