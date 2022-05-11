package leetcode_0378_kth_smallest_element_in_a_sorted_matrix

import (
	"container/heap"
	"sort"
)

// 0378.有序矩阵中第K小的元素
// https://leetcode-cn.com/problems/kth-smallest-element-in-a-sorted-matrix/

// kthSmallest 直接排序
// 时间复杂度: O(n^2*log(n))
// 空间复杂度: O(n^2)
// 思路：转化为一维数组再排序
// 缺点：没有充分利用数组在行/列上的有序性
func kthSmallest(matrix [][]int, k int) int {
	rows, columns := len(matrix), len(matrix[0])
	all := make([]int, rows*columns)
	idx := 0
	for _, row := range matrix {
		for _, num := range row {
			all[idx] = num
			idx++
		}
	}
	sort.Ints(all)
	return all[k-1]
}

// kthSmallest_2 归并排序（使用堆来完成）
// 时间复杂度: O(k*log(n))
// 空间复杂度: O(n)
// 思路：转化为n个数组进行归并排序
// 缺点：没有充分利用数组在列上的有序性
func kthSmallest_2(matrix [][]int, k int) int {
	h := &hp{}
	// 将每行的第一个元素放到小顶堆中
	for i := 0; i < len(matrix); i++ {
		heap.Push(h, [3]int{matrix[i][0], i, 0})
	}

	for i := 0; i < k-1; i++ {
		now := heap.Pop(h).([3]int)
		if now[2] < len(matrix[0])-1 {
			heap.Push(h, [3]int{
				matrix[now[1]][now[2]+1],
				now[1],
				now[2] + 1,
			})
		}
	}
	return heap.Pop(h).([3]int)[0]
}

type hp [][3]int

func (h hp) Len() int { return len(h) }

func (h hp) Less(i, j int) bool { return h[i][0] < h[j][0] } // 小顶堆

func (h hp) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *hp) Push(x interface{}) { *h = append(*h, x.([3]int)) }

func (h *hp) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}

// kthSmallest_3 二维二分搜索（充分利用行列均有序）
// 时间复杂度: O(n*log(r-l))
// 空间复杂度: O(1)
// 思路：采用二维二分搜索。从 左上 到 右下 是递增的
func kthSmallest_3(matrix [][]int, k int) int {
	n := len(matrix)
	left, right := matrix[0][0], matrix[n-1][n-1]
	for left < right {
		mid := left + (right-left)/2
		if check(matrix, mid, k, n) {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left
}

func check(matrix [][]int, mid, k, n int) bool {
	i, j := n-1, 0
	num := 0
	for i >= 0 && j < n {
		if matrix[i][j] <= mid {
			num += i + 1
			j++
		} else {
			i--
		}
	}
	return num >= k
}
