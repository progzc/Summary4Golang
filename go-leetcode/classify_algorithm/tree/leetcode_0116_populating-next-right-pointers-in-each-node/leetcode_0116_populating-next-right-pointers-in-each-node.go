package leetcode_0116_populating_next_right_pointers_in_each_node

// 0116.填充每个节点的下一个右侧节点指针
// https://leetcode-cn.com/problems/populating-next-right-pointers-in-each-node/

type Node struct {
	Val   int
	Left  *Node
	Right *Node
	Next  *Node
}

// connect 迭代
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func connect(root *Node) *Node {
	var (
		queue []*Node
	)
	if root == nil {
		return root
	}
	queue = append(queue, root)
	for len(queue) > 0 {
		size := len(queue)
		cur := &Node{}
		for i := 0; i < size; i++ {
			temp := queue[0]
			queue = queue[1:]
			if i == 0 {
				cur = temp
			} else {
				cur.Next = temp
				cur = cur.Next
			}

			if temp.Left != nil {
				queue = append(queue, temp.Left)
			}
			if temp.Right != nil {
				queue = append(queue, temp.Right)
			}
		}
	}
	return root
}

// TODO 尝试空间复杂度O(1)的解法
