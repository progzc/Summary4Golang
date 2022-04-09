package order_leetcode

import (
	"reflect"
	"testing"
)

// go_knowledge test -v -run=TwoSum$ leetcode_0001_two_sum_test.go_knowledge
func Test_leetcode_0001_twoSum(t *testing.T) {
	tests := []struct {
		nums   []int
		target int
		want   []int
	}{
		{
			nums:   []int{2, 7, 11, 15},
			target: 9,
			want:   []int{0, 1},
		},
		{
			nums:   []int{3, 2, 4},
			target: 6,
			want:   []int{1, 2},
		},
		{
			nums:   []int{3, 3},
			target: 6,
			want:   []int{0, 1},
		},
	}
	for _, test := range tests {
		fact := leetcode_0001_twoSum(test.nums, test.target)
		// util.SameIntSlice是无序的
		if !sameIntSlice(fact, test.want) {
			t.Errorf("nums=%v,target=%d,want=%v,fact=%v",
				test.nums, test.target, test.want, fact)
		}
	}
}

func Test_leetcode_0001_twoSum2(t *testing.T) {
	tests := []struct {
		nums   []int
		target int
		want   []int
	}{
		{
			nums:   []int{2, 7, 11, 15},
			target: 9,
			want:   []int{0, 1},
		},
		{
			nums:   []int{3, 2, 4},
			target: 6,
			want:   []int{1, 2},
		},
		{
			nums:   []int{3, 3},
			target: 6,
			want:   []int{0, 1},
		},
	}
	for _, test := range tests {
		fact := leetcode_0001_twoSum(test.nums, test.target)
		// reflect.DeepEqual是有序的
		if !reflect.DeepEqual(fact, test.want) {
			t.Errorf("nums=%v,target=%d,want=%v,fact=%v\n",
				test.nums, test.target, test.want, fact)
		}
	}
}

func leetcode_0001_twoSum(nums []int, target int) []int {
	hashTable := map[int]int{}
	for i, num := range nums {
		if p, ok := hashTable[target-num]; ok {
			return []int{i, p}
		}
		hashTable[num] = i
	}
	return nil
}
