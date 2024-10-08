package leetcode_0437_path_sum_iii

// 437. 路径总和 III🌟
// https://leetcode.cn/problems/path-sum-iii/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// pathSum
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func pathSum(root *TreeNode, targetSum int) int {
	var (
		ans int
		dfs func(root *TreeNode, targetSum int) int
	)

	if root == nil {
		return ans
	}

	// dfs 以root节点为起始点出发，路径总和为targetSum的路径数量
	dfs = func(root *TreeNode, targetSum int) int {
		var cnt int
		if root == nil {
			return 0
		}
		if root.Val == targetSum {
			cnt++
		}
		cnt += dfs(root.Left, targetSum-root.Val)
		cnt += dfs(root.Right, targetSum-root.Val)
		return cnt
	}
	ans = dfs(root, targetSum)
	// 左子树的路径总和的路径条数
	ans += pathSum(root.Left, targetSum)
	// 右子树的路径总和的路径条数
	ans += pathSum(root.Right, targetSum)
	return ans
}
