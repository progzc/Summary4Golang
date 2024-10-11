package leetcode_2716_minimize_string_length

// 2716. 最小化字符串长度
// https://leetcode.cn/problems/minimize-string-length

// minimizedStringLength 栈
// 怀疑力扣官网的题目或答案有问题
// 输入 s="baadccab", 输出应该是badcab，长度不应该是 6 吗？怎么是 4
func minimizedStringLength(s string) int {
	var stack []byte
	for i := 0; i < len(s); i++ {
		if len(stack) > 0 && stack[len(stack)-1] == s[i] {
			continue
		} else {
			stack = append(stack, s[i])
		}
	}
	return len(stack)
}
