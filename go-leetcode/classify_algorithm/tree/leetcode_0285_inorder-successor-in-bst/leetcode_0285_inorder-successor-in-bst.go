package leetcode_0285_inorder_successor_in_bst

// 0285. 二叉搜索树中的中序后继
// https://leetcode.cn/problems/inorder-successor-in-bst/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// inorderSuccessor dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func inorderSuccessor(root *TreeNode, p *TreeNode) *TreeNode {
	find := false
	var (
		ans *TreeNode
		dfs func(root *TreeNode)
	)
	dfs = func(root *TreeNode) {
		// ans != nil 表示如果找到结果就提前结束
		if root == nil || ans != nil {
			return
		}
		dfs(root.Left)
		if find && ans == nil {
			ans = root
		}
		if root.Val == p.Val && !find {
			find = true
		}
		dfs(root.Right)
	}
	dfs(root)
	return ans
}

// inorderSuccessor_2 bfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func inorderSuccessor_2(root *TreeNode, p *TreeNode) *TreeNode {
	var (
		stack []*TreeNode
		pre   *TreeNode
		cur   = root
	)

	for cur != nil || len(stack) > 0 {
		// 遍历左节点
		for cur != nil {
			stack = append(stack, cur)
			cur = cur.Left
		}
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if pre == p {
			return cur
		}
		pre = cur
		cur = cur.Right
	}
	return nil
}
