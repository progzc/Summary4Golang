package leetcode_0651_4_keys_keyboard

// 0651. 4键键盘
// https://leetcode.cn/problems/4-keys-keyboard/

// maxA 动态规划
// 空间复杂度: O(n)
// 时间复杂度: O(n^2)
// 思路：最后一步，要么是A，要么是Ctrl-V
//	a.若最后一步是A，那么dp[i]=dp[i-1]+1
//	b.若最后一步是Ctrl-V，那么假设第j步是dp[j], 第j+1步执行Ctrl-A，那么第j+2步执行Ctrl-C,剩下的i-(j+2)步执行Ctrl-V
//	  则有dp[i]=max(dp[i], dp[j]*(i-(j+2)+1)) j+2<i
func maxA(n int) int {
	// 巧妙利用dp[0]=0，这样可以使代码更简洁
	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = dp[i-1] + 1
		// j+2<i的缘由：i-(j+2)>0
		for j := 2; j+2 < i; j++ {
			// 解释 i-(j+2)+1：i-(j+2)是Ctrl-V执行的次数, 末尾的+1是指原来的倍数
			dp[i] = max(dp[i], dp[j]*(i-(j+2)+1))
		}
	}
	return dp[n]
}
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
