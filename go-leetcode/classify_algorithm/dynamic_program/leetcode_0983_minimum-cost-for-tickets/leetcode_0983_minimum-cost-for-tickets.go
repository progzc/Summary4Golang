package leetcode_0983_minimum_cost_for_tickets

import "math"

// 0983.最低票价
// https://leetcode-cn.com/problems/minimum-cost-for-tickets/

// mincostTickets 回溯法
// 时间复杂度：O(3*W)，W=365
// 空间复杂度: O(W)
// 思路：
//	定义dp(i)：从第i天开始到一年的结束，我们需要花的钱。对于一年中的任意一天，
//		如果这一天不是必须出行的日期，则：dp(i)=dp(i+1)
//		如果这一天是必须出行的日期，则：dp(i)=min{cost(j)+dp(i+j)},j∈{1,7,30}
func mincostTickets(days []int, costs []int) int {
	mMap := [366]int{}
	dMap := map[int]bool{}
	for _, day := range days {
		dMap[day] = true
	}
	var dp func(day int) int
	dp = func(day int) int {
		if day > 365 {
			return 0
		}
		// 若之前记忆过中间结果，则直接返回
		if mMap[day] > 0 {
			return mMap[day]
		}
		if dMap[day] {
			// 如果这一天是必须出行的日期
			s1 := dp(day+1) + costs[0]
			s2 := dp(day+7) + costs[1]
			s3 := dp(day+30) + costs[2]
			// 记忆中间结果
			mMap[day] = min(min(s1, s2), s3)
		} else {
			// 如果这一天不是必须出行的日期
			mMap[day] = dp(day + 1)
		}
		return mMap[day]
	}
	return dp(1)
}

// mincostTickets_2 回溯法（优化）
// 时间复杂度：O(N)，N为days的长度
// 空间复杂度: O(N)
// 思路：
//	定义：dp(i)：能够完成从第days[i]天到最后的旅行计划的最小花费；j_1：满足days[j_1]>=days[i]+1的最小下标；
//		 j_7：满足days[j_7]>=days[i]+7的最小下标；j_30：满足days[j_30]>=days[i]+30的最小下标。
//	那么：dp(i)=min(dp(j_1)+costs[0],dp(j_7)+costs[1],dp(j_30)+costs[2])
// 技巧：思路是动态规划，但是写法是递归回溯。本质原因是要从后向前动态规划
func mincostTickets_2(days []int, costs []int) int {
	mMap := [366]int{}
	durations := []int{1, 7, 30}

	var dp func(i int) int
	dp = func(i int) int {
		if i >= len(days) {
			return 0
		}
		// 若之前记忆过中间结果，则直接返回
		if mMap[i] > 0 {
			return mMap[i]
		}
		mMap[i] = math.MaxInt32
		j := i
		for idx, duration := range durations {
			for j < len(days) && days[j] < days[i]+duration {
				j++
			}
			mMap[i] = min(mMap[i], dp(j)+costs[idx])
		}
		return mMap[i]
	}
	return dp(0)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
