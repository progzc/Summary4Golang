package leetcode_0623_add_one_row_to_tree

// 623. 在二叉树中增加一行
// https://leetcode.cn/problems/add-one-row-to-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// addOneRow bfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func addOneRow(root *TreeNode, val int, depth int) *TreeNode {
	if depth == 1 || root == nil {
		return &TreeNode{Val: val, Left: root}
	}

	var stack []*TreeNode
	stack = append(stack, root)
	for i := 0; i < depth-2; i++ {
		n := len(stack)
		for j := 0; j < n; j++ {
			node := stack[0]
			stack = stack[1:]
			if node.Left != nil {
				stack = append(stack, node.Left)
			}
			if node.Right != nil {
				stack = append(stack, node.Right)
			}
		}
	}

	for _, node := range stack {
		node.Left = &TreeNode{Val: val, Left: node.Left}
		node.Right = &TreeNode{Val: val, Right: node.Right}
	}
	return root
}
