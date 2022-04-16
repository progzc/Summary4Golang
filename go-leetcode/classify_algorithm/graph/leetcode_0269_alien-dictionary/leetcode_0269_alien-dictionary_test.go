package leetcode_0269_alien_dictionary

import (
	"fmt"
	"strings"
	"testing"
)

// 0269.火星词典
// https://leetcode-cn.com/problems/alien-dictionary/

func TestAlienOrder(t *testing.T) {
	words := []string{"wr", "wr"}
	fmt.Println(alienOrder(words))
}

// alienOrder 有向图
// 思路：构造一个有向图，进行广度优先搜索
func alienOrder(words []string) string {
	// 1.构建字符集
	unknown := map[byte]struct{}{}
	for _, word := range words {
		for i := 0; i < len(word); i++ {
			unknown[word[i]] = struct{}{}
		}
	}
	// 2.建立有向图
	graph := map[byte][]byte{}
	n := len(words)
	for i := 0; i < n-1; i++ {
		w1, w2 := words[i], words[i+1]
		maxLen := max(len(w1), len(w2))
		for j := 0; j < maxLen; j++ {
			if j == len(w2) {
				return ""
			}
			if j == len(w1) {
				break
			}
			if w1[j] != w2[j] {
				graph[w1[j]] = append(graph[w1[j]], w2[j])
				break
			}
		}
	}
	// 3.拓扑排序
	visited := [26]bool{} // 记录结点访问
	onPath := [26]bool{}  // 记录路线访问
	hasCycle := false     // 存在环标志
	path := strings.Builder{}
	var traverse func(key byte)

	traverse = func(key byte) {
		// 若路径曾访问过，表示存在环
		if onPath[key-'a'] {
			hasCycle = true
		}
		// 如果存在环，或者已经访问过该节点，则直接返回
		if hasCycle || visited[key-'a'] {
			return
		}
		visited[key-'a'] = true
		onPath[key-'a'] = true
		for _, v := range graph[key] {
			traverse(v)
		}
		path.WriteByte(key)
		// 这里这里必须回撤
		onPath[key-'a'] = false
	}

	for k := range graph {
		traverse(k)
	}
	// 4.输出结果
	if hasCycle {
		return ""
	}
	pathStr := path.String()
	for unk := range unknown {
		if strings.IndexByte(pathStr, unk) == -1 {
			path.WriteByte(unk)
		}
	}
	return reverse(path.String())
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
