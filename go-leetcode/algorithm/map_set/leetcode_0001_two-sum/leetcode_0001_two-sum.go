package leetcode_0001_two_sum

// 0001.两数之和
// https://leetcode-cn.com/problems/two-sum/

// twoSum 散列表
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func twoSum(nums []int, target int) []int {
	m := map[int]int{}
	for i, num := range nums {
		if v, ok := m[target-num]; ok {
			return []int{i, v}
		}
		m[num] = i
	}
	return nil
}
