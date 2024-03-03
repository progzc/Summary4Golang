package leetcode_0109_convert_sorted_list_to_binary_search_tree

// 0109. 有序链表转换二叉搜索树
// https://leetcode.cn/problems/convert-sorted-list-to-binary-search-tree/description/

type ListNode struct {
	Val  int
	Next *ListNode
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// sortedListToBST 快慢指针+递归
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(log(n))
// 思路：找根节点。可以找出链表元素的中位数作为根节点的值。
// 这里对于中位数的定义为：如果链表中的元素个数为奇数，那么唯一的中间值为中位数；
// 如果元素个数为偶数，那么唯二的中间值都可以作为中位数，而不是常规定义中二者的平均值。
// 找中位数可以采用快慢指针的方法或者二分法！！！
func sortedListToBST(head *ListNode) *TreeNode {
	if head == nil {
		return nil
	}
	if head.Next == nil {
		return &TreeNode{Val: head.Val}
	}
	// 快慢指针找中位值
	dummy := &ListNode{Next: head}
	slow, fast := dummy.Next, dummy.Next
	pre := dummy
	for fast != nil && fast.Next != nil {
		pre = pre.Next
		slow = slow.Next
		fast = fast.Next.Next
	}
	pre.Next = nil
	leftTree := sortedListToBST(dummy.Next)
	rightTree := sortedListToBST(slow.Next)
	slow.Next = nil
	return &TreeNode{
		Val:   slow.Val,
		Left:  leftTree,
		Right: rightTree,
	}
}
