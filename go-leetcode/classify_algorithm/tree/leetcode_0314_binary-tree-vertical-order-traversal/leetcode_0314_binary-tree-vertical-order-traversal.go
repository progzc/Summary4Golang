package leetcode_0314_binary_tree_vertical_order_traversal

import "math"

// 0314. 二叉树的垂直遍历
// https://leetcode.cn/problems/binary-tree-vertical-order-traversal/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// verticalOrder 层序遍历+散列表
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func verticalOrder(root *TreeNode) [][]int {
	var (
		ans      [][]int
		queue    []*TreeNode // 层序遍历的节点队列
		posQueue []int       // 层序遍历节点对应的位置队列
		m        = map[int][]int{}
		minPos   = math.MaxInt32 // 最左侧位置
	)

	if root == nil {
		return ans
	}

	queue = append(queue, root)
	posQueue = append(posQueue, 0)
	for len(queue) > 0 {
		node, pos := queue[0], posQueue[0]
		queue, posQueue = queue[1:], posQueue[1:]
		if v, ok := m[pos]; ok {
			v = append(v, node.Val)
			m[pos] = v
		} else {
			m[pos] = []int{node.Val}
		}

		if node.Left != nil {
			queue = append(queue, node.Left)
			posQueue = append(posQueue, pos-1)
		}

		if node.Right != nil {
			queue = append(queue, node.Right)
			posQueue = append(posQueue, pos+1)
		}

		minPos = min(minPos, pos)
	}

	for i := minPos; i < minPos+len(m); i++ {
		ans = append(ans, m[i])
	}

	return ans
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
