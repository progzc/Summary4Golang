package leetcode_0343_integer_break

import "testing"

func Test_integerBreak(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "用例-1",
			args: args{
				n: 4,
			},
			want: 4,
		},
		{
			name: "用例-2",
			args: args{
				n: 6,
			},
			want: 9,
		},
		{
			name: "用例-3",
			args: args{
				n: 10,
			},
			want: 36,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := integerBreak(tt.args.n); got != tt.want {
				t.Errorf("integerBreak() = %v, want %v", got, tt.want)
			}
		})
	}
}
