package leetcode_0203_remove_linked_list_elements

// 0203. 移除链表元素
// https://leetcode.cn/problems/remove-linked-list-elements/

type ListNode struct {
	Val  int
	Next *ListNode
}

// removeElements 顺序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func removeElements(head *ListNode, val int) *ListNode {
	dummy := &ListNode{0, head}
	pre, cur := dummy, dummy.Next
	for cur != nil {
		if cur.Val == val {
			pre.Next = cur.Next
		} else {
			pre = cur
		}
		cur = cur.Next
	}
	return dummy.Next
}
