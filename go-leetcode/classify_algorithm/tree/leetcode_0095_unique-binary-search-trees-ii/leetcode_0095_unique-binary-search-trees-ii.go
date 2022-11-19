package leetcode_0095_unique_binary_search_trees_ii

// 95. 不同的二叉搜索树 II
// https://leetcode.cn/problems/unique-binary-search-trees-ii/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// generateTrees 递归
// 时间复杂度:
// 空间复杂度:
func generateTrees(n int) []*TreeNode {
	if n == 0 {
		return nil
	}

	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = i + 1
	}
	var dfs func(nums []int) []*TreeNode
	dfs = func(nums []int) []*TreeNode {
		var ans []*TreeNode
		if len(nums) == 0 {
			return []*TreeNode{nil}
		}
		for i := 0; i < len(nums); i++ {
			// 左子树集合
			lefts := dfs(nums[:i])
			// 右子树集合
			rights := dfs(nums[i+1:])
			// 组装成树
			for _, left := range lefts {
				for _, right := range rights {
					root := &TreeNode{Val: nums[i]}
					root.Left = left
					root.Right = right
					ans = append(ans, root)
				}
			}
		}
		return ans
	}
	return dfs(nums)
}
