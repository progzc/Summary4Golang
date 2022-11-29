package baidu_20221129_linked_list_cycle

type ListNode struct {
	Val  int
	Next *ListNode
}

func hasCycle(head *ListNode) bool {
	m := make(map[*ListNode]bool)
	cur := head
	for cur != nil {
		if !m[cur] {
			m[cur] = true
		} else {
			return true
		}
		cur = cur.Next
	}
	return false
}

func hasCycle_2(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}
	dummy := &ListNode{0, head}
	slow, fast := dummy, dummy.Next
	for fast != nil && fast.Next != nil {
		if slow == fast {
			return true
		}
		slow = slow.Next
		fast = fast.Next.Next
	}
	return false
}
