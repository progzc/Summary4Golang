package order_leetcode

type ListNode struct {
	Val  int
	Next *ListNode
}

func initListNode(nums []int) *ListNode {
	var head, tail *ListNode
	for _, num := range nums {
		if head == nil {
			tail = &ListNode{Val: num}
			head = tail
		} else {
			tail.Next = &ListNode{Val: num}
			tail = tail.Next
		}
	}
	return head
}

func listNode2Slice(head *ListNode) []int {
	var nums []int
	for head != nil {
		nums = append(nums, head.Val)
		head = head.Next
	}
	return nums
}
