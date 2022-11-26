package leetcode_0061_rotate_list

// 61. 旋转链表
// https://leetcode.cn/problems/rotate-list/

type ListNode struct {
	Val  int
	Next *ListNode
}

// rotateRight 快慢指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func rotateRight(head *ListNode, k int) *ListNode {
	if k == 0 || head == nil || head.Next == nil {
		return head
	}

	// 确定是否需要旋转
	dummy := &ListNode{0, head}
	cur := dummy
	count := 0
	for cur.Next != nil {
		count++
		cur = cur.Next
	}
	mod := k % count
	if mod == 0 {
		return dummy.Next
	}

	// 快慢指针找到旋转点
	slow, fast := dummy, dummy
	for mod > 0 {
		fast = fast.Next
		mod--
	}
	for fast.Next != nil {
		fast = fast.Next
		slow = slow.Next
	}
	// 进行旋转
	first := dummy.Next
	dummy.Next = slow.Next
	slow.Next = nil
	fast.Next = first
	return dummy.Next
}
