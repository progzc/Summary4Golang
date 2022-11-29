package leetcode_0129_sum_root_to_leaf_numbers

// 129. 求根节点到叶节点数字之和
// https://leetcode.cn/problems/sum-root-to-leaf-numbers/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// sumNumbers dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func sumNumbers(root *TreeNode) int {
	var (
		ans int
		dfs func(root *TreeNode, cur int)
	)

	dfs = func(root *TreeNode, cur int) {
		if root == nil {
			return
		}
		cur = cur*10 + root.Val
		if root.Left == nil && root.Right == nil {
			ans += cur
		}
		if root.Left != nil {
			dfs(root.Left, cur)
		}
		if root.Right != nil {
			dfs(root.Right, cur)
		}
	}

	dfs(root, 0)
	return ans
}
