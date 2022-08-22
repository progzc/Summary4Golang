package leetcode_0244_shortest_word_distance_ii

// 0244. 最短单词距离 II
// https://leetcode.cn/problems/shortest-word-distance-ii/
// 注意：
//	a.word1 != word2
//	b.连续操作

// WordDistance 哈希表+双指针
// 时间复杂度：O(n)
// 空间复杂度：O(n)
type WordDistance struct {
	m map[string][]int
	l int
}

func Constructor(wordsDict []string) WordDistance {
	m, l := map[string][]int{}, 0
	for i, word := range wordsDict {
		m[word] = append(m[word], i)
		l++
	}
	return WordDistance{m: m, l: l}
}

func (this *WordDistance) Shortest(word1 string, word2 string) int {
	ans := this.l
	list1, list2 := this.m[word1], this.m[word2]
	i, j := 0, 0
	m, n := len(list1), len(list2)
	for i < m && j < n {
		index1, index2 := list1[i], list2[j]
		ans = min(ans, abs(index1-index2))
		if index1 < index2 {
			i++
		} else {
			j++
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
