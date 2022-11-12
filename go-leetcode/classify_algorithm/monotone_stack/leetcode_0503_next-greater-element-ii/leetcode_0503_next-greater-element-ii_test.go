package leetcode_0503_next_greater_element_ii

import (
	"reflect"
	"testing"
)

func Test_nextGreaterElements(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "用例1",
			args: args{
				nums: []int{5, 4, 3, 2, 1},
			},
			want: []int{-1, 5, 5, 5, 5},
		},
		{
			name: "用例2",
			args: args{
				nums: []int{1, 2, 3, 2, 1},
			},
			want: []int{2, 3, -1, 3, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextGreaterElements(tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nextGreaterElements() = %v, want %v", got, tt.want)
			}
		})
	}
}
