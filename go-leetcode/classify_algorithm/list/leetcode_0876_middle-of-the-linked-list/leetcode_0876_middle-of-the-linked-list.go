package leetcode_0876_middle_of_the_linked_list

// 0876. 链表的中间结点
// https://leetcode.cn/problems/middle-of-the-linked-list/

type ListNode struct {
	Val  int
	Next *ListNode
}

// middleNode 快慢指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func middleNode(head *ListNode) *ListNode {
	dummy := &ListNode{Next: head}
	slow, fast := dummy, dummy

	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
	}
	if fast == nil {
		return slow
	}
	return slow.Next
}
