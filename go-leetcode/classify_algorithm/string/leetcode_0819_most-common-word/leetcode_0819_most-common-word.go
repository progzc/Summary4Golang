package leetcode_0819_most_common_word

import "strings"

// 0819.最常见的单词
// https://leetcode-cn.com/problems/most-common-word/

// mostCommonWord 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func mostCommonWord(paragraph string, banned []string) string {
	s := strings.ToLower(paragraph)

	bm := map[string]struct{}{}
	for _, ban := range banned {
		bm[ban] = struct{}{}
	}

	ans, max := "", 0
	cnt := map[string]int{}

	i, j := 0, 0
	for i < len(s) {
		for i < len(s) && (s[i] < 'a' || s[i] > 'z') {
			i++
		}
		if i == len(s) {
			break
		}

		j = i
		for j < len(s) && (s[j] >= 'a' && s[j] <= 'z') {
			j++
		}

		key := s[i:j]
		if _, ok := bm[key]; len(key) > 0 && !ok {
			cnt[key]++
			if cnt[key] > max {
				ans = key
				max = cnt[key]
			}
		}
		i = j + 1
	}
	return ans
}
