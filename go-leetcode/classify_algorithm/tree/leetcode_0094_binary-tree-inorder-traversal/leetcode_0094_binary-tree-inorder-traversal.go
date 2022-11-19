package leetcode_0094_binary_tree_inorder_traversal

// 0094.二叉树的中序遍历
// https://leetcode-cn.com/problems/binary-tree-inorder-traversal/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// inorderTraversal_1 递归法
// 时间复杂度: O(n)
// 空间复杂度: O(log(n))，当二叉树退化为一条链时空间复杂度最差,为O(n)
func inorderTraversal_1(root *TreeNode) []int {
	var (
		ans []int
		dfs func(root *TreeNode)
	)
	dfs = func(root *TreeNode) {
		if root == nil {
			return
		}
		dfs(root.Left)
		ans = append(ans, root.Val)
		dfs(root.Right)
	}
	dfs(root)
	return ans
}

// inorderTraversal_2 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func inorderTraversal_2(root *TreeNode) []int {
	var (
		stack []*TreeNode
		ans   []int
	)
	for root != nil || len(stack) > 0 {
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}
		root = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		ans = append(ans, root.Val)
		root = root.Right
	}
	return ans
}
