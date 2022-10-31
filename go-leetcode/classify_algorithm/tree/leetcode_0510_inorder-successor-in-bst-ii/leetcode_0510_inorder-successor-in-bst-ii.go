package leetcode_0510_inorder_successor_in_bst_ii

// 0510. 二叉搜索树中的中序后继 II
// https://leetcode.cn/problems/inorder-successor-in-bst-ii/

type Node struct {
	Val    int
	Left   *Node
	Right  *Node
	Parent *Node
}

// inorderSuccessor
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路:
//  a.若 node 结点有右孩子，则它的后继在树中相对较低的位置。我们向右走一次，再尽可能的向左走，返回最后所在的结点。
//  b.若 node 结点没有右孩子，则它的后继在树中相对较高的位置。我们向上走到直到结点 tmp 的左孩子是 node 的父节点时，则 node 的后继为 tmp。
func inorderSuccessor(node *Node) *Node {
	if node.Right != nil {
		node = node.Right
		for node.Left != nil {
			node = node.Left
		}
		return node
	}

	for node.Parent != nil && node == node.Parent.Right {
		node = node.Parent
	}
	return node.Parent
}
