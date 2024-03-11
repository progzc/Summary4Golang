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
	for cur != nil && cur.Next != nil { // 少于两个节点直接结束
		if cur.Val != cur.Next.Val { // 两个节点值不同，直接向下遍历
			pre = cur
			cur = cur.Next
		} else { // 两个节点值相同，找到下一个不同的节点
			for cur.Next != nil && cur.Val == cur.Next.Val {
				cur = cur.Next
			}
			pre.Next = cur.Next
			cur = cur.Next
		}
	}

	return dummy.Next
}

// deleteDuplicates_2 顺序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func deleteDuplicates_2(head *ListNode) *ListNode {
	dummy := &ListNode{Next: head}
	pre := dummy
	cur := dummy.Next
	for cur != nil && cur.Next != nil {
		first := cur
		for cur.Next != nil && cur.Val == cur.Next.Val {
			cur = cur.Next
		}
		if cur.Next == nil {
			pre.Next = nil
		} else {
			if first == cur {
				pre = cur
			} else {
				pre.Next = cur.Next
			}
		}
		cur = cur.Next
	}
	return dummy.Next
}
