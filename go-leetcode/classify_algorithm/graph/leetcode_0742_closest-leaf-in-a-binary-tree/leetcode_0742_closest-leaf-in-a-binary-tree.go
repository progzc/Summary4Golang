package leetcode_0742_closest_leaf_in_a_binary_tree

// 0742. 二叉树最近的叶节点
// https://leetcode.cn/problems/closest-leaf-in-a-binary-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// findClosestLeaf 转换成图+bfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func findClosestLeaf(root *TreeNode, k int) int {
	g := make(map[*TreeNode][]*TreeNode)
	var dfs func(node *TreeNode, parent *TreeNode)
	dfs = func(node *TreeNode, parent *TreeNode) {
		if node == nil {
			return
		}
		g[node] = append(g[node], parent)
		g[parent] = append(g[parent], node)

		dfs(node.Left, node)
		dfs(node.Right, node)
	}
	dfs(root, nil)

	var (
		stack []*TreeNode
		seen  = make(map[*TreeNode]bool)
	)
	for node := range g {
		if node != nil && node.Val == k {
			stack = append(stack, node)
			seen[node] = true
		}
	}

	for len(stack) > 0 {
		node := stack[0]
		stack = stack[1:]
		if node == nil {
			continue
		}
		if len(g[node]) == 1 {
			return node.Val
		}
		for _, nextNode := range g[node] {
			if !seen[nextNode] {
				seen[nextNode] = true
				stack = append(stack, nextNode)
			}
		}
	}
	return -1
}
