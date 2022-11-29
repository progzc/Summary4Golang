package leetcode_0287_find_the_duplicate_number

// 287. 寻找重复数
// https://leetcode.cn/problems/find-the-duplicate-number/

// findDuplicate 冒泡排序（超时）
// 时间复杂度: O(n^2)
// 空间复杂度: O(1)
// 缺点:
//	未充分利用这个条件: 给定一个包含 n+1 个整数的数组 nums, 其数字都在 [1,n] 范围内（包括 1 和 n）
func findDuplicate(nums []int) int {
	n := len(nums)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if nums[j] < nums[i] {
				nums[i], nums[j] = nums[j], nums[i]
			}
		}
		if i > 0 && nums[i] == nums[i-1] {
			return nums[i]
		}
	}
	return -1
}

// findDuplicate_2 环形链表（快慢指针法）
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：同下面这道题
// 	0142.环形链表 II
//	https://leetcode-cn.com/problems/linked-list-cycle-ii/
func findDuplicate_2(nums []int) int {
	n := len(nums)
	slow, fast := 0, 0
	for fast < n && nums[fast] < n {
		slow = nums[slow]
		fast = nums[nums[fast]]
		if fast == slow {
			p := 0
			for p != slow {
				p = nums[p]
				slow = nums[slow]
			}
			return p
		}
	}
	return 0
}
