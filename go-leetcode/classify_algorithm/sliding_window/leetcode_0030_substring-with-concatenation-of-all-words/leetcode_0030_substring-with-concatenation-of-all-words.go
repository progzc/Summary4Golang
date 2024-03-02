package leetcode_0030_substring_with_concatenation_of_all_words

// 0030. 串联所有单词的子串
// https://leetcode.cn/problems/substring-with-concatenation-of-all-words/description/

// 相似点：
// 0438. 找到字符串中所有字母异位词
// https://leetcode.cn/problems/find-all-anagrams-in-a-string/description/

// findSubstring_2 双哈希map
// 时间复杂度: O(n*m)
// 空间复杂度: O(m)
func findSubstring_2(s string, words []string) []int {
	wordNum := len(words)
	wordLen := len(words[0])

	// 第1个哈希：存所有单词
	allWords := make(map[string]int)
	for _, word := range words {
		allWords[word] = allWords[word] + 1
	}

	var ans []int
	for i := 0; i <= len(s)-wordNum*wordLen; i++ {
		// 第2个哈希：存当前扫描的字符串含有的单词
		hasWords := make(map[string]int)
		count := 0
		sub := s[i : i+wordNum*wordLen] // 对比字串
		for j := 0; j <= len(sub)-wordLen; j += wordLen {
			source := sub[j : j+wordLen]
			if needCnt, ok := allWords[source]; ok {
				hasWords[source] += 1
				if hasWords[source] > needCnt {
					break
				}
				count++
			} else {
				break
			}
		}
		if count == wordNum {
			ans = append(ans, i)
		}
	}
	return ans
}

// findSubstring_3 dfs+哈希map （会超时）
// 时间复杂度: O(n*n!)
// 空间复杂度: O(n)
func findSubstring_3(s string, words []string) []int {
	if len(words) == 0 || len(words[0]) == 0 ||
		len(s) < len(words)*len(words[0]) {
		return []int{}
	}

	var (
		dfs func(pos int, cur string)
		has = map[string]bool{}
		//perms []string
	)

	dfs = func(i int, cur string) {
		if i == len(words) {
			if _, ok := has[cur]; !ok {
				//perms = append(perms,cur)
				has[cur] = true
			}
			return
		}
		for j := i; j < len(words); j++ {
			swap(words, i, j)
			temp := cur + words[i]
			dfs(i+1, temp)
			swap(words, j, i)
		}
		return
	}
	dfs(0, "")

	//fmt.Printf("has: %+v\n",has)
	var ans []int
	step := len(words) * len(words[0])
	for i := 0; i <= len(s)-step; i++ {
		//fmt.Printf("i: %d; sub: %v\n",i, s[i:i+step])
		if has[s[i:i+step]] {
			ans = append(ans, i)
		}
	}
	return ans
}

func swap(words []string, i, j int) {
	words[i], words[j] = words[j], words[i]
}
