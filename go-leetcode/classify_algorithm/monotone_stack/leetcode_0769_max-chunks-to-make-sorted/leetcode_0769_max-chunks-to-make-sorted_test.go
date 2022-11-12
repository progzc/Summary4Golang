package leetcode_0769_max_chunks_to_make_sorted

import "testing"

func Test_maxChunksToSorted(t *testing.T) {
	type args struct {
		arr []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "用例1",
			args: args{
				arr: []int{4, 3, 2, 1, 0},
			},
			want: 1,
		},
		{
			name: "用例2",
			args: args{
				arr: []int{1, 0, 2, 3, 4},
			},
			want: 4,
		},
		{
			name: "用例3",
			args: args{
				arr: []int{3, 0, 2, 1, 4},
			},
			want: 2,
		},
		{
			name: "用例4",
			args: args{
				arr: []int{1, 2, 0, 3},
			},
			want: 2,
		},
		{
			name: "用例5",
			args: args{
				arr: []int{0, 1, 2, 4, 3},
			},
			want: 4,
		},
		{
			name: "用例6",
			args: args{
				arr: []int{0, 1, 3, 4, 2},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxChunksToSorted(tt.args.arr); got != tt.want {
				t.Errorf("maxChunksToSorted() = %v, want %v", got, tt.want)
			}
		})
	}
}
