package leetcode_0165_compare_version_numbers

import (
	"strconv"
	"strings"
)

// 0165. 比较版本号
// https://leetcode.cn/problems/compare-version-numbers/description/

// compareVersion
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func compareVersion(version1 string, version2 string) int {
	vv1 := strings.Split(version1, ".")
	vv2 := strings.Split(version2, ".")
	n1 := len(vv1)
	n2 := len(vv2)
	i, j := 0, 0
	for i < n1 || j < n2 {
		s1 := "0"
		s2 := "0"
		if i < n1 {
			s1 = strings.TrimPrefix(vv1[i], "0")
		}
		if j < n2 {
			s2 = strings.TrimPrefix(vv2[j], "0")
		}
		if s1 == "" {
			s1 = "0"
		}
		if s2 == "" {
			s2 = "0"
		}
		x1, _ := strconv.Atoi(s1)
		x2, _ := strconv.Atoi(s2)
		if x1 > x2 {
			return 1
		} else if x1 < x2 {
			return -1
		}
		i++
		j++
	}
	return 0
}
