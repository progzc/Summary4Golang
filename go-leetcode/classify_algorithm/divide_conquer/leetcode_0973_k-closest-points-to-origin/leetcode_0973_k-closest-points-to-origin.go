package leetcode_0973_k_closest_points_to_origin

import (
	"container/heap"
	"sort"
)

// 0973.最接近原点的 K 个点
// https://leetcode-cn.com/problems/k-closest-points-to-origin/

// kClosest 排序
// 时间复杂度: O(n*log(n))
// 空间复杂度: O(log(n))
func kClosest(points [][]int, k int) [][]int {
	sort.SliceStable(points, func(i, j int) bool {
		p, q := points[i], points[j]
		return p[0]*p[0]+p[1]*p[1] < q[0]*q[0]+q[1]*q[1]
	})
	return points[:k]
}

// kClosest_2 大顶堆
// 时间复杂度: O(n*log(k))
// 空间复杂度: O(k))
func kClosest_2(points [][]int, k int) [][]int {
	h := make(hp, k)
	for i, p := range points[:k] {
		h[i] = pair{p[0]*p[0] + p[1]*p[1], p}
	}
	// 初始化堆
	heap.Init(&h)
	for _, p := range points[k:] {
		if dist := p[0]*p[0] + p[1]*p[1]; dist < h[0].dist {
			h[0] = pair{dist, p}
			// 效率比 pop 后 push 要快
			heap.Fix(&h, 0)
		}
	}
	var ans [][]int
	for _, p := range h {
		ans = append(ans, p.point)
	}
	return ans
}

type pair struct {
	dist  int
	point []int
}

type hp []pair

func (h hp) Len() int { return len(h) }

func (h hp) Less(i, j int) bool { return h[i].dist > h[j].dist } // 大顶堆

func (h hp) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *hp) Push(x interface{}) { *h = append(*h, x.(pair)) }

func (h *hp) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}
