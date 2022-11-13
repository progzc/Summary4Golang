package SwordToOffer_0006_cong_wei_dao_tou_da_yin_lian_biao_lcof

// 剑指 Offer 06. 从尾到头打印链表
// https://leetcode.cn/problems/cong-wei-dao-tou-da-yin-lian-biao-lcof/

type ListNode struct {
	Val  int
	Next *ListNode
}

// reversePrint 顺序遍历+交换
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func reversePrint(head *ListNode) []int {
	var ans []int
	for head != nil {
		ans = append(ans, head.Val)
		head = head.Next
	}
	n := len(ans)
	for i := 0; i < n/2; i++ {
		ans[i], ans[n-i-1] = ans[n-i-1], ans[i]
	}
	return ans
}

// reversePrint_2 顺序遍历+切片头插（效率较低）
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func reversePrint_2(head *ListNode) []int {
	var ans []int
	for head != nil {
		ans = append([]int{head.Val}, ans...)
		head = head.Next
	}
	return ans
}
