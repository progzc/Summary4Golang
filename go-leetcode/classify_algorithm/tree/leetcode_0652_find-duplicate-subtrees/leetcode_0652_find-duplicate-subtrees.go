package leetcode_0652_find_duplicate_subtrees

import "fmt"

// 652. 寻找重复的子树
// https://leetcode.cn/problems/find-duplicate-subtrees/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// findDuplicateSubtrees 使用序列号进行唯一表示+dfs
// 时间复杂度: O(n^2)
// 空间复杂度: O(n^2)
// 思路: 一种容易想到的方法是将每一棵子树都「序列化」成一个字符串，并且保证：
//	i)相同的子树会被序列化成相同的子串；
//	ii)不同的子树会被序列化成不同的子串。
// 常见的序列化方式:
//	a.层序遍历
//	b.递归: 例如一颗以x为根节点值的子树序列化为: x(左子树的序列化结果)(右子树的序列化结果)
func findDuplicateSubtrees(root *TreeNode) []*TreeNode {
	var (
		repeat = make(map[*TreeNode]bool)
		seen   = make(map[string]*TreeNode)
		dfs    func(node *TreeNode) string
	)

	dfs = func(node *TreeNode) string {
		if node == nil {
			return ""
		}
		serial := fmt.Sprintf("%d(%s)(%s)", node.Val, dfs(node.Left), dfs(node.Right))
		if v, ok := seen[serial]; ok {
			repeat[v] = true
		} else {
			seen[serial] = node
		}
		return serial
	}

	dfs(root)
	var (
		ans = make([]*TreeNode, 0, len(repeat))
	)
	for v := range repeat {
		ans = append(ans, v)
	}
	return ans
}
