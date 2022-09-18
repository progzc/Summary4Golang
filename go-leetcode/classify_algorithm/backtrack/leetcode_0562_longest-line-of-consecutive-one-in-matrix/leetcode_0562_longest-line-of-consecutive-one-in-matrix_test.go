package leetcode_0562_longest_line_of_consecutive_one_in_matrix

import "testing"

func Test_longestLine(t *testing.T) {
	type args struct {
		mat [][]int
	}
	var tests = []struct {
		name string
		args args
		want int
	}{
		{
			name: "用例1",
			args: args{
				mat: [][]int{
					{1, 1, 0, 0, 1, 0, 0, 1, 1, 0},
					{1, 0, 0, 1, 0, 1, 1, 1, 1, 1},
					{1, 1, 1, 0, 0, 1, 1, 1, 1, 0},
					{0, 1, 1, 1, 0, 1, 1, 1, 1, 1},
					{0, 0, 1, 1, 1, 1, 1, 1, 1, 0},
					{1, 1, 1, 1, 1, 1, 0, 1, 1, 1},
					{0, 1, 1, 1, 1, 1, 1, 0, 0, 1},
					{1, 1, 1, 1, 1, 0, 0, 1, 1, 1},
					{0, 1, 0, 1, 1, 0, 1, 1, 1, 1},
					{1, 1, 1, 0, 1, 0, 1, 1, 1, 1},
				},
			},
			want: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := longestLine(tt.args.mat); got != tt.want {
				t.Errorf("longestLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
