package leetcode_0023_merge_k_sorted_lists

import (
	"container/heap"
	"math"
)

// 0023.合并K个升序链表
// https://leetcode-cn.com/problems/merge-k-sorted-lists/

type ListNode struct {
	Val  int
	Next *ListNode
}

type IHeap []*ListNode

func (h *IHeap) Len() int { return len(*h) }

func (h *IHeap) Swap(i, j int) { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IHeap) Less(i, j int) bool { return (*h)[i].Val < (*h)[j].Val }

func (h *IHeap) Push(x interface{}) {
	*h = append(*h, x.(*ListNode))
}

func (h *IHeap) Pop() interface{} {
	old := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return old
}

// mergeKLists_3 优先队列
// 时间复杂度: O(k*n*log(n))
// 空间复杂度: O(log(k))
func mergeKLists_3(lists []*ListNode) *ListNode {
	dummy := &ListNode{}
	n := len(lists)
	if n == 0 {
		return dummy.Next
	}

	h := IHeap{}
	heap.Init(&h)
	for i := 0; i < n; i++ {
		if lists[i] != nil {
			heap.Push(&h, lists[i])
		}
	}

	head := dummy
	for h.Len() > 0 {
		x := heap.Pop(&h).(*ListNode)
		head.Next = x
		head = head.Next
		if x.Next != nil {
			heap.Push(&h, x.Next)
		}
	}
	return dummy.Next
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
	mid := (n - 1) / 2
	return merge(mergeKLists_2(lists[:mid+1]), mergeKLists_2(lists[mid+1:]))
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
