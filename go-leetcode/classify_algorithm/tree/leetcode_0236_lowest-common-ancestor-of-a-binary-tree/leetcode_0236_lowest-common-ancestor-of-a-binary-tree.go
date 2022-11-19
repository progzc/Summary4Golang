package leetcode_0236_lowest_common_ancestor_of_a_binary_tree

// 0236.二叉树的最近公共祖先
// https://leetcode-cn.com/problems/lowest-common-ancestor-of-a-binary-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// lowestCommonAncestor 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：
//	定义：f(x)表示x节点的子树中是否包含p节点或q节点,包含则为true，否则为false
//	那么：最近公共祖先一定满足如下条件
//		(f(x_L) && f(x_R)) || ((x==p||x==q) && (f(x_L)||f(x_R)))
//	解释：一共只有如下4种情况
//	f(x_L) && f(x_R)：
//		a.x的左子树包含p,且x的右子树包含q，则x即为最近公共祖先
//		b.x的左子树包含q,且x的右子树包含p，则x即为最近公共祖先
//	(x==p||x==q) && (f(x_L)||f(x_R)：
//		c.x为p节点，且x的左或右子树包含q，则p即为最近公共祖先
//		d.x为q节点，且x的左或右子树包含p，则q即为最近公共祖先
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	if root.Val == p.Val || root.Val == q.Val {
		return root
	}

	left := lowestCommonAncestor(root.Left, p, q)
	right := lowestCommonAncestor(root.Right, p, q)
	if left != nil && right != nil {
		return root
	}
	if left == nil {
		return right
	}
	return left
}
