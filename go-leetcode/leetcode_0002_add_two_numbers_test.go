package go_leetcode

import (
	"reflect"
	"testing"
)

func Test_leetcode_0002_addTwoNumbers(t *testing.T) {
	type params struct {
		l1 []int
		l2 []int
	}
	tests := []struct {
		p    params
		want []int
	}{
		{
			p: params{
				l1: []int{2, 4, 3},
				l2: []int{5, 6, 4},
			},
			want: []int{7, 0, 8},
		},
		{
			p: params{
				l1: []int{0},
				l2: []int{0},
			},
			want: []int{0},
		},
		{
			p: params{
				l1: []int{9, 9, 9, 9, 9, 9, 9},
				l2: []int{9, 9, 9, 9},
			},
			want: []int{8, 9, 9, 9, 0, 0, 0, 1},
		},
	}
	for _, test := range tests {
		l1 := initListNode(test.p.l1)
		l2 := initListNode(test.p.l2)
		fact := listNode2Slice(leetcode_0002_addTwoNumbers(l1, l2))
		if !reflect.DeepEqual(fact, test.want) {
			t.Errorf("params=%v,want=%v,fact=%v",
				test.p, test.want, fact)
		}
	}
}

func leetcode_0002_addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	var tail, head *ListNode
	carry := 0
	for l1 != nil || l2 != nil {
		n1, n2 := 0, 0
		if l1 != nil {
			n1 = l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			n2 = l2.Val
			l2 = l2.Next
		}
		sum := n1 + n2 + carry
		sum, carry = sum%10, sum/10
		if head == nil {
			head = &ListNode{Val: sum}
			tail = head
		} else {
			tail.Next = &ListNode{Val: sum}
			tail = tail.Next
		}
	}
	if carry > 0 {
		tail.Next = &ListNode{Val: carry}
	}
	return head
}
