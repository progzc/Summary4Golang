package leetcode_0725_split_linked_list_in_parts

// 0725.分隔链表
// https://leetcode-cn.com/problems/split-linked-list-in-parts/

type ListNode struct {
	Val  int
	Next *ListNode
}

// splitListToParts 一般思路
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 缺点：写得过于复杂了
func splitListToParts(head *ListNode, k int) []*ListNode {
	var (
		count = 0
		ans   []*ListNode
	)

	// 统计数量
	for node := head; node != nil; count++ {
		node = node.Next
	}

	// 如果链表长度小于k
	if count <= k {
		for node := head; node != nil; {
			ans = append(ans, node)
			next := node.Next
			node.Next = nil
			node = next
		}
		for i := 0; i < k-count; i++ {
			ans = append(ans, nil)
		}
		return ans
	}

	// 否则，链表长度大于k
	m, n := count/k, count%k

	first := head
	for i := 0; i < k; i++ {
		dummy := &ListNode{Next: first}
		node := dummy
		for j := 0; j < m; j++ {
			node = node.Next
		}
		ans = append(ans, dummy.Next)

		if n > 0 {
			n--
			node = node.Next
		}
		nextFirst := node.Next
		node.Next = nil
		first = nextFirst
	}
	return ans
}

// splitListToParts_2 一般思路（优化步骤）
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：记m,n:=count/k,count%k，则在分隔成的k个部分中，前n个部分的长度各为m+1，其余每个部分的长度各为m。
// 注意事项：要简化代码
func splitListToParts_2(head *ListNode, k int) []*ListNode {
	count := 0
	for node := head; node != nil; node = node.Next {
		count++
	}

	m, n := count/k, count%k
	ans := make([]*ListNode, k) // 这样超出的部分就是nil了
	for i, curr := 0, head; i < k && curr != nil; i++ {
		ans[i] = curr
		size := m
		if i < n {
			size++
		}
		// 注意事项：这里j=1开始,应为curr本身已经算成一个节点了
		// 或者 for j := 0; j < size-1; j++
		for j := 1; j < size; j++ {
			curr = curr.Next
		}
		// 下面三行可以简写成：
		// curr, curr.Next = curr.Next, nil
		// 或：
		// curr.Next, curr = nil, curr.Next

		// 但是在力扣上分开写效率更高（表现在分开写消耗更低的内存）?
		// 原因在于：golang给多个字段赋值时会为每个字段创建temp
		// 	参见: http://www.javashuo.com/article/p-mmtmmmoz-ho.html
		// 	例如：curr, curr.Next = curr.Next, nil
		// 	等价于：
		//		temp1 := curr
		//		temp2 := curr.Next
		//		curr = temp2
		//		curr.Next = nil
		next := curr.Next
		curr.Next = nil // 这一步很关键，需要断开
		curr = next
	}
	return ans
}
