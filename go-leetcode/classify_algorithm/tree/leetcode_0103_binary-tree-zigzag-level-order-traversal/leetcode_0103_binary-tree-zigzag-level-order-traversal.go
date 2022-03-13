package leetcode_0103_binary_tree_zigzag_level_order_traversal

// 0103.二叉树的锯齿形层序遍历
// https://leetcode-cn.com/problems/binary-tree-zigzag-level-order-traversal/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// zigzagLevelOrder 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：先按层序遍历正常迭代,碰到特殊的层两两交换倒序即可
func zigzagLevelOrder(root *TreeNode) [][]int {
	var (
		queue []*TreeNode
		ans   [][]int
		level int
	)
	if root == nil {
		return ans
	}
	queue = append(queue, root)
	for len(queue) > 0 {
		layer := []int{}
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
		// 本质上和层序遍历一样，我们只需要把奇数层的元素翻转即可
		if level%2 == 1 {
			n := len(layer)
			for i := 0; i < n/2; i++ {
				layer[i], layer[n-1-i] = layer[n-1-i], layer[i]
			}
		}
		ans = append(ans, layer)
		level++
	}
	return ans
}

// zigzagLevelOrder_1 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：使用标志位来控制是插入顺序（尾插 或 头插）,而不需要交换
func zigzagLevelOrder_2(root *TreeNode) [][]int {
	var (
		queue []*TreeNode
		ans   [][]int
	)
	if root == nil {
		return ans
	}
	queue = append(queue, root)
	flag := false
	for len(queue) > 0 {
		layer := []int{}
		size := len(queue)
		for i := 0; i < size; i++ {
			temp := queue[0]
			queue = queue[1:]
			if !flag {
				layer = append(layer, temp.Val) // 尾插
			} else {
				layer = append([]int{temp.Val}, layer...) // 头插
			}

			if temp.Left != nil {
				queue = append(queue, temp.Left)
			}
			if temp.Right != nil {
				queue = append(queue, temp.Right)
			}
		}
		ans = append(ans, layer)
		flag = !flag
	}
	return ans
}
