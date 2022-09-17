package leetcode_0333_largest_bst_subtree

import "math"

// 0333. 最大 BST 子树
// https://leetcode.cn/problems/largest-bst-subtree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// largestBSTSubtree 中序遍历(常规解法)
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func largestBSTSubtree(root *TreeNode) int {
	if root == nil {
		return 0
	}

	var (
		counter func(node *TreeNode) int
		valid   func(node *TreeNode, l, r int) bool
	)

	counter = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		return counter(node.Left) + counter(node.Right) + 1
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
		return counter(root)
	}
	return max(largestBSTSubtree(root.Left), largestBSTSubtree(root.Right))
}

// largestBSTSubtree_2 中序遍历(优化解法)
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：一棵树如果是二叉搜索树，那么它的左右子树也必然是二叉搜索树，则对于一个节点为根的子树，
// 如果我们已经知道了左右子树是不是二叉搜索树，以及左右子树的值的范围 [l,r]，那么如果左右子树均为二叉搜索树，
// 根据性质我们只要判断该节点的值是不是 大于左子树的最大值 和 小于右子树的最小值
// 即能推断出该节点为根的子树是不是二叉搜索树， 而又因为我们已经拿到了左右子树的信息，所以这个推断只需要 O(1) 的时间复杂度，
// 而方法一不复用信息的话判断一棵子树是不是二叉搜索树则需要 O(n) 的时间复杂度。
func largestBSTSubtree_2(root *TreeNode) int {
	if root == nil {
		return 0
	}

	var (
		dfs func(node *TreeNode) (l, r, sz int)
		ans int
	)

	// l,r表示当前节点为根的二叉搜索树里的值的范围[l,r]，sz为这棵树的节点数，
	// 如果不是BST,则sz == -1，还未递归前l == r == node.val
	dfs = func(node *TreeNode) (l, r, sz int) {
		if node.Left == nil && node.Right == nil {
			ans = max(ans, 1)
			return node.Val, node.Val, 1
		}

		sz = 1
		valid := true
		l, r = node.Val, node.Val
		if node.Left != nil {
			l1, r1, sz1 := dfs(node.Left)
			if sz1 != -1 && node.Val > r1 {
				sz += sz1
				l = l1
			} else {
				valid = false
			}
		}

		if node.Right != nil {
			l2, r2, sz2 := dfs(node.Right)
			if sz2 != -1 && node.Val < l2 {
				sz += sz2
				r = r2
			} else {
				valid = false
			}
		}

		if valid {
			ans = max(ans, sz)
			return l, r, sz
		}

		// 不是BST，sz设为-1标记不是BST，l,r多少都可以
		return -1, -1, -1
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
