package leetcode_0127_word_ladder

import (
	"math"
	"testing"
)

// 0127.单词接龙
// https://leetcode-cn.com/problems/word-ladder/

func TestLadderLength(t *testing.T) {
	beginWord, endWord := "hit", "cog"
	wordList := []string{"hot", "dot", "dog", "lot", "log", "cog"}
	println(ladderLength(beginWord, endWord, wordList))
}

// ladderLength 图+广度优先搜索+虚拟节点
// 时间复杂度：O(N*C^2)
// 空间复杂度：O(N*C^2)
// 思路：最短转换序列的长度，看到最短首先想到的就是广度优先搜索。
//		想到广度优先搜索自然而然的就能想到图。
func ladderLength(beginWord string, endWord string, wordList []string) int {
	wordId := map[string]int{}
	graph := [][]int{}
	// 编号
	addWord := func(word string) int {
		id, has := wordId[word]
		if !has {
			id = len(wordId)
			wordId[word] = id
			graph = append(graph, []int{})
		}
		return id
	}
	// 构造图
	addEdge := func(word string) int {
		id1 := addWord(word)
		s := []byte(word)
		for i, b := range s {
			s[i] = '*'
			id2 := addWord(string(s))
			graph[id1] = append(graph[id1], id2)
			graph[id2] = append(graph[id2], id1)
			s[i] = b
		}
		return id1
	}

	// 初始化图
	for _, word := range wordList {
		addEdge(word)
	}
	beginId := addEdge(beginWord)
	endId, has := wordId[endWord]
	if !has {
		return 0
	}

	// 避免重复计算
	const inf = math.MaxInt64
	dist := make([]int, len(wordId))
	for i := range dist {
		dist[i] = inf
	}

	// 通过栈来广度搜索图
	dist[beginId] = 0
	queue := []int{beginId}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		if v == endId {
			return dist[endId]/2 + 1
		}
		for _, w := range graph[v] {
			// 这个if很重要，避免重复计算
			if dist[w] == inf {
				dist[w] = dist[v] + 1
				queue = append(queue, w)
			}
		}
	}
	return 0
}
