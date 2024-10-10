package tal_20241010_tree_layer

import (
	"fmt"
	"testing"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// TestZigzagLevelOrder
// 好未来面试题：层序打印树，奇数行顺序打印，偶数行倒序打印。
func TestZigzagLevelOrder(t *testing.T) {
	root := &TreeNode{Val: 0}
	root.Left = &TreeNode{Val: 1}
	root.Right = &TreeNode{Val: 2}
	root.Left.Left = &TreeNode{Val: 3}
	root.Left.Right = &TreeNode{Val: 4}
	fmt.Println(zigzagLevelOrder(root))
}

// zigzagLevelOrder
func zigzagLevelOrder(root *TreeNode) [][]int {
	var ans [][]int
	if root == nil {
		return ans
	}
	var stack []*TreeNode
	stack = append(stack, root)
	flag := true // 奇为 true
	for len(stack) > 0 {
		size := len(stack)
		var layer []int
		for i := 0; i < size; i++ {
			node := stack[0]
			stack = stack[1:]
			if flag {
				layer = append(layer, node.Val) // 奇数尾插
			} else {
				layer = append([]int{node.Val}, layer...) // 偶数头插
			}
			if node.Left != nil {
				stack = append(stack, node.Left)
			}
			if node.Right != nil {
				stack = append(stack, node.Right)
			}
		}
		flag = !flag
		ans = append(ans, layer)
	}
	return ans
}

func TestSlice(t *testing.T) {
	nums := []int{1, 2, 3}
	solve1(nums)
	fmt.Println(nums) // [5 2 3]

	nums = []int{1, 2, 3}
	solve2(nums)
	fmt.Println(nums) // [1 2 3]

	nums = []int{1, 2, 3}
	solve3(nums)
	fmt.Println("cap1:", cap(nums)) // cap1: 3
	fmt.Println(nums)               // [1 2 3]

	nums = make([]int, 0, 6)
	nums = append(nums, 1, 2, 3)
	solve4(nums)
	fmt.Println("cap2:", cap(nums)) // cap2: 6
	fmt.Println(nums)               // [7 2 3]
}

func solve1(nums []int) {
	nums[0] = 5
}

func solve2(nums []int) {
	nums = append(nums, 4, 5, 6)
}

func solve3(nums []int) {
	nums = append(nums, 4, 5, 6)
	fmt.Println("solve3 cap:", cap(nums)) // solve3 cap: 6
	nums[0] = 7
}

func solve4(nums []int) {
	nums = append(nums, 4, 5, 6)
	fmt.Println("solve4 cap:", cap(nums)) // solve4 cap: 6
	nums[0] = 7
}
