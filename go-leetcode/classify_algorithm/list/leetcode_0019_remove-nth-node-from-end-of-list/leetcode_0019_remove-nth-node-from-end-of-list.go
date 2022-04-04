package leetcode_0019_remove_nth_node_from_end_of_list

// 0019.删除链表的倒数第 N 个结点
// https://leetcode-cn.com/problems/remove-nth-node-from-end-of-list/

type ListNode struct {
	Val  int
	Next *ListNode
}

// removeNthFromEnd 快慢指针
// 时间复杂度：O(n)
// 空间复杂度：O(1)
// 思路：使用一趟扫描
func removeNthFromEnd(head *ListNode, n int) *ListNode {
	dummy := &ListNode{0, head}
	slow, fast := dummy, dummy
	for n > 0 {
		n--
		fast = fast.Next
	}
	for fast.Next != nil {
		fast = fast.Next
		slow = slow.Next
	}
	slow.Next = slow.Next.Next
	return dummy.Next
}

// removeNthFromEnd_2 栈
// 时间复杂度：O(n)
// 空间复杂度：O(n)
// 思路：根据"倒数"想到要用栈
func removeNthFromEnd_2(head *ListNode, n int) *ListNode {
	var stack []*ListNode
	dummy := &ListNode{0, head}
	node := dummy // 这样避免输入为：head = [1], n = 1时，会出现索引越界
	for node != nil {
		stack = append(stack, node)
		node = node.Next
	}
	prev := stack[len(stack)-1-n]
	prev.Next = prev.Next.Next
	return dummy.Next
}
