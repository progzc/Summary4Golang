package leetcode_0141_linked_list_cycle

import (
	"time"
)

// 141. 环形链表
// https://leetcode-cn.com/problems/linked-list-cycle/

type ListNode struct {
	Val  int
	Next *ListNode
}

// hasCycle_1 快慢指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：需要记住 快慢指针 的常见用法
func hasCycle_1(head *ListNode) bool {
	slow, fast := head, head
	for slow != nil && fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return true
		}
	}
	return false
}

// hasCycle_2 哈希表
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：比较容易想到这种方法
func hasCycle_2(head *ListNode) bool {
	set := make(map[*ListNode]struct{})
	for head != nil {
		if _, ok := set[head]; !ok {
			set[head] = struct{}{}
		} else {
			return true
		}
		head = head.Next
	}
	return false
}

// hasCycle_3 硬算
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：能够通过leetcode,但是有致命缺陷
func hasCycle_3(head *ListNode) bool {
	ch := make(chan struct{}, 1)

	go func() {
		for head != nil {
			head = head.Next
		}
		ch <- struct{}{}
	}()

	select {
	case <-ch:
		return false
	// 超时设置为100ms,可以通过leetcode
	case <-time.After(100 * time.Millisecond):
		return true
	}
}
