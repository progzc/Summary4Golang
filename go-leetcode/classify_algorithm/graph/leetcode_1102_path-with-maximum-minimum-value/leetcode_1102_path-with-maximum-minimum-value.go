package leetcode_1102_path_with_maximum_minimum_value

import (
	"container/heap"
	"math"
)

// 1102. 得分最高的路径
// https://leetcode.cn/problems/path-with-maximum-minimum-value/

type Point struct {
	x   int
	y   int
	val int
}

type hp []*Point

func (h hp) Len() int {
	return len(h)
}

func (h hp) Less(i, j int) bool {
	return h[i].val > h[j].val
}

func (h hp) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *hp) Push(x interface{}) {
	*h = append(*h, x.(*Point))
}

func (h *hp) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}

// maximumMinimumPath 优先队列+迭代
// 时间复杂度: O(m*nlog(n))
// 空间复杂度: O(mn)
// 思路：
//	a.定义一个队列，将下一步可达的地方存入队列。
//	b.每次取出队列中最大的值作为下一步要探索的地方。
//	c.依次记录所有走过的位置中，最小的值作为结果。
func maximumMinimumPath(grid [][]int) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}

	var (
		m       = len(grid)
		n       = len(grid[0])
		dirs    = [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		visited = make([][]bool, m)
		queue   = hp{}
		score   = math.MaxInt32
	)
	for i := 0; i < m; i++ {
		visited[i] = make([]bool, n)
	}
	visited[0][0] = true
	queue = append(queue, &Point{x: 0, y: 0, val: grid[0][0]})
	heap.Init(&queue)

	for len(queue) > 0 {
		p := heap.Pop(&queue).(*Point)
		x, y, val := p.x, p.y, p.val
		score = min(score, val)
		if x == m-1 && y == n-1 {
			return score
		}
		for _, dir := range dirs {
			newX, newY := x+dir[0], y+dir[1]
			if newX >= 0 && newX < m && newY >= 0 && newY < n && !visited[newX][newY] {
				heap.Push(&queue, &Point{x: newX, y: newY, val: grid[newX][newY]})
				visited[newX][newY] = true
			}
		}
	}
	return score
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
