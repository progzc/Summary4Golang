package leetcode_0487_max_consecutive_ones_ii

// 0487. 最大连续1的个数 II
// https://leetcode.cn/problems/max-consecutive-ones-ii/
// 解题方法：前后缀和、动态规划、滑动窗口

// findMaxConsecutiveOnes 预处理+枚举
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：
//	枚举每个0的位置i,则最大连续的1的个数=pre[i-1]+suff[i+1]
//	pre[i-1]表示以i-1结尾往前延伸的最大连续1的个数
//	suff[i+1]表示以i+1开头往后延伸的最大连续1的个数
// 	下面是经过预处理求解pre[i]和suff[i]的递推式
//	递推式：
//		pre[i] = pre[i-1]+1 (if nums[i]==1)
//		pre[i] = 0 (if nums[i]==0)
//		suff[i] = suff[i+1]+1 (if nums[i]==1)
//		suff[i] = 0 (if nums[i]==0)
func findMaxConsecutiveOnes(nums []int) int {
	n, ans := len(nums), 0
	pre, suff := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		if nums[i] == 0 {
			pre[i] = 0
		} else if i == 0 {
			pre[i] = 1
		} else {
			pre[i] = pre[i-1] + 1
		}
		ans = max(ans, pre[i])
	}
	for i := n - 1; i >= 0; i-- {
		if nums[i] == 0 {
			suff[i] = 0
		} else if i == n-1 {
			suff[i] = 1
		} else {
			suff[i] = suff[i+1] + 1
		}
	}
	for i := 0; i < n; i++ {
		if nums[i] == 0 {
			res := 0
			if i > 0 {
				res += pre[i-1]
			}
			if i < n-1 {
				res += suff[i+1]
			}
			ans = max(ans, res+1)
		}
	}
	return ans
}

// findMaxConsecutiveOnes_2 动态规划
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路:
//	定义状态:
//		dp[i][0] 为考虑到以 i 为结尾未使用操作将 [0,i] 某个 0 变成 1 的最大的连续 1 的个数
//		dp[i][1] 为考虑到以 i 为结尾使用操作将 [0,i] 某个 0 变成 1 的最大的连续 1 的个数
//	状态转移方程:
//		若nums[i]==1,则dp[i][0] = dp[i-1][0]+1
//		若nums[i]==0,则dp[i][0] = 0
//
//		若nums[i]==1,则dp[i][1] = dp[i-1][1]+1
//		若nums[i]==0,则dp[i][1] = dp[i-1][0]+1
func findMaxConsecutiveOnes_2(nums []int) int {
	n, ans := len(nums), 0
	dp0, dp1 := 0, 0
	for i := 0; i < n; i++ {
		if nums[i] > 0 {
			dp0++
			dp1++
		} else {
			dp1 = dp0 + 1
			dp0 = 0
		}
		ans = max(ans, max(dp0, dp1))
	}
	return ans
}

// findMaxConsecutiveOnes_3 滑动窗口
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：题目等价于给定一个区间，该区间中最多只能包含1个0，求出该区间的最大长度
// 解法：只用维护一个区间，这个区间中最多只包含一个0。当区间中包含两个0的时候，直接移动左边界l直到区间只包含一个0即可。
//		这个过程中去更新最大区间长度，最后就能得到答案。
func findMaxConsecutiveOnes_3(nums []int) int {
	ans, count, n := 0, 0, len(nums)
	for l, r := 0, 0; r < n; r++ {
		if nums[r] == 0 {
			count++
			for count > 1 {
				if nums[l] == 0 {
					count--
					l++
				} else {
					l++
				}
			}
		}
		ans = max(ans, r-l+1)
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
