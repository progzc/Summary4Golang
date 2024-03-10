package leetcode_0082_remove_duplicates_from_sorted_list_ii

import "math"

// 0082. 删除排序链表中的重复元素 II
// https://leetcode.cn/problems/remove-duplicates-from-sorted-list-ii/

type ListNode struct {
	Val  int
	Next *ListNode
}

// deleteDuplicates 顺序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func deleteDuplicates(head *ListNode) *ListNode {
	dummy := &ListNode{math.MinInt32, head}
	cur := dummy.Next
	pre := dummy
	for cur != nil && cur.Next != nil {
		if cur.Val != cur.Next.Val {
			pre = cur
			cur = cur.Next
		} else {
			for cur.Next != nil && cur.Val == cur.Next.Val {
				cur = cur.Next
			}
			pre.Next = cur.Next
			cur = cur.Next
		}
	}

	return dummy.Next
}
