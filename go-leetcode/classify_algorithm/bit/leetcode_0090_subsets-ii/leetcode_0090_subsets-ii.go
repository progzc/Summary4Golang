package leetcode_0090_subsets_ii

import "sort"

// 0090.子集 II
// https://leetcode-cn.com/problems/subsets-ii/

// subsetsWithDup 位操作
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(n)
func subsetsWithDup(nums []int) [][]int {
	ans := [][]int{}
	n := len(nums)
	sort.Ints(nums)
	for mask := 0; mask < 1<<n; mask++ {
		set := []int{}
		flag := true
		for i, num := range nums {
			if mask&(1<<i) > 0 {
				// 对于当前选择的数x，若前面有与其相同的数y，且没有选择y，此时包含x的子集，必然会出现在包含y的所有子集中
				if i > 0 && (mask&(1<<(i-1)) == 0) && nums[i-1] == nums[i] {
					flag = false
					break
				}
				set = append(set, num)
			}
		}
		// 注意：这里要引入flag，不然会出错
		if flag {
			ans = append(ans, set)
		}
	}
	return ans
}

// subsetsWithDup_2 递归解法
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(2n)
func subsetsWithDup_2(nums []int) [][]int {
	var (
		ans [][]int
		set []int
		dfs func(choosePre bool, cur int)
	)
	n := len(nums)
	sort.Ints(nums)
	dfs = func(choosePre bool, cur int) {
		if cur == n {
			ans = append(ans, append([]int(nil), set...))
			//下面这句会出bug,原因在于切片是传的地址
			//ans = append(ans, set)
			return
		}
		dfs(false, cur+1)
		// 没有选择上一个数，且当前数与上一个数相同，则跳过
		if !choosePre && cur > 0 && nums[cur-1] == nums[cur] {
			return
		}
		set = append(set, nums[cur])
		dfs(true, cur+1)
		set = set[:len(set)-1]
	}
	dfs(false, 0)
	return ans
}
