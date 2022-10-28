package leetcode_1258_synonymous_sentences

import (
	"sort"
	"strings"
)

// 1258. 近义词句子
// https://leetcode.cn/problems/synonymous-sentences/

// generateSentences 并查集+dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func generateSentences(synonyms [][]string, text string) []string {
	var (
		ans    []string
		parent = make(map[string]string)
		rank   = make(map[string]int)
		dic    = make(map[string][]string)
		union  func(s1, s2 string)
		find   func(s string) string
	)
	find = func(s string) string {
		// 压缩路径
		if parent[s] != s {
			parent[s] = find(parent[s])
		}
		return parent[s]
	}
	union = func(s1, s2 string) {
		if fs1, fs2 := find(s1), find(s2); fs1 != fs2 {
			// 必须要按秩合并
			if rank[fs1] <= rank[fs2] {
				parent[fs1] = fs2
			} else {
				parent[fs2] = fs1
			}
			if rank[fs1] == rank[fs2] {
				rank[fs1]++
			}
		}
	}

	// 快速判断
	if len(text) == 0 || len(synonyms) == 0 {
		ans = append(ans, text)
		return ans
	}

	// 初始化
	for _, synonym := range synonyms {
		for _, s := range synonym {
			parent[s] = s
			rank[s] = 1
		}
	}
	// 并合并
	for _, synonym := range synonyms {
		if find(synonym[0]) != find(synonym[1]) {
			union(synonym[0], synonym[1])
		}
	}
	// 归类
	for k := range parent {
		p := find(k)
		dic[p] = append(dic[p], k)
	}

	// 待替换的单词数组
	textArray := strings.Split(text, " ")
	n := len(textArray)
	// dfs 从textArray的第idx个单词开始替换近义词，替换后的新句子放到ans中
	var dfs func(idx int, strs []string)
	dfs = func(idx int, strs []string) {
		if idx == n {
			ans = append(ans, strings.Join(strs, " "))
			return
		}
		// 若当前单词不在近义词字典中,则跳过当前单词
		if v, ok := parent[textArray[idx]]; !ok {
			strs = append(strs, textArray[idx])
			dfs(idx+1, strs)
		} else {
			// 获取当前待替换的单词的所有近义词
			bros := dic[v]
			sort.SliceStable(bros, func(i, j int) bool {
				return bros[i] < bros[j]
			})
			for _, bro := range bros {
				strs = append(strs, bro)
				dfs(idx+1, strs)
				strs = strs[:len(strs)-1]
			}
		}
	}
	dfs(0, nil)
	return ans
}
