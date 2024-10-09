package leetcode_0155_min_stack

// 0155.最小栈🌟
// https://leetcode-cn.com/problems/min-stack/

// MinStack 两个栈
// 时间复杂度: O(1)
// 空间复杂度: O(n)
type MinStack struct {
	m []int // 正常栈
	n []int // 最小栈
}

func Constructor() MinStack {
	return MinStack{
		m: make([]int, 0),
		n: make([]int, 0),
	}
}

func (this *MinStack) Push(val int) {
	this.m = append(this.m, val)
	if len(this.n) == 0 {
		this.n = append(this.n, val)
		return
	}
	curMin := this.n[len(this.n)-1]
	if curMin < val {
		this.n = append(this.n, curMin)
	} else {
		this.n = append(this.n, val)
	}
}

func (this *MinStack) Pop() {
	this.m = this.m[:len(this.m)-1]
	this.n = this.n[:len(this.n)-1]

}

func (this *MinStack) Top() int {
	return this.m[len(this.m)-1]
}

func (this *MinStack) GetMin() int {
	return this.n[len(this.m)-1]
}

//type MinStack struct {
//	stack1 []int
//	stack2 []int
//}
//
//func Constructor() MinStack {
//	m := MinStack{
//		stack1: make([]int, 0),
//		stack2: make([]int, 0),
//	}
//	return m
//}
//
//func (this *MinStack) Push(val int) {
//	this.stack1 = append(this.stack1, val)
//	if len(this.stack2) == 0 || val <= this.stack2[len(this.stack2)-1] {
//		this.stack2 = append(this.stack2, val)
//	}
//}
//
//func (this *MinStack) Pop() {
//	val := this.stack1[len(this.stack1)-1]
//	this.stack1 = this.stack1[:len(this.stack1)-1]
//	if val == this.stack2[len(this.stack2)-1] {
//		this.stack2 = this.stack2[:len(this.stack2)-1]
//	}
//}
//
//func (this *MinStack) Top() int {
//	return this.stack1[len(this.stack1)-1]
//}
//
//func (this *MinStack) GetMin() int {
//	return this.stack2[len(this.stack2)-1]
//}
