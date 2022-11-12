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
		if len(queue) > 0 {
			queue[len(queue)-1].Next = nil
		}
	}
}
