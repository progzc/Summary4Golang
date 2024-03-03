package leetcode_0147_insertion_sort_list

// 0147. 对链表进行插入排序
// https://leetcode.cn/problems/insertion-sort-list/

type ListNode struct {
	Val  int
	Next *ListNode
}

// insertionSortList 链表+插入排序
// 时间复杂度：O(n^2)
// 空间复杂度：O(1)
// 思路：递归+哨兵节点
func insertionSortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	head.Next = insertionSortList(head.Next)

	dummy := &ListNode{Next: head}
	pre := dummy
	cur := dummy.Next
	for cur != nil && cur.Next != nil && cur.Val > cur.Next.Val {
		pre.Next = cur.Next
		tmp := cur.Next
		cur.Next = cur.Next.Next
		tmp.Next = cur
		pre = tmp
	}
	return dummy.Next
}
