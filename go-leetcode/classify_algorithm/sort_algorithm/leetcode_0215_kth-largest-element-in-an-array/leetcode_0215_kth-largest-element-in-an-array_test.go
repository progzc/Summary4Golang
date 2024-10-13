package leetcode_0215_kth_largest_element_in_an_array

import (
	"container/heap"
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

// 0215.数组中的第K个最大元素🌟
// https://leetcode-cn.com/problems/kth-largest-element-in-an-array/
// 复习巩固：快速排序+堆排序

// findKthLargest 使用标准库
// 时间复杂度: O(n*log(n))
// 空间复杂度: O(log(n))
func findKthLargest(nums []int, k int) int {
	sort.Ints(nums)
	return nums[len(nums)-k]
}

// findKthLargest_2 使用快速排序
// 时间复杂度: O(n*log(n))
// 空间复杂度: O(log(n))
func findKthLargest_2(nums []int, k int) int {
	rand.Seed(time.Now().UnixNano())
	return quickSelect(nums, 0, len(nums)-1, len(nums)-k)
}

func quickSelect(a []int, l, r, index int) int {
	q := randomPartition(a, l, r)
	if q == index {
		return a[q]
	} else if q < index {
		return quickSelect(a, q+1, r, index)
	}
	return quickSelect(a, l, q-1, index)
}

func randomPartition(a []int, l, r int) int {
	i := rand.Int()%(r-l+1) + l
	a[i], a[r] = a[r], a[i]
	return partition(a, l, r)
}

func partition(a []int, l, r int) int {
	x := a[r]
	i := l - 1
	for j := l; j < r; j++ {
		if a[j] <= x {
			i++
			a[i], a[j] = a[j], a[i]
		}
	}
	a[i+1], a[r] = a[r], a[i+1]
	return i + 1
}

// findKthLargest_3 使用堆排序
// 时间复杂度: O(n*log(n))
// 空间复杂度: O(log(n))
func findKthLargest_3(nums []int, k int) int {
	heapSize := len(nums)
	buildMaxHeap(nums, heapSize)
	for i := len(nums) - 1; i >= len(nums)-k+1; i-- {
		nums[0], nums[i] = nums[i], nums[0]
		heapSize--
		maxHeapify(nums, 0, heapSize)
	}
	return nums[0]
}

func buildMaxHeap(a []int, heapSize int) {
	for i := heapSize / 2; i >= 0; i-- {
		maxHeapify(a, i, heapSize)
	}
}

func maxHeapify(a []int, i, heapSize int) {
	l, r, largest := i*2+1, i*2+2, i
	if l < heapSize && a[l] > a[largest] {
		largest = l
	}
	if r < heapSize && a[r] > a[largest] {
		largest = r
	}
	if largest != i {
		a[i], a[largest] = a[largest], a[i]
		maxHeapify(a, largest, heapSize)
	}
}

// -------------------------堆的使用------------------------------
func TestFindKthLargest_4(t *testing.T) {
	nums, k := []int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4
	fmt.Println(findKthLargest_4(nums, k))
}

type hp []int

func (h *hp) Len() int { return len(*h) }

func (h *hp) Less(i, j int) bool { return (*h)[i] > (*h)[j] } // 大顶堆

func (h *hp) Swap(i, j int) { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *hp) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *hp) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}

// findKthLargest_4 采用标准库的堆排序
func findKthLargest_4(nums []int, k int) int {
	h := hp(nums)
	heap.Init(&h)
	for i := 0; i < k-1; i++ {
		heap.Pop(&h)
	}
	return heap.Pop(&h).(int)
}
