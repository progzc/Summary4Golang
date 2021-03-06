package leetcode_0494_target_sum

// 0494.目标和
// https://leetcode-cn.com/problems/target-sum/

// findTargetSumWays 动态规划（0/1背包问题）
// 时间复杂度: O(N*T)
// 空间复杂度: O(T)
// 思路：转化为如下的0/1背包问题
// 数组和sum,目标和target, 正数和x,负数和y,则x+y=sum,x-y=target,
// 那么x=(target+sum)/2=newTarget,y=(target-sum)/2
// ==> 选nums里的数得到target的种数
// 特点：0-1背包不考虑元素顺序的组合问题
func findTargetSumWays(nums []int, target int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	// 注意点1: 由于y是负数,所以y=(target-sum)/2<0==>sum<target
	if (sum+target)%2 != 0 || sum < target {
		return 0
	}
	newTarget := (sum + target) / 2
	dp := make([]int, newTarget+1)
	// 注意点2: 当没有任何元素可以选取时，元素和只能是0,应的方案数是1
	dp[0] = 1
	for _, num := range nums {
		for j := newTarget; j >= num; j-- {
			dp[j] += dp[j-num]
		}
	}
	return dp[newTarget]
}
