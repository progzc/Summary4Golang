package leetcode_0508_most_frequent_subtree_sum

import "math"

// 0508.出现次数最多的子树元素和
// https://leetcode-cn.com/problems/most-frequent-subtree-sum/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// findFrequentTreeSum 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func findFrequentTreeSum(root *TreeNode) []int {
	cm := map[int]int{}
	var (
		ans        []int
		calTreeSum func(root *TreeNode) int
	)
	// calTreeSum 计算树的元素和
	calTreeSum = func(root *TreeNode) int {
		sum := root.Val
		if root.Left != nil {
			sum += calTreeSum(root.Left)
		}
		if root.Right != nil {
			sum += calTreeSum(root.Right)
		}
		cm[sum]++

		return sum
	}
	if root == nil {
		return ans
	}
	calTreeSum(root)
	maxCount := math.MinInt32
	for key, value := range cm {
		if value > maxCount {
			maxCount = value
			ans = []int{key}
		} else if value == maxCount {
			ans = append(ans, key)
		}
	}
	return ans
}

// findFrequentTreeSum 递归优化
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func findFrequentTreeSum_2(root *TreeNode) []int {
	cm := map[int]int{}
	maxCount := math.MinInt32
	var (
		ans        []int
		calTreeSum func(root *TreeNode) int
	)
	// calTreeSum 计算树的元素和
	calTreeSum = func(root *TreeNode) int {
		sum := root.Val
		if root.Left != nil {
			sum += calTreeSum(root.Left)
		}
		if root.Right != nil {
			sum += calTreeSum(root.Right)
		}
		cm[sum]++
		value := cm[sum]
		if value > maxCount {
			maxCount = value
			ans = []int{sum}
		} else if value == maxCount {
			ans = append(ans, sum)
		}
		return sum
	}
	if root == nil {
		return ans
	}
	calTreeSum(root)
	return ans
}
