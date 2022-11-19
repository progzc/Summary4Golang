package leetcode_0145_binary_tree_postorder_traversal

// 145. 二叉树的后序遍历
// https://leetcode.cn/problems/binary-tree-postorder-traversal/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// postorderTraversal 递归法
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func postorderTraversal(root *TreeNode) []int {
	var (
		ans []int
		dfs func(root *TreeNode)
	)
	dfs = func(root *TreeNode) {
		if root == nil {
			return
		}

		dfs(root.Left)
		dfs(root.Right)
		ans = append(ans, root.Val)
	}
	dfs(root)
	return ans
}

// postorderTraversal_2 迭代法
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func postorderTraversal_2(root *TreeNode) []int {
	var (
		ans   []int
		stack []*TreeNode
	)
	if root == nil {
		return ans
	}
	stack = append(stack, root)
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		ans = append(ans, node.Val)
		if node.Left != nil {
			stack = append(stack, node.Left)
		}
		if node.Right != nil {
			stack = append(stack, node.Right)
		}
	}
	reverse(ans)
	return ans
}

func reverse(nums []int) {
	n := len(nums)
	for i := 0; i < n/2; i++ {
		nums[i], nums[n-i-1] = nums[n-i-1], nums[i]
	}
}
