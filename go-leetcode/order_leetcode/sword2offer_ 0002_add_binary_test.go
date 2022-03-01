package order_leetcode

import (
	"strconv"
	"testing"
)

func Test_sword2offer_0002_addBinary(t *testing.T) {
	type params struct {
		a string
		b string
	}
	tests := []struct {
		p    params
		want string
	}{
		{
			p: params{
				a: "11",
				b: "10",
			},
			want: "101",
		},
		{
			p: params{
				a: "1010",
				b: "1011",
			},
			want: "10101",
		},
	}
	for _, test := range tests {
		fact := sword2offer_0002_addBinaryMethod1(test.p.a, test.p.b)
		if fact != test.want {
			t.Errorf("params=%v,want=%v,fact=%v",
				test.p, test.want, fact)
		}
	}
}

// sword2offer_0002_addBinaryMethod1 模拟二进制的运算
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func sword2offer_0002_addBinaryMethod1(a string, b string) string {
	ans := ""
	carry := 0
	lenA, lenB := len(a), len(b)
	n := max(lenA, lenB)
	for i := 0; i < n; i++ {
		if i < lenA {
			carry += int(a[lenA-i-1] - '0')
		}
		if i < lenB {
			carry += int(b[lenB-i-1] - '0')
		}
		ans = strconv.Itoa(carry%2) + ans
		carry /= 2
	}
	if carry > 0 {
		ans = "1" + ans
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
