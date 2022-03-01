package leetcode_0142_linked_list_cycle_ii

// 0142.环形链表 II
// https://leetcode-cn.com/problems/linked-list-cycle-ii/

type ListNode struct {
	Val  int
	Next *ListNode
}

// detectCycle_1 哈希表
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：比较容易想到这种方法
func detectCycle_1(head *ListNode) *ListNode {
	set := map[*ListNode]struct{}{}
	for head != nil {
		if _, ok := set[head]; ok {
			return head
		}
		set[head] = struct{}{}
		head = head.Next
	}
	return nil
}

// detectCycle_2 快慢指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：凡是链表的问题务必要想到 哨兵节点 / 快慢指针 / 递归 这三种常规解法上面去
// 快慢指针的思路：其实是一道 龟兔赛跑 的路径题，主要是画图分析出 a = c + (n-1)(b+c)
func detectCycle_2(head *ListNode) *ListNode {
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
