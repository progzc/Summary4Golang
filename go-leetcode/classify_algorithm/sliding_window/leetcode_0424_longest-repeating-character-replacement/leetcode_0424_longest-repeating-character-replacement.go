package leetcode_0424_longest_repeating_character_replacement

// 0424. 替换后的最长重复字符
// https://leetcode.cn/problems/longest-repeating-character-replacement/

// 与 https://leetcode.cn/problems/max-consecutive-ones-iii/ 类似

// characterReplacement 滑动窗口
// 时间复杂度: O(n)
// 空间复杂度: O(A),A=26
// 对比暴力法，滑动窗口优化的点：
//	a.如果找到了一个长度为 L 且替换 k 个字符以后全部相等的子串，就没有必要考虑长度小于等于 L 的子串，
// 	  因为题目只让我们找到符合题意的最长的长度。
//	b.如果找到了一个长度为 L 且替换 k 个字符以后不能全部相等的子串，左边界相同、
// 	  长度更长的子串一定不符合要求（原因我们放在最后说）。
func characterReplacement(s string, k int) int {
	n := len(s)
	if n <= k {
		return n
	}
	ans, maxFreq := 0, 0
	freq := make([]int, 26)
	for l, r := 0, 0; r < n; r++ {
		freq[s[r]-'A']++
		maxFreq = max(maxFreq, freq[s[r]-'A'])
		// 这里很巧妙：我们如何判断一个字符串改变 K 个字符，能够变成一个连续串，
		// 如果当前字符串中的出现次数最多的字母个数 +K 大于串长度，那么这个串就是满足条件的。
		// 这里可以不使用for循环
		if r-l+1 > maxFreq+k {
			// 此时说明k不够用
			// 把其它不是最多出现的字符替换以后，都不能填满这个滑动的窗口，这个时候须要考虑左边界向右移动
			// 移出滑动窗口的时候，频数数组须要相应地做减法
			freq[s[l]-'A']--
			l++
		} else {
			ans = max(ans, r-l+1)
		}
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
