package order_leetcode

// leetcode_0142_detectCycle_method2 快慢指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func leetcode_0142_detectCycle_method2(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast != nil {
		slow = slow.Next
		if fast.Next == nil {
			return nil
		}
		fast = fast.Next.Next
		if fast == slow {
			p := head
			for p != slow {
				p = p.Next
				slow = slow.Next
			}
			return p
		}
	}
	return nil
}

// leetcode_0142_detectCycle_method1 哈希表
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func leetcode_0142_detectCycle_method1(head *ListNode) *ListNode {
	listMap := map[*ListNode]struct{}{}
	for head != nil {
		if _, ok := listMap[head]; ok {
			return head
		}
		listMap[head] = struct{}{}
		head = head.Next
	}
	return nil
}
