package go_leetcode

import (
	"testing"
)

func TestSameIntSlice(t *testing.T) {
	type params struct {
		x []int
		y []int
	}
	tests := []struct {
		p    params
		want bool
	}{
		{
			p: params{
				x: []int{2, 3, 4},
				y: []int{3, 4, 2},
			},
			want: true,
		},
		{
			p: params{
				x: []int{2, 3, 4},
				y: []int{2, 4, 2},
			},
			want: false,
		},
	}
	for _, test := range tests {
		fact := sameIntSlice(test.p.x, test.p.y)
		if fact != test.want {
			t.Errorf("params=%v,want=%T,fact=%T",
				test.p, test.want, fact)
		}
	}
}
