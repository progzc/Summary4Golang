package leetcode_0250_count_univalue_subtrees

// 250.统计同值子树
// https://leetcode-cn.com/problems/count-univalue-subtrees/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// countUnivalSubtrees 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func countUnivalSubtrees(root *TreeNode) int {
	var (
		ans         int
		isEqualTree func(node *TreeNode) bool
	)
	// isEqualTree 判断一个子树是否是同值子树
	isEqualTree = func(node *TreeNode) bool {
		if node == nil {
			return false
		}
		if node.Left != nil && node.Val != node.Left.Val {
			return false
		}
		if node.Right != nil && node.Val != node.Right.Val {
			return false
		}
		return (node.Left == nil || isEqualTree(node.Left)) && (node.Right == nil || isEqualTree(node.Right))
	}

	if isEqualTree(root) {
		ans++
	}
	if root != nil && root.Left != nil {
		ans += countUnivalSubtrees(root.Left)
	}
	if root != nil && root.Right != nil {
		ans += countUnivalSubtrees(root.Right)
	}
	return ans
}

// countUnivalSubtrees_2 递归优化
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func countUnivalSubtrees_2(root *TreeNode) int {
	var (
		ans         int
		isEqualTree func(node *TreeNode) bool
	)
	// isEqualTree 判断一个子树是否是同值子树
	isEqualTree = func(node *TreeNode) bool {
		if node.Left == nil && node.Right == nil {
			ans++
			return true
		}
		isUnival := true
		if node.Left != nil {
			isUnival = isEqualTree(node.Left) && node.Left.Val == node.Val && isUnival
		}
		if node.Right != nil {
			isUnival = isEqualTree(node.Right) && node.Right.Val == node.Val && isUnival
		}
		if !isUnival {
			return false
		}
		ans++
		return true
	}
	if root == nil {
		return 0
	}
	isEqualTree(root)
	return ans
}
