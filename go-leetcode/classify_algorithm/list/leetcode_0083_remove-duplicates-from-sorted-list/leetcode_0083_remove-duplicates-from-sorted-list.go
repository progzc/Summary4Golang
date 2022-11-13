package leetcode_0083_remove_duplicates_from_sorted_list

// 0083. 删除排序链表中的重复元素
// https://leetcode.cn/problems/remove-duplicates-from-sorted-list/

type ListNode struct {
	Val  int
	Next *ListNode
}

// deleteDuplicates 顺序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func deleteDuplicates(head *ListNode) *ListNode {
	first := head
	for first != nil && first.Next != nil {
		if first.Val == first.Next.Val {
			first.Next = first.Next.Next
		} else {
			first = first.Next
		}
	}
	return head
}
