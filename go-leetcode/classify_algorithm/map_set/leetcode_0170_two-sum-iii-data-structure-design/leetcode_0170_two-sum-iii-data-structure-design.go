package leetcode_0170_two_sum_iii_data_structure_design

// 0170.两数之和 III - 数据结构设计
// https://leetcode.cn/problems/two-sum-iii-data-structure-design/

// TwoSum 哈希表
// 时间复杂度：
// - add O(1)
// - find O(n)
// 空间复杂度：O(n)
type TwoSum struct {
	m map[int]int
}

func Constructor() TwoSum {
	return TwoSum{
		m: map[int]int{},
	}
}

func (this *TwoSum) Add(number int) {
	if v, ok := this.m[number]; ok {
		this.m[number] = v + 1
	} else {
		this.m[number] = 1
	}
}

func (this *TwoSum) Find(value int) bool {
	for k, v := range this.m {
		complement := value - k
		if complement != k {
			if _, ok := this.m[complement]; ok {
				return true
			}
		} else {
			if v > 1 {
				return true
			}
		}
	}
	return false
}
