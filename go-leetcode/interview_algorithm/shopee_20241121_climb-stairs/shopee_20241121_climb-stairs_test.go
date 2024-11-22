package shopee_20241121_climb_stairs

import (
	"fmt"
	"testing"
)

// 虾皮（一面）
// 一个人每次最多走m阶台阶，问走n阶台阶一共有多少种走法？输出m=4, n=5的结果

func TestClimbStairs(t *testing.T) {
	m, n := 2, 3
	fmt.Println(climbStairs(m, n)) // 3

	m, n = 4, 5
	fmt.Println(climbStairs(m, n)) // 15
}

func TestClimbStairs2(t *testing.T) {
	m, n := 2, 3
	fmt.Println(climbStairs2(m, n)) // 3

	m, n = 4, 5
	fmt.Println(climbStairs2(m, n)) // 15
}

func climbStairs(m, n int) int {
	dp := make([]int, n+1)
	dp[0] = 1
	for i := 1; i <= n; i++ { // 先遍历背包
		for j := 1; j <= m; j++ { // 再遍历物品
			if i-j >= 0 {
				dp[i] += dp[i-j]
			}
		}
	}
	return dp[n]
}

func climbStairs2(m, n int) int {
	var ans int
	var dfs func(leave int)
	dfs = func(leave int) {
		if leave == 0 {
			ans++
			return
		}
		for i := 1; i <= m; i++ {
			if leave-i >= 0 {
				dfs(leave - i)
			}
		}
	}
	dfs(n)
	return ans
}
