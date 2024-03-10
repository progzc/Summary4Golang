package leetcode_0143_reorder_list

// 0143. 重排链表
// https://leetcode.cn/problems/reorder-list/

type ListNode struct {
	Val  int
	Next *ListNode
}

// reorderList 队列
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func reorderList(head *ListNode) {
	var queue []*ListNode
	for head != nil {
		queue = append(queue, head)
		head = head.Next
	}

	for len(queue) > 2 {
		first, last := queue[0], queue[len(queue)-1]
		queue = queue[1 : len(queue)-1]
		last.Next = first.Next
		first.Next = last
		queue[len(queue)-1].Next = nil
	}
}

// reorderList_2 列表
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func reorderList_2(head *ListNode) {
	cur := head
	var stack []*ListNode
	for cur != nil {
		stack = append(stack, cur)
		cur = cur.Next
	}

	for len(stack) >= 3 {
		first := stack[0]
		next := first.Next

		before := stack[len(stack)-2]
		last := stack[len(stack)-1]

		first.Next = last
		last.Next = next
		before.Next = nil
		stack = stack[1 : len(stack)-1]
	}
	return
}
