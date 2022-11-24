package futu_20221124_move_right

// 题目:
// 将一个数组的所有元素向右移动若干单位，并把数组右侧溢出的元素填补
// 在数组左侧的空缺中，这种经操作称为数组的循环平移。
//
// 给你一个不小于 3 个元素的数组 a，已知 a 是从一个有序且不包含
// 重复元素的数组平移 k(k 大于等于 0 且小于数组长度)个单位而来；
// 请写一个函数，输入 int 类型数组 a，返回 k 的值。
//
// 例如，对于数组 a = {5, 1, 2, 3, 4}，它由有序数组
// {1, 2, 3, 4, 5} 循环平移 1 个单位而来，因此 k = 1。

// 力扣类似题:
// 153. 寻找旋转排序数组中的最小值
// https://leetcode.cn/problems/find-minimum-in-rotated-sorted-array/

// translate 二分查找
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
func translate(nums []int) int {
	low, high := 0, len(nums)-1
	for low < high {
		mid := low + (high-low)/2
		if nums[mid] < nums[high] {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}
