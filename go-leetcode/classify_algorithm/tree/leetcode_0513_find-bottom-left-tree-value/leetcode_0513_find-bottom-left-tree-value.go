package leetcode_0513_find_bottom_left_tree_value

// 513. 找树左下角的值
// https://leetcode.cn/problems/find-bottom-left-tree-value/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// findBottomLeftValue 迭代法
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func findBottomLeftValue(root *TreeNode) int {
	var (
		ans   int
		stack []*TreeNode
	)
	if root == nil {
		return ans
	}
	stack = append(stack, root)
	for len(stack) > 0 {
		ans = stack[0].Val
		n := len(stack)
		for i := 0; i < n; i++ {
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
	return ans
}
