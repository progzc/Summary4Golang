package leetcode_0666_path_sum_iv

// 666. 路径总和 IV
// https://leetcode.cn/problems/path-sum-iv/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// pathSum 常规思路
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路: 很容易想到的思路,先重建树,再求路径总和
//  在树的构造过程中，我们有深度、位置和权值这些信息，我们可以根据条件 pos-1 < 2*(depth-2) 来判断结点在右边还是左边。
func pathSum(nums []int) int {
	root := &TreeNode{Val: nums[0] % 10}
	for i, num := range nums {
		if i == 0 {
			continue
		}
		depth, pos, val := num/100, num/10%10, num%10
		pos--
		cur := root
		for d := depth - 2; d >= 0; d-- {
			if pos < 1<<d {
				if cur.Left == nil {
					cur.Left = &TreeNode{Val: val}
				}
				cur = cur.Left
			} else {
				if cur.Right == nil {
					cur.Right = &TreeNode{Val: val}
				}
				cur = cur.Right
			}
			pos %= 1 << d
		}
	}

	var (
		dfs func(root *TreeNode, sum int)
		ans int
	)

	dfs = func(root *TreeNode, sum int) {
		if root == nil {
			return
		}
		sum += root.Val
		if root.Left == nil && root.Right == nil {
			ans += sum
		} else {
			dfs(root.Left, sum)
			dfs(root.Right, sum)
		}
	}

	dfs(root, 0)
	return ans
}

// pathSum_2 map+直接遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路: 直接进行遍历，使用depth+pos来唯一标识一个节点
//  我们根据等式 root = num / 10 = 10 * depth + pos 作为根节点的唯一标识符。则:
//  左子结点的标识符是 left = 10 * (depth + 1) + 2 * pos - 1，
//  右子节点则是 right = left + 1。
func pathSum_2(nums []int) int {
	var (
		m   = make(map[int]int, len(nums))
		ans int
		dfs func(root, sum int)
	)

	for _, num := range nums {
		m[num/10] = num % 10
	}

	dfs = func(root, sum int) {
		if _, ok := m[root]; !ok {
			return
		}
		sum += m[root]
		depth, pos := root/10, root%10
		left := (depth+1)*10 + 2*pos - 1
		right := left + 1

		_, ok1 := m[left]
		_, ok2 := m[right]
		if !ok1 && !ok2 {
			ans += sum
		} else {
			dfs(left, sum)
			dfs(right, sum)
		}
	}

	dfs(nums[0]/10, 0)
	return ans
}
