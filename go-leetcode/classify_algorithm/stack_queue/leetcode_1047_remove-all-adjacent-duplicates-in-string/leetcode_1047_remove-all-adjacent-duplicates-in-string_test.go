package leetcode_1047_remove_all_adjacent_duplicates_in_string

// 1047. 删除字符串中的所有相邻重复项
// https://leetcode.cn/problems/remove-all-adjacent-duplicates-in-string

// removeDuplicates 栈
// 时间复杂度：O(n)
// 空间复杂度: O(n)
func removeDuplicates(s string) string {
	var stack []byte
	for i := 0; i < len(s); i++ {
		if len(stack) > 0 && stack[len(stack)-1] == s[i] {
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, s[i])
		}
	}
	return string(stack)
}
