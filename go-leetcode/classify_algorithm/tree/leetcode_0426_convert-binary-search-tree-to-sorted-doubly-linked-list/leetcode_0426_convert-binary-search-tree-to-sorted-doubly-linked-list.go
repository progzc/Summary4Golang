package leetcode_0426_convert_binary_search_tree_to_sorted_doubly_linked_list

// 0426. 将二叉搜索树转化为排序的双向链表
// https://leetcode.cn/problems/convert-binary-search-tree-to-sorted-doubly-linked-list/

type Node struct {
	Val   int
	Left  *Node
	Right *Node
}

// treeToDoublyList 中序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func treeToDoublyList(root *Node) *Node {
	if root == nil {
		return nil
	}
	var (
		dfs   func(node *Node)
		first *Node
		last  *Node
	)
	dfs = func(node *Node) {
		if node == nil {
			return
		}
		dfs(node.Left)
		if last != nil {
			last.Right = node
			node.Left = last
		} else {
			first = node
		}
		last = node
		dfs(node.Right)
	}
	dfs(root)
	last.Right = first
	first.Left = last
	return first
}
