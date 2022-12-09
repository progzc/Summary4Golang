package leetcode_0337_house_robber_iii

// 337. 打家劫舍 III
// https://leetcode.cn/problems/house-robber-iii/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// rob dfs 超时
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func rob(root *TreeNode) int {
	var (
		dfs func(root *TreeNode, pre bool) int
	)

	dfs = func(root *TreeNode, pre bool) int {
		if root == nil {
			return 0
		}

		var ans int
		if !pre {
			// 父节点未偷盗
			// 选择偷盗本节点
			ans = max(ans, root.Val+dfs(root.Left, true)+dfs(root.Right, true))
			// 不选择偷盗本节点
			ans = max(ans, dfs(root.Left, false)+dfs(root.Right, false))
		} else {
			// 父节点已偷盗
			// 不选择偷盗本节点
			ans = max(ans, dfs(root.Left, false)+dfs(root.Right, false))
		}
		return ans
	}
	return dfs(root, false)
}

// rob_2 dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func rob_2(root *TreeNode) int {
	var (
		dfs func(root *TreeNode) []int
	)

	// 以root为起始节点,返回分别选择或不选择偷盗root节点的最大偷盗金额
	dfs = func(root *TreeNode) []int {
		if root == nil {
			return []int{0, 0}
		}
		l, r := dfs(root.Left), dfs(root.Right)
		selected := root.Val + l[1] + r[1]
		noSelected := max(l[0], l[1]) + max(r[0], r[1])
		return []int{selected, noSelected}
	}
	chooses := dfs(root)
	return max(chooses[0], chooses[1])
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
