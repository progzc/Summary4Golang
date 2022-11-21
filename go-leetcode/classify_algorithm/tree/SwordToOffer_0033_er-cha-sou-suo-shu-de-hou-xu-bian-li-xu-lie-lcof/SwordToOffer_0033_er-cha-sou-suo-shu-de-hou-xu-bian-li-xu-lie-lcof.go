package SwordToOffer_0033_er_cha_sou_suo_shu_de_hou_xu_bian_li_xu_lie_lcof

// 剑指 Offer 33. 二叉搜索树的后序遍历序列
// https://leetcode.cn/problems/er-cha-sou-suo-shu-de-hou-xu-bian-li-xu-lie-lcof/

// verifyPostorder 递归（常规思路）
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func verifyPostorder(postorder []int) bool {
	n := len(postorder)
	if n == 0 {
		return true
	}

	num := postorder[len(postorder)-1]
	idx := -1
	for i := 0; i < n-1; i++ {
		if postorder[i] > num {
			idx = i
			break
		}
	}

	if idx > -1 {
		for i := idx + 1; i < n-1; i++ {
			if postorder[i] <= num {
				return false
			}
		}
		return verifyPostorder(postorder[:idx]) && verifyPostorder(postorder[idx:n-1])
	}
	return verifyPostorder(postorder[:n-1])
}
