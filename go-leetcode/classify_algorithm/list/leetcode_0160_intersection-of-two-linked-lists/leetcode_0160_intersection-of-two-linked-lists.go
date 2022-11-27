package leetcode_0160_intersection_of_two_linked_lists

// 0160.相交链表
// https://leetcode-cn.com/problems/intersection-of-two-linked-lists/

type ListNode struct {
	Val  int
	Next *ListNode
}

// getIntersectionNode_1 哈希表
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func getIntersectionNode_1(headA, headB *ListNode) *ListNode {
	seen := map[*ListNode]struct{}{}
	v1, v2 := headA, headB
	for v1 != nil {
		seen[v1] = struct{}{}
		v1 = v1.Next
	}
	for v2 != nil {
		if _, ok := seen[v2]; ok {
			return v2
		}
		v2 = v2.Next
	}
	return nil
}

// getIntersectionNode_2 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func getIntersectionNode_2(headA, headB *ListNode) *ListNode {
	if headA == nil || headB == nil {
		return nil
	}
	pa, pb := headA, headB
	for pa != pb {
		if pa == nil {
			pa = headB
		} else {
			pa = pa.Next
		}

		if pb == nil {
			pb = headA
		} else {
			pb = pb.Next
		}
	}
	return pa
}

// getIntersectionNode_3 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func getIntersectionNode_3(headA, headB *ListNode) *ListNode {
	dummy1 := &ListNode{0, headA}
	dummy2 := &ListNode{0, headB}
	x, y := dummy1, dummy2
	f1, f2 := true, true
	for x != nil && y != nil {
		if x == y {
			return x
		}
		x = x.Next
		if x == nil && f1 {
			x = dummy2.Next
			f1 = !f1
		}

		y = y.Next
		if y == nil && f2 {
			y = dummy1.Next
			f2 = !f2
		}
	}
	return nil
}
