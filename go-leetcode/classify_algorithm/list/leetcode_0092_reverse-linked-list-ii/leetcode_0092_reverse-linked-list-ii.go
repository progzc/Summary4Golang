package leetcode_0092_reverse_linked_list_ii

// 0092.反转链表 II
// https://leetcode-cn.com/problems/reverse-linked-list-ii/

type ListNode struct {
	Val  int
	Next *ListNode
}

// reverseBetween 常规思路
// 时间复杂度: O(2n)
// 空间复杂度: O(1)
// 思路：
//	找准4个节点，先切割，再反转，最后再接回来。
//	使用虚拟节点可以简化 边界条件
// 缺点：需要遍历两次，一次是找到leftNode和rightNode；另一次是反转链表
func reverseBetween(head *ListNode, left int, right int) *ListNode {
	dummyNode := &ListNode{Next: head}

	// 寻找preNode
	preNode := dummyNode
	for i := 0; i < left-1; i++ {
		preNode = preNode.Next
	}

	// 寻找rightNode
	rightNode := preNode
	for i := 0; i < right-left+1; i++ {
		rightNode = rightNode.Next
	}

	// 关键：需要切断链表，截取需要反转的部分
	// 截取待反转子链表
	leftNode := preNode.Next
	nextNode := rightNode.Next
	// 切断链表
	preNode.Next = nil
	rightNode.Next = nil

	// 反转链表
	reverse(leftNode)

	// 接回切断的部分
	preNode.Next = rightNode
	leftNode.Next = nextNode

	return dummyNode.Next
}

func reverse(head *ListNode) {
	var pre *ListNode
	cur := head
	for cur != nil {
		next := cur.Next
		cur.Next = pre
		pre = cur
		cur = next
	}
}

// reverseBetween_2 头插法
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 优点：只需要遍历一次
func reverseBetween_2(head *ListNode, left int, right int) *ListNode {
	dummy := &ListNode{Next: head}

	pre := dummy
	for i := 0; i < left-1; i++ {
		pre = pre.Next
	}

	cur := pre.Next
	for i := 0; i < right-left; i++ {
		// 缓存next
		next := cur.Next
		// 穿针引线
		cur.Next = next.Next
		// 下面这一步最关键：说明是头插法
		// 若写成 next.Next = cur 则会出错
		next.Next = pre.Next
		pre.Next = next
	}
	return dummy.Next
}
