package leetcode_0572_subtree_of_another_tree

// 0572.另一棵树的子树
// https://leetcode-cn.com/problems/subtree-of-another-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// isSubtree 双重递归
// 时间复杂度: O(m*n)
// 空间复杂度: O(n)
// 思路：
//	依次将各个子树与subRoot进行比较，如果有一个相同，则subRoot是该树的子树
func isSubtree(root *TreeNode, subRoot *TreeNode) bool {
	if root == nil {
		return false
	}
	// check 判断两个树是否相等
	var check func(x, y *TreeNode) bool
	check = func(x, y *TreeNode) bool {
		if x == nil && y == nil {
			return true
		}
		if x == nil || y == nil {
			return false
		}
		if x.Val == y.Val {
			return check(x.Left, y.Left) && check(x.Right, y.Right)
		}
		return false
	}
	return check(root, subRoot) || isSubtree(root.Left, subRoot) || isSubtree(root.Right, subRoot)
}

// isSubtree_wrong 错误的解法
func isSubtree_wrong(root *TreeNode, subRoot *TreeNode) bool {
	if root == nil && subRoot == nil {
		return true
	}
	if root == nil || subRoot == nil {
		return false
	}

	if root.Val == subRoot.Val {
		return isSubtree_wrong(root.Left, subRoot.Left) && isSubtree_wrong(root.Right, subRoot.Right)
	}
	return isSubtree_wrong(root.Left, subRoot) || isSubtree_wrong(root.Right, subRoot)
}
