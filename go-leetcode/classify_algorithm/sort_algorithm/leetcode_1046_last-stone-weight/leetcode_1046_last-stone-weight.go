package leetcode_1046_last_stone_weight

import "container/heap"

// 1046. 最后一块石头的重量
// https://leetcode.cn/problems/last-stone-weight/

type hp []int

func (h hp) Len() int {
	return len(h)
}

func (h hp) Less(i, j int) bool {
	return h[i] >= h[j]
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

// lastStoneWeight 大顶堆(优先队列)
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
func lastStoneWeight(stones []int) int {
	n := len(stones)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return stones[0]
	}

	h := hp(stones)
	heap.Init(&h)
	for len(h) > 1 {
		x := heap.Pop(&h).(int)
		y := heap.Pop(&h).(int)
		leave := x - y
		if leave > 0 {
			heap.Push(&h, leave)
		}
	}
	if len(h) == 1 {
		return heap.Pop(&h).(int)
	}
	return 0
}
