package leetcode_0110_balanced_binary_tree

// 0110.平衡二叉树
// https://leetcode-cn.com/problems/balanced-binary-tree/
// 关于高度平衡的二叉树的定义：一个二叉树每个节点 的左右两个子树的高度差的绝对值不超过 1

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// isBalanced 递归
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func isBalanced(root *TreeNode) bool {
	if root == nil {
		return true
	}

	var (
		max func(x, y int) int
		abs func(x int) int
		cal func(root *TreeNode) int
	)

	max = func(x, y int) int {
		if x > y {
			return x
		}
		return y
	}

	abs = func(x int) int {
		if x < 0 {
			return -1 * x
		}
		return x
	}

	cal = func(root *TreeNode) int {
		if root == nil {
			return 0
		}
		return max(cal(root.Left), cal(root.Right)) + 1
	}

	return abs(cal(root.Left)-cal(root.Right)) <= 1 && isBalanced(root.Left) && isBalanced(root.Right)
}
