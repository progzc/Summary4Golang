package leetcode_0234_palindrome_linked_list

// 0234.回文链表
// https://leetcode-cn.com/problems/palindrome-linked-list/

type ListNode struct {
	Val  int
	Next *ListNode
}

// isPalindrome：反转+双指针
// 时间复杂度：O(n)
// 空间复杂度: O(1)
// 思路：
//	我们可以将链表的后半部分反转（修改链表结构），然后将前半部分和后半部分进行比较。比较完成后我们应该将链表恢复原样。
//	虽然不需要恢复也能通过测试用例，但是使用该函数的人通常不希望链表结构被更改。
// 具体步骤：
//	a.找到前半部分链表的尾节点。（具体找的话可以使用快慢指针；也可以先计数，再定位到中间。不过推荐使用前一种方法）
//	b.反转后半部分链表。（反转链表）
//	c.判断是否回文。
//	d.恢复链表。
//	e.返回结果。
func isPalindrome(head *ListNode) bool {
	if head == nil {
		return true
	}

	// a.找到前半部分链表的尾节点
	firstHalfEnd := endOfFirstHalf(head)
	// b.反转后半部分链表（注意：不需要切断与前半部分的联系）
	secondHalfStart := reverse(firstHalfEnd.Next)
	// c.判断是否回文
	p1, p2 := head, secondHalfStart
	result := true
	for p2 != nil {
		if p2.Val != p1.Val {
			result = false
			break
		}
		p1 = p1.Next
		p2 = p2.Next
	}
	// d.恢复链表
	firstHalfEnd.Next = reverse(secondHalfStart)

	// e.返回结果
	return result
}

func endOfFirstHalf(head *ListNode) *ListNode {
	dummy := &ListNode{Next: head}
	slow, fast := dummy, dummy
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
	}
	return slow
}

func reverse(head *ListNode) *ListNode {
	var pre *ListNode
	cur := head
	for cur != nil {
		next := cur.Next
		cur.Next = pre
		pre = cur
		cur = next
	}
	return pre
}

// isPalindrome_2 借助数组/栈
// 时间复杂度：O(n+n/2)
// 空间复杂度: O(n)
func isPalindrome_2(head *ListNode) bool {
	var vals []int
	for head != nil {
		vals = append(vals, head.Val)
		head = head.Next
	}

	n := len(vals)
	for i := 0; i < n/2; i++ {
		if vals[i] != vals[n-i-1] {
			return false
		}
	}
	return true
}

// isPalindrome_3 递归（具有栈后进先出的特点）
// 时间复杂度：O(n)
// 空间复杂度: O(n)
// 缺点：性能比较低，一般不推荐使用
func isPalindrome_3(head *ListNode) bool {
	front := head
	var check func(cur *ListNode) bool
	check = func(cur *ListNode) bool {
		if cur != nil {
			if !check(cur.Next) {
				return false
			}
			if cur.Val != front.Val {
				return false
			}
			front = front.Next
		}
		return true
	}
	return check(head)
}
