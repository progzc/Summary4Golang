package leetcode_0050_powx_n

// 50. Pow(x, n)
// https://leetcode.cn/problems/powx-n/

// myPow 快速幂（分治法）+ 递归
// 时间复杂度: O(log(n))
// 空间复杂度: O(log(n))
// 思路: 二分+递归
func myPow(x float64, n int) float64 {
	var quickMul func(x float64, n int) float64
	quickMul = func(x float64, n int) float64 {
		if n == 0 {
			return 1
		}
		y := quickMul(x, n/2)
		if n%2 == 0 {
			return y * y
		}
		return y * y * x
	}
	if n >= 0 {
		return quickMul(x, n)
	}
	return 1.0 / quickMul(x, -n)
}

// myPow_2 快速幂（分治法）+ 迭代
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
// 思路: 将n进行二进制展开。
//	例如：求x^n, x=0.1, n=10;
//		n=(1010)二进制
//		x^10 = x^(1*0)*x^(2*1)*x^(4*0)*x^(8*1)
func myPow_2(x float64, n int) float64 {
	if x == 0.0 {
		return 0.0
	}

	res := 1.0
	if n < 0 {
		x = 1.0 / x
		n = -n
	}
	for n > 0 {
		if n&1 == 1 {
			res *= x
		}
		x *= x
		n >>= 1
	}
	return res
}
