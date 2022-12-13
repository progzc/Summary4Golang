package leetcode_0435_non_overlapping_intervals

import "testing"

func Test_eraseOverlapIntervals(t *testing.T) {
	type args struct {
		intervals [][]int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "用例-1",
			args: args{
				//[[-52,31],[-73,-26],[82,97],[-65,-11],[-62,-49],[95,99],[58,95],[-31,49],[66,98],[-63,2],[30,47],[-40,-26]]
				intervals: [][]int{
					{-52, 31}, {-73, -26}, {82, 97}, {-65, -11}, {-62, -49},
					{95, 99}, {58, 95}, {-31, 49}, {66, 98}, {-63, 2},
					{30, 47}, {-40, -26}},
			},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := eraseOverlapIntervals(tt.args.intervals); got != tt.want {
				t.Errorf("eraseOverlapIntervals() = %v, want %v", got, tt.want)
			}
		})
	}
}
