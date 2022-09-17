package leetcode_1373_maximum_sum_bst_in_binary_tree

import "math"

// 1373. 二叉搜索子树的最大键值和
// https://leetcode.cn/problems/maximum-sum-bst-in-binary-tree/
// 注意事项: 空的子树也是BST，而空树的键值之和为0

// 同下面题
// 0333. 最大 BST 子树
// https://leetcode.cn/problems/largest-bst-subtree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// maxSumBST 层次遍历(常规解法)
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
// 缺点：会超出时间限制
func maxSumBST(root *TreeNode) int {
	if root == nil {
		return 0
	}

	var (
		sum   func(node *TreeNode) int
		valid func(node *TreeNode, l, r int) bool
		ans   int
	)

	sum = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		return node.Val + sum(node.Left) + sum(node.Right)
	}

	valid = func(node *TreeNode, l, r int) bool {
		if node == nil {
			return true
		}
		if node.Val <= l || node.Val >= r {
			return false
		}
		return valid(node.Left, l, node.Val) && valid(node.Right, node.Val, r)
	}

	if valid(root, math.MinInt32, math.MaxInt32) {
		// 注意,与【0333. 最大 BST 子树】不同,这里不能直接return
		// 因为节点值可为负数
		ans = max(ans, sum(root))
	}
	ans = max(ans, max(maxSumBST(root.Left), maxSumBST(root.Right)))
	return ans
}

// maxSumBST_2 后序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 缺点：会超出时间限制
func maxSumBST_2(root *TreeNode) int {
	if root == nil {
		return 0
	}

	var (
		dfs func(node *TreeNode) (l, r, sum int, isBst bool)
		ans int
	)

	// l,r表示当前节点为根的二叉搜索树里的值的范围[l,r]，isBst为这棵树是否为二叉树，
	// 如果不是BST,则isBst == false，还未递归前l == r == node.val
	dfs = func(node *TreeNode) (l, r, sum int, isBst bool) {
		if node.Left == nil && node.Right == nil {
			ans = max(ans, node.Val)
			return node.Val, node.Val, node.Val, true
		}

		l, r, sum, isBst = node.Val, node.Val, node.Val, true
		if node.Left != nil {
			l1, r1, sum1, isBst1 := dfs(node.Left)
			if isBst1 && node.Val > r1 {
				sum += sum1
				l = l1
			} else {
				isBst = false
			}
		}

		if node.Right != nil {
			l2, r2, sum2, isBst2 := dfs(node.Right)
			if isBst2 && node.Val < l2 {
				sum += sum2
				r = r2
			} else {
				isBst = false
			}
		}

		// 是BST
		if isBst {
			ans = max(ans, sum)
			return
		}
		// 不是BST
		return
	}
	dfs(root)
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
