package leetcode_0291_word_pattern_ii

// 0291. 单词规律 II
// https://leetcode.cn/problems/word-pattern-ii/

// wordPatternMatch dfs+map
// 时间复杂度: O(mn)
// 空间复杂度: O(n)
// 注意：双射映射 说明需要两个map
// 推荐题解: https://leetcode.cn/problems/word-pattern-ii/solution/hui-su-hashmapbao-cun-pi-pei-hashsetbao-645u7/
func wordPatternMatch(pattern string, s string) bool {
	var (
		pLen = len(pattern)
		sLen = len(s)
		m    = make(map[byte]string)
		set  = make(map[string]bool)
		dfs  func(idx1, idx2 int) bool
	)
	dfs = func(idx1, idx2 int) bool {
		if idx1 == pLen {
			// 匹配完成
			if idx2 == sLen {
				return true
			} else {
				// s中有部分未被匹配到
				return false
			}
		}

		// 之前添加过
		if v, ok := m[pattern[idx1]]; ok {
			if idx2+len(v) <= sLen && s[idx2:idx2+len(v)] == v {
				return dfs(idx1+1, idx2+len(v))
			} else {
				return false
			}
		}
		// 之前未添加过
		for i := idx2 + 1; i <= sLen; i++ {
			str := s[idx2:i]
			if !set[str] {
				set[str] = true
				m[pattern[idx1]] = str
				if dfs(idx1+1, i) {
					return true
				}
				delete(m, pattern[idx1])
				delete(set, str)
			}
		}
		return false
	}
	return dfs(0, 0)
}
