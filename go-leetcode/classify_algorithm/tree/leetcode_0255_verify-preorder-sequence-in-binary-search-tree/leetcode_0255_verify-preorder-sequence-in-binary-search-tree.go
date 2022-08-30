package leetcode_0255_verify_preorder_sequence_in_binary_search_tree

import "math"

// 0255. 验证前序遍历序列二叉搜索树
// https://leetcode.cn/problems/verify-preorder-sequence-in-binary-search-tree/

// verifyPreorder 单调栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func verifyPreorder(preorder []int) bool {
	var (
		// 单调递减栈
		stack  []int
		preEle = math.MinInt32
	)

	for i := 0; i < len(preorder); i++ {
		// 右子树元素必须要大于递减栈被peek访问的元素,否则就不是二叉搜索树
		if preorder[i] < preEle {
			return false
		}
		for len(stack) > 0 && preorder[i] > stack[len(stack)-1] {
			// 数组元素大于单调栈的元素了，表示往右子树走了，记录下上个根节点
			// 找到这个右子树对应的根节点，之前左子树全部弹出，不在记录，因为不可能在往根节点的左子树走了
			preEle = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, preorder[i])
	}
	return true
}
