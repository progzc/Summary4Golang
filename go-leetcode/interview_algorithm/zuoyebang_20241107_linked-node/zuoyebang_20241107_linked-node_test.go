package zuoyebang_20241107_linked_node

import (
	"fmt"
	"testing"
)

// 作业帮（一面）
// 先将链表相加，再进行反转

func TestAddReverse(t *testing.T) {
	// 链表 1：1->9->9
	// 链表 2：9
	//   相加：0->0->0->1
	//   反转: 1->0->0->0
	l1 := &ListNode{Val: 1, Next: &ListNode{Val: 9, Next: &ListNode{Val: 9}}}
	l2 := &ListNode{Val: 9}

	root := Reverse(Add(l1, l2))
	for root != nil {
		if root.Next != nil {
			fmt.Printf("%d->", root.Val)
		} else {
			fmt.Printf("%d", root.Val)
		}
		root = root.Next
	}
	fmt.Println()
	// Output:
	// 1->0->0->0
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func Reverse(root *ListNode) *ListNode {
	if root == nil || root.Next == nil {
		return root
	}
	var (
		pre *ListNode
		cur = root
	)
	for cur != nil {
		next := cur.Next
		cur.Next = pre
		pre = cur
		cur = next
	}
	return pre
}

func Add(l1, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	}
	if l2 == nil {
		return l1
	}

	dummy := &ListNode{}
	head := dummy
	last := 0
	for l1 != nil || l2 != nil {
		var x, y int
		if l1 != nil {
			x = l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			y = l2.Val
			l2 = l2.Next
		}
		mod := (last + x + y) % 10
		last = (last + x + y) / 10
		head.Next = &ListNode{Val: mod}
		head = head.Next
	}
	if last > 0 {
		head.Next = &ListNode{Val: last}
		head = head.Next
	}
	return dummy.Next
}
