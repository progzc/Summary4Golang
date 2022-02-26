package leetcode_0703_kth_largest_element_in_a_stream

import (
	"container/heap"
	"sort"
)

// 0703.数据流中的第 K 大元素
// https://leetcode-cn.com/problems/kth-largest-element-in-a-stream/

type KthLargest struct {
	sort.IntSlice
	k int
}

func Constructor(k int, nums []int) KthLargest {
	kthLargest := KthLargest{
		k: k,
	}
	for _, v := range nums {
		kthLargest.Add(v)
	}
	return kthLargest
}

func (this *KthLargest) Add(val int) int {
	// heap 默认是小顶堆（因为sort.IntSlice默认是升序）
	heap.Push(this, val)
	if this.Len() > this.k {
		heap.Pop(this)
	}
	return this.IntSlice[0]
}

func (this *KthLargest) Push(v interface{}) {
	this.IntSlice = append(this.IntSlice, v.(int))
}

func (this *KthLargest) Pop() interface{} {
	v := this.IntSlice[len(this.IntSlice)-1]
	this.IntSlice = this.IntSlice[:len(this.IntSlice)-1]
	return v
}
