package leetcode_0002_add_two_numbers

// 0002.两数相加
// https://leetcode-cn.com/problems/add-two-numbers/

type ListNode struct {
	Val  int
	Next *ListNode
}

// addTwoNumbers 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	var tail, head *ListNode
	carry := 0
	for l1 != nil || l2 != nil {
		n1, n2 := 0, 0
		if l1 != nil {
			n1 = l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			n2 = l2.Val
			l2 = l2.Next
		}
		sum := n1 + n2 + carry
		sum, carry = sum%10, sum/10
		if head == nil {
			head = &ListNode{Val: sum}
			tail = head
		} else {
			tail.Next = &ListNode{Val: sum}
			tail = tail.Next
		}
	}
	if carry > 0 {
		tail.Next = &ListNode{Val: carry}
	}
	return head
}
