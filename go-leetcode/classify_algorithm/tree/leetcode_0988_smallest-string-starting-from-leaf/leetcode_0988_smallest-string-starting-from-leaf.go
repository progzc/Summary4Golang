package leetcode_0988_smallest_string_starting_from_leaf

// 988. 从叶结点开始的最小字符串
// https://leetcode.cn/problems/smallest-string-starting-from-leaf/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// smallestFromLeaf dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func smallestFromLeaf(root *TreeNode) string {
	var (
		dfs func(root *TreeNode, str string)
		ans string
	)
	dfs = func(root *TreeNode, str string) {
		if root == nil {
			return
		}
		str = string(byte(root.Val)+'a') + str
		if root.Left == nil && root.Right == nil {
			if len(ans) == 0 {
				ans = str
			} else if str < ans {
				ans = str
			}
			return
		}
		if root.Left != nil {
			dfs(root.Left, str)
		}
		if root.Right != nil {
			dfs(root.Right, str)
		}
	}
	dfs(root, "")
	return ans
}
