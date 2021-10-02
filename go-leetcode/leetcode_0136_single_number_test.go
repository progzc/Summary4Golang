package go_leetcode

import "testing"

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
			t.Errorf("params=%v,want=%T,fact=%T",
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
