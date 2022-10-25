package leetcode_1167_minimum_cost_to_connect_sticks

import "container/heap"

// 1167. 连接棒材的最低费用
// https://leetcode.cn/problems/minimum-cost-to-connect-sticks/

type hp []int

func (h hp) Len() int {
	return len(h)
}

func (h hp) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h hp) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *hp) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *hp) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}

// connectSticks 优先队列
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
// 思路：越短的棒材要被使用越多次
func connectSticks(sticks []int) int {
	h := hp(sticks)
	heap.Init(&h)
	sum := 0
	for h.Len() > 1 {
		v1 := heap.Pop(&h).(int)
		v2 := heap.Pop(&h).(int)
		v := v1 + v2
		sum += v
		heap.Push(&h, v)
	}
	return sum
}
