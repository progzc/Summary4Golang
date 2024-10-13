package leetcode_0347_top_k_frequent_elements

import (
	"container/heap"
	"sort"
)

// 0347.前 K 个高频元素
// https://leetcode.cn/problems/top-k-frequent-elements/

// topKFrequent 基于快速排序
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
func topKFrequent(nums []int, k int) []int {
	type item struct {
		num int
		cnt int
	}
	n := len(nums)
	m := make(map[int]*item)
	for i := 0; i < n; i++ {
		if v, ok := m[nums[i]]; ok {
			v.cnt++
			m[nums[i]] = v
		} else {
			m[nums[i]] = &item{nums[i], 1}
		}
	}

	var (
		ss  []*item
		ans []int
	)
	for _, value := range m {
		ss = append(ss, value)
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].cnt >= ss[j].cnt
	})

	for i := 0; i < len(ss) && i < k; i++ {
		ans = append(ans, ss[i].num)
	}
	return ans
}

type IHeap [][2]int

func (h *IHeap) Push(x any) {
	*h = append(*h, x.([2]int))
}

func (h *IHeap) Pop() any {
	if len(*h) == 0 {
		return nil
	}
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}

func (h *IHeap) Len() int { return len(*h) }

func (h *IHeap) Less(i, j int) bool { return (*h)[i][1] > (*h)[j][1] } // 大顶堆

func (h *IHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

// topKFrequent_2 基于大顶推
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
func topKFrequent_2(nums []int, k int) []int {
	var (
		ans []int
		m   = make(map[int]int)
	)
	for _, num := range nums {
		m[num] = m[num] + 1
	}

	h := &IHeap{}
	heap.Init(h)
	for key, value := range m {
		heap.Push(h, [2]int{key, value})
	}
	for i := 0; i < k; i++ {
		ans = append(ans, heap.Pop(h).([2]int)[0])
	}
	return ans
}

type IHeap2 struct {
	IHeap
}

func (h *IHeap2) Less(i, j int) bool { return !h.IHeap.Less(i, j) } // 小顶堆

// topKFrequent_3 基于小顶推
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
func topKFrequent_3(nums []int, k int) []int {
	var m = make(map[int]int)
	for _, num := range nums {
		m[num] = m[num] + 1
	}

	h := &IHeap2{}
	heap.Init(h)
	for key, value := range m {
		heap.Push(h, [2]int{key, value})
		if h.Len() > k {
			heap.Pop(h)
		}
	}

	ans := make([]int, k)
	for i := 0; i < k; i++ {
		ans[k-i-1] = heap.Pop(h).([2]int)[0]
	}
	return ans
}
