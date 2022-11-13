package leetcode_0328_odd_even_linked_list

// 0328.奇偶链表
// https://leetcode-cn.com/problems/odd-even-linked-list/

type ListNode struct {
	Val  int
	Next *ListNode
}

// oddEvenList
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func oddEvenList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil || head.Next.Next == nil {
		return head
	}
	// odd代表奇数的链表
	// even代表偶数的链表
	odd, even := head, head.Next
	curOdd, curEven := odd, even

	head = head.Next.Next
	i := 0
	for head != nil {
		i++
		if i%2 == 1 {
			curOdd.Next = head
			curOdd = curOdd.Next
		} else {
			curEven.Next = head
			curEven = curEven.Next
		}
		head = head.Next
	}

	curOdd.Next = even
	curEven.Next = nil
	return odd
}
