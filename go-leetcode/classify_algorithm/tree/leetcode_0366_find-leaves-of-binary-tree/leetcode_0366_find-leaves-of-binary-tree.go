package leetcode_0366_find_leaves_of_binary_tree

// 0366. 寻找二叉树的叶子节点
// https://leetcode.cn/problems/find-leaves-of-binary-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// findLeaves 深度优先遍历+后序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func findLeaves(root *TreeNode) [][]int {
	var (
		ans [][]int
		dfs func(root *TreeNode) int
	)

	dfs = func(root *TreeNode) int {
		if root == nil {
			return -1
		}

		left := dfs(root.Left)
		right := dfs(root.Right)
		cur := max(left, right) + 1
		if cur >= len(ans) {
			ans = append(ans, []int{})
		}
		ans[cur] = append(ans[cur], root.Val)
		return cur
	}

	dfs(root)
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
