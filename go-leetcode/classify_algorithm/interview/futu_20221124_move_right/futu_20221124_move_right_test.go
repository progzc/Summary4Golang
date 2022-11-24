package futu_20221124_move_right

import "testing"

func Test_translate(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "用例1",
			args: args{
				nums: []int{5, 1, 2, 3, 4},
			},
			want: 1,
		},
		{
			name: "用例2",
			args: args{
				nums: []int{5, 6, 1, 2, 3, 4},
			},
			want: 2,
		},
		{
			name: "用例3",
			args: args{
				nums: []int{1, 2, 3, 4},
			},
			want: 0,
		},
		{
			name: "用例4",
			args: args{
				nums: []int{5, 4, 3, 2, 1},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := translate(tt.args.nums); got != tt.want {
				t.Errorf("translate() = %v, want %v", got, tt.want)
			}
		})
	}
}
