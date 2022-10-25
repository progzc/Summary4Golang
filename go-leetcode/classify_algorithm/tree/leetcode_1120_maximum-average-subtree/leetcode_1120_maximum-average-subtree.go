package leetcode_1120_maximum_average_subtree

// 1120. 子树的最大平均值
// https://leetcode.cn/problems/maximum-average-subtree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// maximumAverageSubtree dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func maximumAverageSubtree(root *TreeNode) float64 {
	if root == nil {
		return float64(0)
	}
	var ans float64
	var dfs func(cur *TreeNode) (int, int)
	dfs = func(cur *TreeNode) (int, int) {
		if cur == nil {
			return 0, 0
		}

		var sum, cnt int
		if cur.Left != nil {
			s1, c1 := dfs(cur.Left)
			if c1 > 0 {
				sum += s1
				cnt += c1
				ans = max(ans, float64(s1)/float64(c1))
			}
		}
		if cur.Right != nil {
			s2, c2 := dfs(cur.Right)
			if c2 > 0 {
				sum += s2
				cnt += c2
				ans = max(ans, float64(s2)/float64(c2))
			}
		}
		ans = max(ans, float64(sum+cur.Val)/float64(cnt+1))
		return sum + cur.Val, cnt + 1
	}
	dfs(root)
	return ans
}

func max(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}
