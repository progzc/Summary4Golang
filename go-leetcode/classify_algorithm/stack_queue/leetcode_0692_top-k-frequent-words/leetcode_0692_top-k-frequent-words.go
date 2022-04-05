package leetcode_0692_top_k_frequent_words

import (
	"container/heap"
	"sort"
)

// 0692.前K个高频单词
// https://leetcode-cn.com/problems/top-k-frequent-words/

// topKFrequent 哈希表+排序
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
func topKFrequent(words []string, k int) []string {
	cnt := map[string]int{}
	for _, word := range words {
		cnt[word]++
	}

	uw := make([]string, 0, len(cnt))
	for w := range cnt {
		uw = append(uw, w)
	}

	sort.SliceStable(uw, func(i, j int) bool {
		s1, s2 := uw[i], uw[j]
		return cnt[s1] > cnt[s2] || (cnt[s1] == cnt[s2] && s1 < s2)
	})
	return uw[:k]
}

// topKFrequent_2 优先队列
// 时间复杂度: O(nlog(k))
// 空间复杂度: O(n)
func topKFrequent_2(words []string, k int) []string {
	cnt := map[string]int{}
	for _, w := range words {
		cnt[w]++
	}
	h := &hp{}
	for w, c := range cnt {
		heap.Push(h, Item{w, c})
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	ans := make([]string, k)
	for i := k - 1; i >= 0; i-- {
		ans[i] = heap.Pop(h).(Item).word
	}
	return ans
}

type Item struct {
	word  string
	count int
}

type hp []Item

func (h hp) Len() int {
	return len(h)
}

func (h hp) Less(i, j int) bool {
	a, b := h[i], h[j]
	// 小顶堆
	return a.count < b.count || (a.count == b.count && a.word > b.word)
}

func (h hp) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *hp) Push(x interface{}) {
	*h = append(*h, x.(Item))
}

func (h *hp) Pop() interface{} {
	a := *h
	v := a[len(a)-1]
	*h = a[:len(a)-1]
	return v
}
