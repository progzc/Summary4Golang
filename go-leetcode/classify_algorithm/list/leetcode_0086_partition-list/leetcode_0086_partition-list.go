package leetcode_0086_partition_list

// 0086.分隔链表
// https://leetcode-cn.com/problems/partition-list/

type ListNode struct {
	Val  int
	Next *ListNode
}

// partition 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：使用两个虚拟节点
func partition(head *ListNode, x int) *ListNode {
	dummy1, dummy2 := &ListNode{}, &ListNode{}
	copy1, copy2 := dummy1, dummy2

	for node := head; node != nil; node = node.Next {
		if node.Val < x {
			dummy1.Next = node
			dummy1 = dummy1.Next
		} else {
			dummy2.Next = node
			dummy2 = dummy2.Next
		}
	}
	// 注意事项：dummy2.Next = nil 这一行不能掉
	// 考虑两种情况：
	//	情况1：原始链表的最后一个节点的值大于等于x，可以去掉 dummy2.Next = nil
	//	情况2: 原始链表的最后一个节点的值小于x，不能去掉 dummy2.Next = nil，否则就会出现一个环。
	// 		  例如：输入 head = [1,4,3,2,5,2], x = 3
	// 		  若漏掉dummy2.Next = nil 这一行，则会出现环：1—>2->2->4->3-5->2->4->3->5->...loop
	dummy2.Next = nil
	dummy1.Next = copy2.Next
	return copy1.Next
}
