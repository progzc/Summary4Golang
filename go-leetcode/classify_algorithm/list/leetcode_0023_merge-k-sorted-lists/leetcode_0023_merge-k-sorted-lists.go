package leetcode_0023_merge_k_sorted_lists

import "math"

// 0023.合并K个升序链表
// https://leetcode-cn.com/problems/merge-k-sorted-lists/

type ListNode struct {
	Val  int
	Next *ListNode
}

// mergeKLists 递归
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
// 思路：自己想出来的一种很自然的解法
func mergeKLists(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}
	var (
		result *ListNode
		min    = math.MaxInt32
		idx    = -1
	)
	// 找到最小的那个节点
	for i := 0; i < len(lists); i++ {
		if lists[i] == nil {
			continue
		}
		if lists[i].Val < min {
			min = lists[i].Val
			idx = i
		}
	}
	// 更新
	if idx != -1 {
		result = lists[idx]
		lists[idx] = lists[idx].Next
	}
	// 递归
	if result != nil {
		result.Next = mergeKLists(lists)
	}
	return result
}

// mergeKLists_2 二分+分治
// 时间复杂度: O(k*n*log(n))
// 空间复杂度: O(log(k))
// 思路：由合并两个链表==拓展想到使用分治
func mergeKLists_2(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}
	return merge(lists, 0, len(lists)-1)
}

func merge(lists []*ListNode, left, right int) *ListNode {
	if left == right {
		return lists[left]
	}
	if left > right {
		return nil
	}
	mid := (left + right) >> 1
	return mergeTwoLists(merge(lists, left, mid), merge(lists, mid+1, right))

}

func mergeTwoLists(list1, list2 *ListNode) *ListNode {
	prev := new(ListNode)
	sentry := prev
	for list1 != nil && list2 != nil {
		if list1.Val < list2.Val {
			prev.Next = list1
			list1 = list1.Next
		} else {
			prev.Next = list2
			list2 = list2.Next
		}
		prev = prev.Next
	}
	if list1 == nil {
		prev.Next = list2
	} else {
		prev.Next = list1
	}
	return sentry.Next
}
