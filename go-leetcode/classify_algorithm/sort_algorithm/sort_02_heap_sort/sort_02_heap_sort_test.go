package sort_02_heap_sort

import (
	"fmt"
	"testing"
)

// 参考文章：https://www.cnblogs.com/onepixel/articles/7674659.html#!comments

// heapSort 大顶堆排序
// 时间复杂度(平均): O(n*log(n))
// 时间复杂度(最坏): O(n*log(n))
// 时间复杂度(最好): O(n*log(n))
// 空间复杂度: O(1): O(1)
// 稳定性: 不稳定
// 思路：利用堆这种数据结构所设计的一种排序算法。堆积是一个近似完全二叉树的结构，并同时满足堆积的性质：即子结点的键值或索引总是小于（或者大于）它的父节点。
//	a.构建大顶堆
//	b.调整堆
func heapSort(arr []int) {
	arrLen := len(arr)
	buildMaxHeap(arr)
	for i := arrLen - 1; i >= 0; i-- {
		heapify(arr, i)
	}
}

// buildMaxHeap 建堆
func buildMaxHeap(arr []int) {
	arrLen := len(arr)
	for i := arrLen / 2; i >= 0; i-- {
		heapify(arr, i)
	}
}

// heapify 调整堆
func heapify(arr []int, i int) {
	left := 2*i + 1
	right := 2*i + 2
	largest := i
	arrLen := len(arr)
	if left < arrLen && arr[left] > arr[largest] {
		largest = left
	}
	if right < arrLen && arr[right] > arr[largest] {
		largest = right
	}
	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		heapify(arr, largest)
	}
}

func TestHeapSort(t *testing.T) {
	arr := []int{5, 6, 7, 4, 2}
	heapSort(arr)
	fmt.Println(arr) // [7 6 5 4 2]

	arr2 := []int{2, 3, 2}
	heapSort(arr2)
	fmt.Println(arr2) // [3 2 2]
}
