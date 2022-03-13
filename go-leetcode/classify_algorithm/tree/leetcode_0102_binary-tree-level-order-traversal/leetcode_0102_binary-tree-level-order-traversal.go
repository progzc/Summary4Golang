package leetcode_0102_binary_tree_level_order_traversal

// 0102.二叉树的层序遍历
// https://leetcode-cn.com/problems/binary-tree-level-order-traversal/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// levelOrder 迭代(BFS)
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func levelOrder(root *TreeNode) [][]int {
	var (
		queue []*TreeNode
		ans   [][]int
	)
	if root == nil {
		return ans
	}
	queue = append(queue, root)
	for len(queue) > 0 {
		var layer []int
		size := len(queue)
		for i := 0; i < size; i++ {
			temp := queue[0]
			queue = queue[1:]
			layer = append(layer, temp.Val)
			if temp.Left != nil {
				queue = append(queue, temp.Left)
			}
			if temp.Right != nil {
				queue = append(queue, temp.Right)
			}
		}
		ans = append(ans, layer)
	}
	return ans
}
