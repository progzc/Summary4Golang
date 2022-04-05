package sort_01_quick_sort

import (
	"fmt"
	"testing"
)

// 参考文章：https://www.cnblogs.com/onepixel/articles/7674659.html#!comments

// quickSort 快速排序
// 时间复杂度(平均): O(n*log(n))
// 时间复杂度(最坏): O(n^2)
// 时间复杂度(最好): O(n*log(n))
// 空间复杂度: O(n*log(n))
// 稳定性: 不稳定
// 思路：通过一趟排序将待排记录分隔成独立的两部分，其中一部分记录的关键字均比另一部分的关键字小，则可分别对这两部分记录继续进行排序，以达到整个序列有序。
//	a.挑选基准：从数列中挑出一个元素，称为"基准"(pivot)
//	b.分治：重新排序数列，所有元素比基准值小的摆放在基准前面，所有元素比基准值大的摆在基准的后面（相同的数可以到任一边）。
// 	  在这个分区退出之后，该基准就处于数列的中间位置。这个称为分区（partition）操作。
//	c.递归：把小于基准值元素的子数列和大于基准值元素的子数列排序。
func quickSort(arr []int, left, right int) {
	if len(arr) < 2 || left < 0 || right >= len(arr) || left >= right {
		return
	}

	pdx := partition(arr, left, right)
	quickSort(arr, left, pdx-1)
	quickSort(arr, pdx+1, right)
}

func partition(arr []int, left, right int) int {
	var (
		pivot = left
		index = pivot + 1
	)
	for i := index; i <= right; i++ {
		if arr[i] < arr[pivot] {
			arr[i], arr[index] = arr[index], arr[i]
			index++
		}
	}
	arr[pivot], arr[index-1] = arr[index-1], arr[pivot]
	return index - 1
}

func TestQuickSort(t *testing.T) {
	arr := []int{5, 6, 7, 4, 2}
	quickSort(arr, 0, len(arr)-1)
	fmt.Println(arr) // [2 4 5 6 7]

	arr2 := []int{2, 3, 2}
	quickSort(arr2, 0, len(arr2)-1)
	fmt.Println(arr2) // [2 2 3]
}
