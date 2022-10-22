package leetcode_2024_maximize_the_confusion_of_an_exam

// 2024. 考试的最大困扰度
// https://leetcode.cn/problems/maximize-the-confusion-of-an-exam/

// 与 https://leetcode.cn/problems/longest-repeating-character-replacement/ 基本一样

// maxConsecutiveAnswers 滑动窗口
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func maxConsecutiveAnswers(answerKey string, k int) int {
	n := len(answerKey)
	if n <= k {
		return n
	}

	m := map[byte]int{}
	ans, maxFreq := 0, 0
	for l, r := 0, 0; r < n; r++ {
		m[answerKey[r]]++
		maxFreq = max(maxFreq, m[answerKey[r]])
		// 注意这里的技巧
		for r-l+1 > maxFreq+k {
			m[answerKey[l]]--
			if m[answerKey[l]] == 0 {
				delete(m, answerKey[l])
			}
			l++
		}
		ans = max(ans, r-l+1)
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
