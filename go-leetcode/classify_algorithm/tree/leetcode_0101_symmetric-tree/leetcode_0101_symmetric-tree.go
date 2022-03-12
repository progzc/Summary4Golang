package leetcode_0101_symmetric_tree

// 0101.对称二叉树
// https://leetcode-cn.com/problems/symmetric-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// isSymmetric_1 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路: 一个树的左子树和右子树镜像对称 可转化为 两个树在什么情况下互为镜像?
//	而两个树互为镜像的条件:
//	a.它们的两个根结点具有相同的值
//	b.每个树的右子树都与另一个树的左子树镜像对称
func isSymmetric_1(root *TreeNode) bool {
	var check func(p, q *TreeNode) bool
	check = func(p, q *TreeNode) bool {
		if p == nil && q == nil {
			return true
		}
		if p == nil || q == nil {
			return false
		}
		return p.Val == q.Val && check(p.Left, q.Right) && check(p.Right, q.Left)
	}
	return check(root, root)
}

// isSymmetric_1 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路: 一个树的左子树和右子树镜像对称 可转化为 两个树在什么情况下互为镜像?
//	而两个树互为镜像的条件:
//	a.它们的两个根结点具有相同的值
//	b.每个树的右子树都与另一个树的左子树镜像对称
func isSymmetric_2(root *TreeNode) bool {
	var (
		stack []*TreeNode
		u, v  *TreeNode
	)
	u, v = root, root
	stack = append(stack, u)
	stack = append(stack, v)
	for len(stack) > 0 {
		u, v = stack[0], stack[1]
		stack = stack[2:]
		if u == nil && v == nil {
			continue
		}
		if u == nil || v == nil {
			return false
		}
		if u.Val != v.Val {
			return false
		}
		stack = append(stack, u.Left)
		stack = append(stack, v.Right)

		stack = append(stack, u.Right)
		stack = append(stack, v.Left)
	}
	return true
}
