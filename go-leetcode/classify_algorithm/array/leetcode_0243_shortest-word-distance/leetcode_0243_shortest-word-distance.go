package leetcode_0243_shortest_word_distance

import "math"

// 0243.最短单词距离
// https://leetcode.cn/problems/shortest-word-distance/
// 注意：word1 != word2

// shortestDistance_2 一次遍历
// 时间复杂度：O(n)
// 空间复杂度：O(1)
func shortestDistance_2(wordsDict []string, word1 string, word2 string) int {
	ans := len(wordsDict)
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
	return ans
}

// shortestDistance 一般解法
// 时间复杂度：O(n^2)
// 空间复杂度：O(n)
func shortestDistance(wordsDict []string, word1 string, word2 string) int {
	m := map[string][]int{}
	for i, word := range wordsDict {
		if v, ok := m[word]; !ok {
			m[word] = []int{i}
		} else {
			v = append(v, i)
			m[word] = v
		}
	}

	ans := math.MaxInt32
	x1s, x2s := m[word1], m[word2]
	for _, x1 := range x1s {
		for _, x2 := range x2s {
			ans = min(abs(x1-x2), ans)
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
	if x < y {
		return x
	}
	return y
}
