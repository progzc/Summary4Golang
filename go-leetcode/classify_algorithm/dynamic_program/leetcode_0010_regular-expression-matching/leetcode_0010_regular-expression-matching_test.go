package leetcode_0010_regular_expression_matching

import "testing"

func Test_isMatch(t *testing.T) {
	type args struct {
		s string
		p string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "用例-1",
			args: args{
				s: "aa",
				p: "a*",
			},
			want: true,
		},
		{
			name: "用例-2",
			args: args{
				s: "aab",
				p: "c*a*b",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMatch(tt.args.s, tt.args.p); got != tt.want {
				t.Errorf("isMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
