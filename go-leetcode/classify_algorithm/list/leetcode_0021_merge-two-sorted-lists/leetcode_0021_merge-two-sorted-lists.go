package leetcode_0021_merge_two_sorted_lists

// 0021.合并两个有序链表
// https://leetcode-cn.com/problems/merge-two-sorted-lists/

type ListNode struct {
	Val  int
	Next *ListNode
}

// mergeTwoLists_1 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func mergeTwoLists_1(list1 *ListNode, list2 *ListNode) *ListNode {
	if list1 == nil {
		return list2
	} else if list2 == nil {
		return list1
	} else if list1.Val > list2.Val {
		list1.Next = mergeTwoLists_1(list1.Next, list2)
		return list1
	} else {
		list2.Next = mergeTwoLists_1(list1, list2.Next)
		return list2
	}
}
