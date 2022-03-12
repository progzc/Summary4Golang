package leetcode_0104_maximum_depth_of_binary_tree

// 0104.二叉树的最大深度
// https://leetcode-cn.com/problems/maximum-depth-of-binary-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// maxDepth_1 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func maxDepth_1(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return max(maxDepth_1(root.Left), maxDepth_1(root.Right)) + 1
}

// maxDepth_2 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func maxDepth_2(root *TreeNode) int {
	if root == nil {
		return 0
	}
	var (
		queue []*TreeNode
		ans   int
	)

	queue = append(queue, root)
	for len(queue) > 0 {
		size := len(queue)
		for size > 0 {
			node := queue[0]
			queue = queue[1:]
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
			size--
		}
		ans++
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
