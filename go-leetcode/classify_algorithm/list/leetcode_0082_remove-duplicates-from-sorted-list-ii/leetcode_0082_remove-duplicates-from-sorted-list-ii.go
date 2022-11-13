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
// 空间复杂度: O(n)
func deleteDuplicates(head *ListNode) *ListNode {
	dummy := &ListNode{math.MinInt32, head}
	first := dummy.Next
	pre := dummy
	for first != nil && first.Next != nil {
		if first.Val != first.Next.Val {
			pre = first
			first = first.Next
		} else {
			for first.Next != nil && first.Val == first.Next.Val {
				first = first.Next
			}
			pre.Next = first.Next
			first = first.Next
		}
	}

	return dummy.Next
}
