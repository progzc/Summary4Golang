package leetcode_0025_reverse_nodes_in_k_group

// 0025. K 个一组翻转链表
// https://leetcode.cn/problems/reverse-nodes-in-k-group/

type ListNode struct {
	Val  int
	Next *ListNode
}

// reverseKGroup
// 时间复杂度: O(nk)
// 空间复杂度: O(1)
func reverseKGroup(head *ListNode, k int) *ListNode {
	dummy := &ListNode{0, head}
	pre, end := dummy, dummy
	for end.Next != nil {
		for i := 0; i < k && end != nil; i++ {
			end = end.Next
		}
		if end == nil {
			break
		}
		start, next := pre.Next, end.Next
		end.Next = nil
		pre.Next = reverse(start)
		start.Next = next
		pre = start
		end = pre
	}
	return dummy.Next
}

func reverse(head *ListNode) *ListNode {
	var (
		pre *ListNode
		cur = head
	)
	for cur != nil {
		next := cur.Next
		cur.Next = pre
		pre = cur
		cur = next
	}
	return pre
}
