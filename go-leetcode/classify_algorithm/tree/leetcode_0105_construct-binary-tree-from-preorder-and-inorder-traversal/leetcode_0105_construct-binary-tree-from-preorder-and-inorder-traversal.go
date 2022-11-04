package leetcode_0105_construct_binary_tree_from_preorder_and_inorder_traversal

// 0105. 从前序与中序遍历序列构造二叉树
// https://leetcode.cn/problems/construct-binary-tree-from-preorder-and-inorder-traversal/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// buildTree dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func buildTree(preorder []int, inorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}

	root := &TreeNode{
		Val: preorder[0],
	}
	var idx int
	for i, num := range inorder {
		if num == preorder[0] {
			idx = i
			break
		}
	}
	root.Left = buildTree(preorder[1:idx+1], inorder[:idx])
	root.Right = buildTree(preorder[idx+1:], inorder[idx+1:])
	return root
}
