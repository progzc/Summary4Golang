package leetcode_0114_flatten_binary_tree_to_linked_list

// 0114. 二叉树展开为链表
// https://leetcode.cn/problems/flatten-binary-tree-to-linked-list/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// flatten dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func flatten(root *TreeNode) {
	if root == nil {
		return
	}
	flatten(root.Left)
	flatten(root.Right)
	temp := root.Right
	root.Right = root.Left
	root.Left = nil
	first := root
	for first.Right != nil {
		first = first.Right
	}
	first.Right = temp
}

// flatten bfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func flatten_2(root *TreeNode) {
	if root == nil {
		return
	}
	var (
		stack []*TreeNode
		pre   *TreeNode
	)
	stack = append(stack, root)
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if pre != nil {
			pre.Left = nil
			pre.Right = curr
		}
		if curr.Right != nil {
			stack = append(stack, curr.Right)
		}
		if curr.Left != nil {
			stack = append(stack, curr.Left)
		}
		pre = curr
	}
}
