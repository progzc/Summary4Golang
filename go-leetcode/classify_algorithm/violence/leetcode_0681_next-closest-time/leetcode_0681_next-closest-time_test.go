package leetcode_0681_next_closest_time

import "testing"

func Test_nextClosestTime_2(t *testing.T) {
	type args struct {
		time string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "用例1",
			args: args{
				time: "19:34",
			},
			want: "19:39",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextClosestTime_2(tt.args.time); got != tt.want {
				t.Errorf("nextClosestTime_2() = %v, want %v", got, tt.want)
			}
		})
	}
}
