package leetcode_0237_delete_node_in_a_linked_list

// 0237.删除链表中的节点
// https://leetcode-cn.com/problems/delete-node-in-a-linked-list/

type ListNode struct {
	Val  int
	Next *ListNode
}

// deleteNode 重点在理解题意
// 时间复杂度: O(1)
// 空间复杂度: O(1)
// 思路：
//	给定链表 4-->5-->1-->9，要被删除的节点是5，即链表中的第2个节点。 可以通过如下两步操作实现删除节点的操作：
//		a.将第2个节点的值修改为第3个节点的值，即将节点5的值修改为1，此时链表如下：4-->1-->1-->9
//		b.删除第3个节点，此时链表如下：4-->1-->9
func deleteNode(node *ListNode) {
	node.Val = node.Next.Val
	node.Next = node.Next.Next
}
