package order_leetcode

import "testing"

func Test_leetcode_0287_findDuplicate(t *testing.T) {
	type params struct {
		nums []int
	}
	tests := []struct {
		p    params
		want int
	}{
		{
			p: params{
				nums: []int{1, 3, 4, 2, 2},
			},
			want: 2,
		},
		{
			p: params{
				nums: []int{3, 1, 3, 4, 2},
			},
			want: 3,
		},
		{
			p: params{
				nums: []int{1, 1},
			},
			want: 1,
		},
		{
			p: params{
				nums: []int{1, 1},
			},
			want: 1,
		},
		{
			p: params{
				nums: []int{1, 1, 2},
			},
			want: 1,
		},
	}
	for _, test := range tests {
		fact := leetcode_0287_findDuplicate(test.p.nums)
		if fact != test.want {
			t.Errorf("params=%v,want=%v,fact=%v",
				test.p, test.want, fact)
		}
	}
}

// leetcode_0287_findDuplicate 快慢指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func leetcode_0287_findDuplicate(nums []int) int {
	slow, fast := 0, 0
	slow = nums[slow]
	fast = nums[nums[fast]]
	for slow != fast {
		slow = nums[slow]
		fast = nums[nums[fast]]
	}
	pre1, pre2 := 0, slow
	for pre1 != pre2 {
		pre1 = nums[pre1]
		pre2 = nums[pre2]
	}
	return pre1
}
