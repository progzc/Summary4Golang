package leetcode_0099_recover_binary_search_tree

// 0099.恢复二叉搜索树
// https://leetcode-cn.com/problems/recover-binary-search-tree/
// 关键信息: 恰好有两个节点被错误交换

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// recoverTree_1 中序遍历 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	a.根据中序遍历获得二叉搜索树的节点值的切片
//	b.找到乱序的两个数
//	  这里注意要分两种情况: 交换的是相邻的两个数；交换的不是相邻的两个数
func recoverTree_1(root *TreeNode) {
	var (
		dfs  func(root *TreeNode)
		find func(nums []int) (int, int)
		swap func(root *TreeNode, x, y, cnt int)
		nums []int
	)

	dfs = func(root *TreeNode) {
		if root == nil {
			return
		}
		dfs(root.Left)
		nums = append(nums, root.Val)
		dfs(root.Right)
	}

	find = func(nums []int) (int, int) {
		idx1, idx2 := -1, -1
		for i := 0; i < len(nums)-1; i++ {
			if nums[i+1] < nums[i] {
				idx2 = i + 1
				if idx1 == -1 {
					idx1 = i
				} else {
					// 提前退出循环
					break
				}
			}
		}
		return nums[idx1], nums[idx2]
	}

	swap = func(root *TreeNode, x, y, cnt int) {
		if root == nil || cnt == 0 {
			return
		}
		if root.Val == x || root.Val == y {
			if root.Val == x {
				root.Val = y
			} else {
				root.Val = x
			}
			cnt--
		}
		// 这个顺序不能挪到上面的if语句上,不然会出错
		swap(root.Left, x, y, cnt)
		swap(root.Right, x, y, cnt)
	}

	dfs(root)
	x, y := find(nums)
	swap(root, x, y, 2)
}

// recoverTree_2 中序遍历 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	recoverTree_1中进行了多余的两次遍历和查找，可以针对这项进行优化
func recoverTree_2(root *TreeNode) {
	var (
		stack     []*TreeNode
		x, y, pre *TreeNode
	)

	for root != nil || len(stack) > 0 {
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}
		root = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if pre != nil && root.Val < pre.Val {
			y = root
			if x == nil {
				x = pre
			} else {
				break
			}
		}
		pre = root
		root = root.Right
	}
	x.Val, y.Val = y.Val, x.Val
}
