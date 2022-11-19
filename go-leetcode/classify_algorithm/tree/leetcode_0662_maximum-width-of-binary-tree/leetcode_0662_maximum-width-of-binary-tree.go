package leetcode_0662_maximum_width_of_binary_tree

// 0662. 二叉树最大宽度
// https://leetcode.cn/problems/maximum-width-of-binary-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type Pair struct {
	Node  *TreeNode
	Index int
}

// widthOfBinaryTree 层序遍历+满二叉树节点编号的规律
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	此题求二叉树所有层的最大宽度，比较直观的方法是求出每一层的宽度，然后求出最大值。
//	a.求每一层的宽度时，因为两端点间的 null 节点也需要计入宽度，因此可以对节点进行编号。
//	  一个编号为 index 的左子节点的编号记为 2×index，右子节点的编号记为 2×index+1，计算每层宽度时，
//	  用每层节点的最大编号减去最小编号再加 1 即为宽度。根节点编号为 1.
//	b.遍历节点时，可以用广度优先搜索来遍历每一层的节点，并求出最大值。
func widthOfBinaryTree(root *TreeNode) int {
	if root == nil {
		return 0
	}

	var (
		ans   int
		stack []Pair
	)
	stack = append(stack, Pair{root, 1})
	for len(stack) > 0 {
		n := len(stack)
		ans = max(ans, stack[n-1].Index-stack[0].Index+1)
		for i := 0; i < n; i++ {
			node := stack[0]
			stack = stack[1:]
			if node.Node.Left != nil {
				stack = append(stack, Pair{node.Node.Left, node.Index * 2})
			}
			if node.Node.Right != nil {
				stack = append(stack, Pair{node.Node.Right, node.Index*2 + 1})
			}
		}
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
