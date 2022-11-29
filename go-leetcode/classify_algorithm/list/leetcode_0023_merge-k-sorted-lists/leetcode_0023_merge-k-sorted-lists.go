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
// 思路：由合并两个链表==>拓展想到使用分治
func mergeKLists_2(lists []*ListNode) *ListNode {
	n := len(lists)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return lists[0]
	}

	l, r := 0, n-1
	mid := l + (r-l)>>2
	return merge(mergeKLists_2(lists[l:mid+1]), mergeKLists_2(lists[mid+1:]))
}

func merge(list1, list2 *ListNode) *ListNode {
	dummy := &ListNode{0, nil}
	cur := dummy
	for list1 != nil && list2 != nil {
		if list1.Val < list2.Val {
			cur.Next = list1
			list1 = list1.Next
		} else {
			cur.Next = list2
			list2 = list2.Next
		}
		cur = cur.Next
	}

	if list1 != nil {
		cur.Next = list1
	}
	if list2 != nil {
		cur.Next = list2
	}
	return dummy.Next
}
