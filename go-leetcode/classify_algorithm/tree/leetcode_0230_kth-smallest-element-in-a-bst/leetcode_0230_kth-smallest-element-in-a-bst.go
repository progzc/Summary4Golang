package leetcode_0230_kth_smallest_element_in_a_bst

// 0230. 二叉搜索树中第 K 小的元素
// https://leetcode.cn/problems/kth-smallest-element-in-a-bst

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// kthSmallest 中序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func kthSmallest(root *TreeNode, k int) int {
	var (
		dfs func(root *TreeNode)
		ans []int
	)
	dfs = func(root *TreeNode) {
		if root == nil {
			return
		}
		dfs(root.Left)
		ans = append(ans, root.Val)
		dfs(root.Right)
	}
	dfs(root)
	return ans[k-1]
}

// kthSmallest_2 中序遍历(使用递归法无须遍历整颗树)
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func kthSmallest_2(root *TreeNode, k int) int {
	var (
		dfs func(root *TreeNode)
		ans int
	)
	dfs = func(root *TreeNode) {
		if root == nil {
			return
		}
		dfs(root.Left)
		k--
		if k == 0 {
			ans = root.Val
			return
		}
		dfs(root.Right)
	}
	dfs(root)
	return ans
}

// kthSmallest_3 中序遍历(使用迭代法无须遍历整颗树)
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func kthSmallest_3(root *TreeNode, k int) int {
	// TODO
	return 0
}
