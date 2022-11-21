package leetcode_0106_construct_binary_tree_from_inorder_and_postorder_traversal

// 106. 从中序与后序遍历序列构造二叉树
// https://leetcode.cn/problems/construct-binary-tree-from-inorder-and-postorder-traversal/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// buildTree dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func buildTree(inorder []int, postorder []int) *TreeNode {
	n := len(postorder)
	if n == 0 {
		return nil
	}
	num := postorder[n-1]
	root := &TreeNode{Val: num}

	idx := -1
	for i := 0; i < len(inorder); i++ {
		if inorder[i] == num {
			idx = i
			break
		}
	}
	root.Left = buildTree(inorder[:idx], postorder[:idx])
	root.Right = buildTree(inorder[idx+1:], postorder[idx:n-1])
	return root
}
