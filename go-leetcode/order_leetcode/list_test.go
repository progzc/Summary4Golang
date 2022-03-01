package order_leetcode

import (
	"reflect"
	"testing"
)

func TestListNode(t *testing.T) {
	tests := []struct {
		nums []int
		want []int
	}{
		{
			nums: []int{2, 4, 3},
			want: []int{2, 4, 3},
		},
		{
			nums: []int{},
			want: nil,
		},
	}
	for _, test := range tests {
		fact := listNode2Slice(initListNode(test.nums))
		if !reflect.DeepEqual(fact, test.want) {
			t.Errorf("nums=%v,want=%v,fact=%v",
				test.nums, test.want, fact)
		}
	}
}
