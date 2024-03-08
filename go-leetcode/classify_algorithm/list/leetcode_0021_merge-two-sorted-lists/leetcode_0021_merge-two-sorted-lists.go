package leetcode_0021_merge_two_sorted_lists

// 0021.合并两个有序链表
// https://leetcode-cn.com/problems/merge-two-sorted-lists/

type ListNode struct {
	Val  int
	Next *ListNode
}

// mergeTwoLists_1 递归
// 时间复杂度: O(n+m)
// 空间复杂度: O(n+m)
func mergeTwoLists_1(list1 *ListNode, list2 *ListNode) *ListNode {
	if list1 == nil {
		return list2
	}
	if list2 == nil {
		return list1
	}
	if list1.Val < list2.Val {
		list1.Next = mergeTwoLists_1(list1.Next, list2)
		return list1
	} else {
		list2.Next = mergeTwoLists_1(list1, list2.Next)
		return list2
	}
}

// mergeTwoLists_2 迭代
// 时间复杂度: O(n+m)
// 空间复杂度: O(1)
func mergeTwoLists_2(list1 *ListNode, list2 *ListNode) *ListNode {
	dummy := new(ListNode)
	prev := dummy
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
	return dummy.Next
}
