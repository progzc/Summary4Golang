package SwordToOffer_0026_shu_de_zi_jie_gou_lcof

// 剑指 Offer 26. 树的子结构
// https://leetcode.cn/problems/shu-de-zi-jie-gou-lcof/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// isSubStructure_2 递归(对 isSubStructure 代码进行精简)
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func isSubStructure_2(A *TreeNode, B *TreeNode) bool {
	return (A != nil && B != nil) && (isSubStructure_2(A.Left, B) || isSubStructure_2(A.Right, B) || isSimilar_2(A, B))
}
func isSimilar_2(A, B *TreeNode) bool {
	if B == nil {
		return true
	}
	if A == nil {
		return false
	}
	if A.Val != B.Val {
		return false
	}
	return isSimilar_2(A.Left, B.Left) && isSimilar_2(A.Right, B.Right)
}

// ---------------------------------------------------------------------------------------------------------------------

// isSubStructure 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func isSubStructure(A *TreeNode, B *TreeNode) bool {
	if A == nil && B == nil {
		return false
	}

	if A == nil || B == nil {
		return false
	}
	if A.Val != B.Val {
		return isSubStructure(A.Left, B) || isSubStructure(A.Right, B)
	}
	return isSubStructure(A.Left, B) || isSubStructure(A.Right, B) || isSimilar(A, B)
}

func isSimilar(A, B *TreeNode) bool {
	if A == nil && B == nil {
		return true
	}
	if A == nil {
		return false
	}
	if B == nil {
		return true
	}
	if A.Val != B.Val {
		return false
	}
	return isSimilar(A.Left, B.Left) && isSimilar(A.Right, B.Right)
}
