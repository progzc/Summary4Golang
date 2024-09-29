package zhihu_1116_1st_merge_two_disorder_array

import "sort"

// 知乎一面算: 合并两个无序数组,变成升序排列,并去重

func combine(nums1, nums2 []int) []int {
	var ans []int
	sort.Slice(nums1, func(i, j int) bool {
		return nums1[i] < nums1[j]
	})
	sort.Slice(nums2, func(i, j int) bool {
		return nums2[i] < nums2[j]
	})

	len1, len2 := len(nums1), len(nums2)

	p1, p2 := 0, 0
	for p1 < len1 && p2 < len2 {
		if nums1[p1] <= nums2[p2] {
			if len(ans) == 0 {
				ans = append(ans, nums1[p1])
			}
			if len(ans) > 0 && ans[len(ans)-1] != nums1[p1] {
				ans = append(ans, nums1[p1])
			}
			p1++
		} else {
			if len(ans) == 0 {
				ans = append(ans, nums2[p2])
			}
			if len(ans) > 0 && ans[len(ans)-1] != nums2[p2] {
				ans = append(ans, nums2[p2])
			}
			p2++
		}
	}

	for p1 < len1 {
		if len(ans) > 0 && ans[len(ans)-1] != nums1[p1] {
			ans = append(ans, nums1[p1])
		}
		p1++
	}
	for p2 < len2 {
		if len(ans) > 0 && ans[len(ans)-1] != nums2[p2] {
			ans = append(ans, nums2[p2])
		}
		p2++
	}
	return ans
}
