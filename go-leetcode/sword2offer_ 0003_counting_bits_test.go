package go_leetcode

import (
	"math/bits"
	"reflect"
	"testing"
)

func TestCountBits(t *testing.T) {
	type params struct {
		n int
	}
	tests := []struct {
		p    params
		want []int
	}{
		{
			p: params{
				n: 2,
			},
			want: []int{0, 1, 1},
		},
		{
			p: params{
				n: 5,
			},
			want: []int{0, 1, 1, 2, 1, 2},
		},
	}
	for _, test := range tests {
		fact := countBitsMethod3(test.p.n)
		if !reflect.DeepEqual(test.want, fact) {
			t.Errorf("params=%v,want=%v,fact=%v",
				test.p, test.want, fact)
		}
	}
}

// countBitsMethod4 位运算
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func countBitsMethod4(n int) []int {
	res := make([]int, n+1)
	// 注意迭代从1开始
	for i := 1; i <= n; i++ {
		// 若i为偶数,i&(i-1)=0,若i为奇数,i&(i-1)=i-1
		// 巧妙利用了默认值res[0]=0,且偶数的二进制中只有一个1
		res[i] = res[i&(i-1)] + 1
	}
	return res
}

// countBitsMethod3 分奇偶讨论
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func countBitsMethod3(n int) []int {
	res := make([]int, n+1)
	for i := 0; i <= n; i++ {
		if i%2 == 0 {
			// 如果当前num为偶数，则二进制中1的个数为 num/2 的二进制中1的个数
			res[i] = res[i/2]
		} else {
			// 如果当前num为奇数，则二进制中1的个数为前一个数二进制1的个数+1
			res[i] = res[i-1] + 1
		}
	}
	return res
}

// countBitsMethod2 标准库调用
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func countBitsMethod2(n int) []int {
	res := make([]int, n+1)
	for i := 0; i <= n; i++ {
		res[i] = bits.OnesCount(uint(i))
	}
	return res
}

// countBitsMethod1 常规计算方法
// 时间复杂度: O(n*31)
// 空间复杂度: O(n)
func countBitsMethod1(n int) []int {
	var res []int
	for i := 0; i <= n; i++ {
		c := 0
		a := i
		for j := 0; j < 32; j++ {
			if a&1 != 0 {
				c++
			}
			a = a >> 1
		}
		res = append(res, c)
	}
	return res
}
