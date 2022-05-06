package leetcode_0725_split_linked_list_in_parts

// 0725.分隔链表
// https://leetcode-cn.com/problems/split-linked-list-in-parts/

type ListNode struct {
	Val  int
	Next *ListNode
}

// splitListToParts 一般思路
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 缺点：写得过于复杂了
func splitListToParts(head *ListNode, k int) []*ListNode {
	var (
		count = 0
		ans   []*ListNode
	)

	// 统计数量
	for node := head; node != nil; count++ {
		node = node.Next
	}

	// 如果链表长度小于k
	if count <= k {
		for node := head; node != nil; {
			ans = append(ans, node)
			next := node.Next
			node.Next = nil
			node = next
		}
		for i := 0; i < k-count; i++ {
			ans = append(ans, nil)
		}
		return ans
	}

	// 否则，链表长度大于k
	m, n := count/k, count%k

	first := head
	for i := 0; i < k; i++ {
		dummy := &ListNode{Next: first}
		node := dummy
		for j := 0; j < m; j++ {
			node = node.Next
		}
		ans = append(ans, dummy.Next)

		if n > 0 {
			n--
			node = node.Next
		}
		nextFirst := node.Next
		node.Next = nil
		first = nextFirst
	}
	return ans
}

// splitListToParts_2 一般思路（优化步骤）
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func splitListToParts_2(head *ListNode, k int) []*ListNode {
	// TODO
}
