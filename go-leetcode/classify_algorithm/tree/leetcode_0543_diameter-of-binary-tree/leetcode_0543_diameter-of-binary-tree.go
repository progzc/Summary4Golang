package leetcode_0543_diameter_of_binary_tree

// 0543.二叉树的直径
// https://leetcode-cn.com/problems/diameter-of-binary-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// diameterOfBinaryTree 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：
//	a.假设某节点的左、右子树的最大深度分别为L、R，则以该节点为起点的路径经过节点数的最大值即为L+R+1。
//	  从而以该节点为顶点的二叉树的直径即为(L+R+1)-1=L+R
//	b.枚举所有节点为顶点的二叉树的直径，记录最大的直径即为该二叉树的直径。
// 该题是 0104.二叉树的最大深度 的衍生题
func diameterOfBinaryTree(root *TreeNode) int {
	var (
		ans   int
		depth func(node *TreeNode) int
	)

	// depth函数即为该二叉树的最大深度
	depth = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		L := depth(node.Left)
		R := depth(node.Right)
		ans = max(ans, L+R+1)
		return max(L, R) + 1
	}
	depth(root)
	return ans - 1
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
