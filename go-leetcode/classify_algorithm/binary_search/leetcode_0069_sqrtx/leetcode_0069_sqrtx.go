package leetcode_40069_sqrtx

// 0069. x 的平方根
// https://leetcode.cn/problems/sqrtx/

// mySqrt 二分查找
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
func mySqrt(x int) int {
	l, r := 0, x
	ans := -1
	for l <= r {
		mid := l + (r-l)>>2
		if mid*mid <= x {
			ans = mid
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return ans
}
