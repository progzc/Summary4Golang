package leetcode_0024_swap_nodes_in_pairs

// 两两交换链表中的节点
// https://leetcode-cn.com/problems/swap-nodes-in-pairs/

type ListNode struct {
	Val  int
	Next *ListNode
}

// swapPairs_1 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路: 使用 哑节点 （或称 哨兵节点）
func swapPairs_1(head *ListNode) *ListNode {
	sentinel := &ListNode{
		Val:  0,
		Next: head,
	}

	cur := sentinel
	for cur.Next != nil && cur.Next.Next != nil {
		node1 := cur.Next
		node2 := cur.Next.Next

		cur.Next = node2
		node1.Next = node2.Next
		node2.Next = node1

		cur = node1
	}
	return sentinel.Next
}

// swapPairs_2 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路: 思考终止条件 + 画图
func swapPairs_2(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	node1 := head
	node2 := head.Next
	node3 := swapPairs_2(head.Next.Next)

	node2.Next = node1
	node1.Next = node3

	return node2
}
