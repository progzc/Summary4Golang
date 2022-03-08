package binary_search

// 关于二分搜索及其4个变种的总结:
// 1.经典的二分搜索法
// 2.4个基本变种:
//	a.查找第一个与 target 相等的元素
//	b.查找最后一个与 target 相等的元素
//	c.查找第一个大于等于 target 的元素
//	d.查找最后一个小于等于 target 的元素
// 3.其他变种:
//	a.在山峰数组中找山峰
//	b.在旋转有序数组中找分界点
// 4.力扣中的经典题目：33、81、153、154、162、852

// binarySearchMatrix 经典的二分搜索法
// 注意事项:
//	a.二分搜索的前提是 有序(一般而言是升序)
// 	b.循环退出条件，注意是 low <= high，而不是 low < high
// 	c.mid 的取值，mid := low + (high-low)/2
//	d.low 和 high 的更新，low = mid + 1，high = mid - 1
func binarySearchMatrix(nums []int, target int) int {
	low, high := 0, len(nums)-1
	for low <= high {
		mid := low + (high-low)>>1
		if nums[mid] == target {
			return mid
		} else if nums[mid] > target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return -1
}

// searchFirstEqualElement 变种1: 查找第一个与 target 相等的元素
func searchFirstEqualElement(nums []int, target int) int {
	low, high := 0, len(nums)-1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if nums[mid] > target {
			high = mid - 1
		} else if nums[mid] < target {
			low = mid + 1
		} else {
			if (mid == 0) || (nums[mid-1] != target) { // 找到第一个与 target 相等的元素
				return mid
			}
			high = mid - 1
		}
	}
	return -1
}

// searchLastEqualElement 变种2: 查找最后一个与 target 相等的元素
func searchLastEqualElement(nums []int, target int) int {
	low, high := 0, len(nums)-1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if nums[mid] > target {
			high = mid - 1
		} else if nums[mid] < target {
			low = mid + 1
		} else {
			if (mid == len(nums)-1) || (nums[mid+1] != target) { // 找到最后一个与 target 相等的元素
				return mid
			}
			low = mid + 1
		}
	}
	return -1
}

// searchFirstGreaterElement 变种3: 查找第一个大于等于 target 的元素
func searchFirstGreaterElement(nums []int, target int) int {
	low, high := 0, len(nums)-1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if nums[mid] >= target {
			if (mid == 0) || (nums[mid-1] < target) { // 找到第一个大于等于 target 的元素
				return mid
			}
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return -1
}

// searchLastLessElement 变种4: searchLastLessElement
func searchLastLessElement(nums []int, target int) int {
	low, high := 0, len(nums)-1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if nums[mid] <= target {
			if (mid == len(nums)-1) || (nums[mid+1] > target) { // 找到最后一个小于等于 target 的元素
				return mid
			}
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}
