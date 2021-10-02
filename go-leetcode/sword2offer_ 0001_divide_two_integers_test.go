package go_leetcode

import (
	"math"
	"testing"
)

func TestDivide(t *testing.T) {
	type params struct {
		a int
		b int
	}
	tests := []struct {
		p    params
		want int
	}{
		{
			p: params{
				a: math.MinInt32,
				b: -1,
			},
			want: math.MaxInt32,
		},
		{
			p: params{
				a: 15,
				b: 2,
			},
			want: 7,
		},
		{
			p: params{
				a: 7,
				b: -3,
			},
			want: -2,
		},
		{
			p: params{
				a: 20,
				b: math.MinInt32,
			},
			want: 0,
		},
	}
	for _, test := range tests {
		fact := divide(test.p.a, test.p.b)
		if fact != test.want {
			t.Errorf("params=%v,want=%d,fact=%d",
				test.p, test.want, fact)
		}
	}
}

func divide(a int, b int) int {
	if a == math.MinInt32 && b == -1 {
		return math.MaxInt32
	}
	sign := 1
	if (a > 0 && b < 0) || (a < 0 && b > 0) {
		sign = -1
	}
	a = abs(a)
	b = abs(b)
	res := 0
	for i := 31; i >= 0; i-- {
		if (a>>i)-b >= 0 { // 等价于 a-(b<<i)>=0,但这种写法有可能b<<i会越界
			a = a - (b << i)
			res += 1 << i
		}
	}
	if sign == -1 {
		return -res
	}
	return res
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
