package leetcode_0148_sort_list

// 0148.排序链表
// https://leetcode-cn.com/problems/sort-list/

type ListNode struct {
	Val  int
	Next *ListNode
}

// sortList 归并排序+递归（即自顶向下归并排序）
// 空间复杂度: O(n*log(n))
// 时间复杂度: O(log(n))
// 思路：寻找链表的中点，可以考虑使用快慢指针。
//		以中点为分界，将链表拆分成两个子链表。寻找链表的中点可以使用快慢指针的做法，快指针每次移动 22 步，慢指针每次移动 11 步，
//		当快指针到达链表末尾时，慢指针指向的链表节点即为链表的中点。
func sortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	// 找中点
	slow, fast := head, head.Next
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	// 此时中点就是slow
	head2 := slow.Next
	slow.Next = nil
	return merge(sortList(head), sortList(head2))
}

// sortList_2 归并排序+迭代（即自底向上归并排序）
// 空间复杂度: O(n*log(n))
// 时间复杂度: O(1)
// 思路：寻找链表的中点，可以考虑使用快慢指针。
//		以中点为分界，将链表拆分成两个子链表。寻找链表的中点可以使用快慢指针的做法，快指针每次移动 22 步，慢指针每次移动 11 步，
//		当快指针到达链表末尾时，慢指针指向的链表节点即为链表的中点。
func sortList_2(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	length := 0
	for node := head; node != nil; node = node.Next {
		length++
	}

	dummy := &ListNode{Next: head}

	// interval 标识每轮合并链表长度
	for interval := 1; interval < length; interval <<= 1 {
		// 每次尝试走interval-1步就断开一个链表，a、b分别记录断开的两个子链表头节点
		var a, b *ListNode
		var cur, pre *ListNode
		cur, pre = dummy.Next, dummy
		dummy.Next = nil

		for cur != nil {
			var (
				nextCur *ListNode
				i       int
			)
			i = 0   // cur尝试往后走interval-1步
			a = cur // 记录首个链表头
			for i < interval-1 && cur.Next != nil {
				cur = cur.Next
				i++
			}

			b = cur.Next   // 记录第二个链表头
			cur.Next = nil // 断开首个链表

			i = 0
			cur = b
			if cur != nil { // 注意第二个链表可能不存在，如果为nil的话，直接拿a和nil合并
				for i < interval-1 && cur.Next != nil {
					cur = cur.Next
					i++
				}
				nextCur = cur.Next // 记录下一轮的cur
				cur.Next = nil     // 断开第二个链表
			}
			// 合并两个断开的链表
			node := merge(a, b)
			pre.Next = node

			// 更新新一轮尾部
			for pre.Next != nil {
				pre = pre.Next
			}
			cur = nextCur
		}
	}
	return dummy.Next
}

// merge 合并两个有序链表
// 参考：
//	0021.合并两个有序链表
//	https://leetcode-cn.com/problems/merge-two-sorted-lists/
func merge(list1, list2 *ListNode) *ListNode {
	prev := new(ListNode)
	sentry := prev
	for list1 != nil && list2 != nil {
		if list1.Val < list2.Val {
			prev.Next = list1
			list1 = list1.Next
		} else {
			prev.Next = list2
			list2 = list2.Next
		}
		prev = prev.Next
	}
	if list1 == nil {
		prev.Next = list2
	} else {
		prev.Next = list1
	}
	return sentry.Next
}
