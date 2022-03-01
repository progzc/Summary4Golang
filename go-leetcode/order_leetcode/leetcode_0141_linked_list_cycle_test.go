package order_leetcode

// leetcode_0141_hasCycle_method2 快慢指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func leetcode_0141_hasCycle_method2(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}
	slow, fast := head, head.Next
	for fast != slow {
		if fast == nil || fast.Next == nil {
			return false
		}
		slow = slow.Next
		fast = fast.Next.Next
	}
	return true
}

// leetcode_0141_hasCycle_method1 哈希表
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func leetcode_0141_hasCycle_method1(head *ListNode) bool {
	listMap := map[*ListNode]struct{}{}
	for head != nil {
		if _, ok := listMap[head]; ok {
			return true
		}
		listMap[head] = struct{}{}
		head = head.Next
	}
	return false
}
