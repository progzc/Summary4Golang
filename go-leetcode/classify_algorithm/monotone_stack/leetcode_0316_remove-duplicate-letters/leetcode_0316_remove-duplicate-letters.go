package leetcode_0316_remove_duplicate_letters

// 0316. 去除重复字母
// https://leetcode.cn/problems/remove-duplicate-letters/

// removeDuplicateLetters 单调递增(非严格)栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	i)先遍历字符串，记录每个元素出现次数
//	ii)接着元素挨个入栈, 维持栈的核心要素是：保持栈内所有元素都小于即将进栈的元素，即字典序；
//	   如果不满足，则需要把栈内元素出栈，直到满足要求为止。
func removeDuplicateLetters(s string) string {
	var (
		// count 统计每个字符的出现次数
		count [26]int
		// exist 是否在栈中
		exist [26]bool
		// stack 单调栈
		stack []byte
	)

	for _, c := range []byte(s) {
		count[c-'a']++
	}

	for _, c := range []byte(s) {
		// 若栈中已有c,跳过
		if exist[c-'a'] {
			// 同时减少这个字符出现的次数
			count[c-'a']--
			continue
		}
		// 出栈的核心判断要素：
		// 1. 栈里有元素
		// 2. 栈顶元素大于当前元素c
		// 3. 栈顶元素在后续出现
		for len(stack) > 0 && stack[len(stack)-1] > c && count[stack[len(stack)-1]-'a'] > 0 {
			// 进入这里说明栈顶元素大于当前元素，所以不符合字典序，要把栈顶元素出栈
			// 【重要】for 循环确保栈里面保存的都是比当前元素小的，因为大的都被pop掉了
			// 标记为栈中不含栈顶元素
			exist[stack[len(stack)-1]-'a'] = false
			// 删除栈顶元素
			stack = stack[:len(stack)-1]
		}
		// 添加新字符
		stack = append(stack, c)
		// 减少该字符出现次数
		count[c-'a']--
		// 标记栈中有该字符
		exist[c-'a'] = true
	}
	return string(stack)
}
