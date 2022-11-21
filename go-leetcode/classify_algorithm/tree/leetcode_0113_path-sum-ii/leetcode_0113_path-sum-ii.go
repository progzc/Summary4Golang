package leetcode_0113_path_sum_ii

// 113. 路径总和 II
// https://leetcode.cn/problems/path-sum-ii/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// pathSum 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func pathSum(root *TreeNode, targetSum int) [][]int {
	var (
		ans  [][]int
		path []int
		dfs  func(root *TreeNode, targetSum int)
	)

	dfs = func(root *TreeNode, sum int) {
		if root == nil {
			return
		}
		sum += root.Val
		path = append(path, root.Val)
		defer func() {
			path = path[:len(path)-1]
		}()
		if root.Left == nil && root.Right == nil {
			if sum == targetSum {
				ans = append(ans, append([]int(nil), path...))
				return
			}
		}
		dfs(root.Left, sum)
		dfs(root.Right, sum)
	}
	dfs(root, 0)
	return ans
}
