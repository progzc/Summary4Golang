package leetcode_0098_validate_binary_search_tree

import "math"

// 验证二叉搜索树
// https://leetcode-cn.com/problems/validate-binary-search-tree/solution/yan-zheng-er-cha-sou-suo-shu-by-leetcode-solution/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// isValidBST_1 递归解法
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func isValidBST_1(root *TreeNode) bool {
	return isValid(root, math.MinInt32, math.MaxInt32)
}

func isValid(root *TreeNode, lower, upper int) bool {
	if root == nil {
		return true
	}
	if root.Val <= lower || root.Val >= upper {
		return false
	}
	return isValid(root.Left, lower, root.Val) && isValid(root.Right, root.Val, upper)
}

// isValidBST_2 中序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func isValidBST_2(root *TreeNode) bool {
	var stack []*TreeNode
	inorder := math.MinInt64
	for len(stack) > 0 || root != nil {
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}
		root = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if root.Val <= inorder {
			return false
		}
		inorder = root.Val
		root = root.Right
	}
	return true
}
