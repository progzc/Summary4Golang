package leetcode_0070_climbing_stairs

// 0070.爬楼梯
// https://leetcode-cn.com/problems/climbing-stairs/

// climbStairs 动态规划
// 时间复杂度：O(n)
// 空间复杂度：O(1)
// 思路：f(n)=f(n-1)+f(n-2)
func climbStairs(n int) int {
	if n == 1 {
		return 1
	}
	if n == 2 {
		return 2
	}

	f1, f2 := 1, 2
	for i := 3; i <= n; i++ {
		f1, f2 = f2, f1+f2
	}
	return f2
}
