package leetcode_0232_implement_queue_using_stacks

// 0232.用栈实现队列
// https://leetcode-cn.com/problems/implement-queue-using-stacks/

type MyQueue struct {
	s1 []int
	s2 []int
}

func Constructor() MyQueue {
	return MyQueue{
		s1: []int{},
		s2: []int{},
	}
}

// 时间复杂度: O(1)
// 空间复杂度: O(n)
func (this *MyQueue) Push(x int) {
	this.s1 = append(this.s1, x)
}

// 时间复杂度: O(n)
// 空间复杂度: O(1)
func (this *MyQueue) Pop() int {
	if len(this.s2) == 0 {
		for len(this.s1) != 0 {
			// 栈1 pop
			v1 := this.s1[len(this.s1)-1]
			this.s1 = this.s1[:len(this.s1)-1]
			// 站2 push
			this.s2 = append(this.s2, v1)
		}
	}

	// 栈2 pop
	v2 := this.s2[len(this.s2)-1]
	this.s2 = this.s2[:len(this.s2)-1]
	return v2
}

// 时间复杂度: O(n)
// 空间复杂度: O(1)
func (this *MyQueue) Peek() int {
	if len(this.s2) == 0 {
		for len(this.s1) != 0 {
			// 栈1 pop
			v1 := this.s1[len(this.s1)-1]
			this.s1 = this.s1[:len(this.s1)-1]
			// 站2 push
			this.s2 = append(this.s2, v1)
		}
	}
	return this.s2[len(this.s2)-1]
}

// 时间复杂度: O(1)
// 空间复杂度: O(1)
func (this *MyQueue) Empty() bool {
	return len(this.s1) == 0 && len(this.s2) == 0
}
