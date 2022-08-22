package leetcode_0245_shortest_word_distance_iii

// 0245. 最短单词距离 III
// https://leetcode.cn/problems/shortest-word-distance-iii/
// 注意：word1 可能等于 word2

// shortestWordDistance 一次遍历
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func shortestWordDistance(wordsDict []string, word1 string, word2 string) int {
	ans := len(wordsDict)
	// 分情况讨论
	// 情况一：word1 == word2
	if word1 == word2 {
		pre := -1
		for i, word := range wordsDict {
			if word == word1 {
				if pre >= 0 {
					ans = min(ans, i-pre)
				}
				pre = i
			}
		}
	} else {
		// 情况二：word1 != word2
		index1, index2 := -1, -1
		for i, word := range wordsDict {
			if word == word1 {
				index1 = i
			} else if word == word2 {
				index2 = i
			}
			if index1 >= 0 && index2 >= 0 {
				ans = min(ans, abs(index1-index2))
			}
		}
	}
	return ans
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
