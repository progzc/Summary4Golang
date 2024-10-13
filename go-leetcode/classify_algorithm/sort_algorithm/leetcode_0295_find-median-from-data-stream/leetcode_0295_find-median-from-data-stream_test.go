package leetcode_0295_find_median_from_data_stream

import (
	"container/heap"
	"fmt"
	"testing"
)

// 0295. 数据流的中位数🌟
// https://leetcode.cn/problems/find-median-from-data-stream

// MedianFinder 大顶堆+小顶堆
// 时间复杂度: log(nlog(n))
// 空间复杂度: O(n)
// 思路: 前半部分用大顶堆存，后半部分用小顶堆存；始终保证 len(大顶堆) == len(小顶堆) 或 len(大顶堆)+1 == len(小顶堆)
type MedianFinder struct {
	maxIHeap *MaxIHeap
	minIHeap *MinIHeap
}

func Constructor() MedianFinder {
	return MedianFinder{
		maxIHeap: &MaxIHeap{},
		minIHeap: &MinIHeap{},
	}
}

func (this *MedianFinder) AddNum(num int) {
	if this.maxIHeap.Len() == this.minIHeap.Len() {
		if this.maxIHeap.Len() == 0 || num >= (*this.maxIHeap)[0] {
			heap.Push(this.minIHeap, num)
		} else {
			heap.Push(this.minIHeap, heap.Pop(this.maxIHeap).(int))
			heap.Push(this.maxIHeap, num)
		}
	} else {
		if num <= (*this.minIHeap)[0] {
			heap.Push(this.maxIHeap, num)
		} else {
			heap.Push(this.maxIHeap, heap.Pop(this.minIHeap).(int))
			heap.Push(this.minIHeap, num)
		}
	}
}

func (this *MedianFinder) FindMedian() float64 {
	if this.maxIHeap.Len() == this.minIHeap.Len() {
		return float64((*this.maxIHeap)[0]+(*this.minIHeap)[0]) / 2.0
	}
	return float64((*this.minIHeap)[0])
}

type MaxIHeap []int

func (h *MaxIHeap) Len() int { return len(*h) }

func (h *MaxIHeap) Swap(i, j int) { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MaxIHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] } // 大顶推

func (h *MaxIHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MaxIHeap) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}

type MinIHeap []int

func (h *MinIHeap) Len() int { return len(*h) }

func (h *MinIHeap) Swap(i, j int) { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MinIHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] } // 小顶推

func (h *MinIHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MinIHeap) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}

func TestMedianFinder(t *testing.T) {
	c := Constructor()
	c.AddNum(-1)
	fmt.Printf("Find: %v\n", c.FindMedian()) // Find: -1

	c.AddNum(-2)
	fmt.Printf("Find: %v\n", c.FindMedian()) // Find: -1.5

	c.AddNum(-3)
	fmt.Printf("Find: %v\n", c.FindMedian()) // Find: -2

	c.AddNum(-4)
	fmt.Printf("Find: %v\n", c.FindMedian()) // Find: -2.5

	c.AddNum(-5)
	fmt.Printf("Find: %v\n", c.FindMedian()) // Find: -3
}
