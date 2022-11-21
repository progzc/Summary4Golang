package leetcode_0108_convert_sorted_array_to_binary_search_tree

// 108. 将有序数组转换为二叉搜索树
// https://leetcode.cn/problems/convert-sorted-array-to-binary-search-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// sortedArrayToBST
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func sortedArrayToBST(nums []int) *TreeNode {
	n := len(nums)
	if n == 0 {
		return nil
	}
	root := &TreeNode{
		Val: nums[n/2],
	}
	root.Left = sortedArrayToBST(nums[:n/2])
	root.Right = sortedArrayToBST(nums[n/2+1:])
	return root
}
