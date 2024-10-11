package leetcode_1209_remove_all_adjacent_duplicates_in_string_ii

// 1209.删除字符串中的所有相邻重复项 II🌟
// https://leetcode.cn/problems/remove-all-adjacent-duplicates-in-string-ii/

// removeDuplicates 记忆计数（会超时）
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func removeDuplicates(s string, k int) string {
	n := len(s)
	if n < k {
		return s
	}

	count := make([]int, n)
	ss := s
	for i := 0; i < len(ss); i++ {
		if i == 0 || ss[i] != ss[i-1] {
			count[i] = 1
		} else {
			count[i] = count[i-1] + 1
			if count[i] == k {
				ss = ss[:i-k+1] + ss[i+1:]
				temp := append([]int(nil), count[:i-k+1]...)
				temp = append(temp, count[i+1:]...)
				i = i - k
			}
		}
	}
	return ss
}
