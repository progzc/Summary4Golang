package leetcode_0347_top_k_frequent_elements

import "sort"

// 0347.前 K 个高频元素
// https://leetcode.cn/problems/top-k-frequent-elements/

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
