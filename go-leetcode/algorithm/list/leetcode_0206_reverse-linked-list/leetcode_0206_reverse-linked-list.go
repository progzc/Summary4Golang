package leetcode_0206_reverse_linked_list

// 反转链表
// https://leetcode-cn.com/problems/reverse-linked-list/

type ListNode struct {
	Val  int
	Next *ListNode
}

// reverseList_1 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func reverseList_1(head *ListNode) *ListNode {
	var pre *ListNode
	cur := head
	for cur != nil {
		next := cur.Next
		cur.Next = pre
		pre = cur
		cur = next
	}
	return pre
}

// reverseList_2 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func reverseList_2(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	next := reverseList_2(head.Next)
	head.Next.Next = head
	head.Next = nil
	return next
}
