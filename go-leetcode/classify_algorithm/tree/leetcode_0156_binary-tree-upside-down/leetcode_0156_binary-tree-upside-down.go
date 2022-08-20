package leetcode_0156_binary_tree_upside_down

// 0156.上下翻转二叉树
// https://leetcode.cn/problems/binary-tree-upside-down/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// upsideDownBinaryTree 层序遍历
// 时间复杂度: O(n),n为数的高度
// 空间复杂度: O(1)
func upsideDownBinaryTree(root *TreeNode) *TreeNode {
	var right, father *TreeNode
	for root != nil {
		// 为了继续遍历,先记录下原来的左子节点防止丢失
		left := root.Left
		// 当前节点的左子节点更新为父节点的右子节点
		root.Left = right

		// 记录下当前节点的右子节点
		right = root.Right
		// 当前节点的右子节点更新为原父节点
		root.Right = father

		// 记录下当前节点作为下一个待遍历节点的父节点
		father = root
		// 继续下一层
		root = left
	}
	return father
}
