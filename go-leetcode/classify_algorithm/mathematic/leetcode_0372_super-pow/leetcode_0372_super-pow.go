package leetcode_0372_super_pow

// 372. 超级次方
// https://leetcode.cn/problems/super-pow/

// superPow
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(1)
func superPow(a int, b []int) int {
	ans := 1
	const mod = 1337
	for i := len(b) - 1; i >= 0; i-- {
		ans = ans * pow(a, b[i]) % mod
		a = pow(a, 10)
	}
	return ans
}

// pow x,n 均为正数
func pow(x, n int) int {
	const mod = 1337
	res := 1
	for n > 0 {
		if n&1 == 1 {
			res = res * x % mod
		}
		x = x * x % mod
		n >>= 1
	}
	return res
}
