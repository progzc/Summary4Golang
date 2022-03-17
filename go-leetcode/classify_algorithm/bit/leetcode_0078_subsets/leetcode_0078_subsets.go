package leetcode_0078_subsets

// 0078.子集
// https://leetcode-cn.com/problems/subsets/

// subsets 位操作
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(n)
func subsets(nums []int) [][]int {
	ans := [][]int{}
	n := len(nums)
	for mask := 0; mask < 1<<n; mask++ {
		set := []int{}
		for i, num := range nums {
			if mask&(1<<i) > 0 {
				set = append(set, num)
			}
		}
		ans = append(ans, set)
	}
	return ans
}

// subsets_2 递归
// 时间复杂度: O(n)
// 空间复杂度: O()
func subsets_2(nums []int) [][]int {
	var (
		ans [][]int
		set []int
		dfs func(cur int)
	)
	n := len(nums)
	dfs = func(cur int) {
		if cur == n {
			ans = append(ans, append([]int(nil), set...))
			//下面这句会出bug,原因在于切片是传的地址
			//ans = append(ans, set) // 输入:[1,2,3]; 输出:[[3,3,3],[3,3],[3,3],[3],[3,3],[3],[3],[]]
			return
		}
		set = append(set, nums[cur])
		dfs(cur + 1)
		set = set[:len(set)-1]
		dfs(cur + 1)
	}
	dfs(0)
	return ans
}
