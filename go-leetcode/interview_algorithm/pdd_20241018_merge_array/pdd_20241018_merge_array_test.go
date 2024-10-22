package pdd_20241018_merge_array

import (
	"container/heap"
	"fmt"
	"testing"
)

// 拼多多: 合并 k 个有序数组
// 例如: [][]int{{1,2,3},{2,3,4},{2,5,5}}

func TestMerge(t *testing.T) {
	var ls [][]LogSegment
	ls = append(ls, []LogSegment{{1, "x"}, {2, ""}, {3, ""}})
	ls = append(ls, []LogSegment{{2, ""}, {2, ""}, {5, ""}})
	ls = append(ls, []LogSegment{{2, ""}, {3, ""}, {4, ""}})
	ls = append(ls, []LogSegment{{1, ""}, {6, ""}})
	fmt.Println(ls)        // [[{1 } {2 } {3 }] [{2 } {2 } {5 }] [{2 } {3 } {4 }] [{1 } {6 }]]
	fmt.Println(merge(ls)) // [{1 } {1 } {2 } {2 } {2 } {2 } {3 } {3 } {4 } {5 } {6 }]
}

type LogSegment struct {
	timestamp int
	msg       string
}

type LogItem struct {
	i         int
	j         int
	timestamp int
}

type IHeap []LogItem

func (h *IHeap) Len() int { return len(*h) }

func (h *IHeap) Swap(i, j int) { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IHeap) Less(i, j int) bool { return (*h)[i].timestamp < (*h)[j].timestamp }

func (h *IHeap) Push(x interface{}) {
	*h = append(*h, x.(LogItem))
}

func (h *IHeap) Pop() interface{} {
	old := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return old
}

func merge(ls [][]LogSegment) []LogSegment {
	var ans []LogSegment
	m := len(ls)
	if m == 0 {
		return ans
	}

	h := IHeap{}
	heap.Init(&h)
	for i := 0; i < m; i++ {
		if len(ls[i]) > 0 {
			heap.Push(&h, LogItem{
				i:         i,
				j:         0,
				timestamp: ls[i][0].timestamp,
			})
		}
	}

	for h.Len() > 0 {
		item := heap.Pop(&h).(LogItem)
		ans = append(ans, ls[item.i][item.j])
		if item.j+1 < len(ls[item.i]) {
			heap.Push(&h, LogItem{
				i:         item.i,
				j:         item.j + 1,
				timestamp: ls[item.i][item.j+1].timestamp,
			})
		}
	}
	return ans
}
