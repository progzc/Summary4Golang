package leetcode_0225_implement_stack_using_queues

// 0225.用队列实现栈
// https://leetcode-cn.com/problems/implement-stack-using-queues/
type MyStack struct {
	s1 []int
}

func Constructor() MyStack {
	return MyStack{
		s1: []int{},
	}
}

func (this *MyStack) Push(x int) {
	this.s1 = append(this.s1, x)
}

func (this *MyStack) Pop() int {
	v := this.s1[len(this.s1)-1]
	this.s1 = this.s1[:len(this.s1)-1]
	return v
}

func (this *MyStack) Top() int {
	return this.s1[len(this.s1)-1]
}

func (this *MyStack) Empty() bool {
	return len(this.s1) == 0
}
