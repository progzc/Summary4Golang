package leetcode_1101_the_earliest_moment_when_everyone_become_friends

import "sort"

// 1101. 彼此熟识的最早时间
// https://leetcode.cn/problems/the-earliest-moment-when-everyone-become-friends/

// earliestAcq 排序+并查集
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
func earliestAcq(logs [][]int, n int) int {
	if len(logs) == 0 || len(logs[0]) == 0 {
		return -1
	}

	var (
		p     = make([]int, n)
		r     = make([]int, n)
		find  func(x int) int
		union func(x, y int)
		count = 0
	)
	for i := 0; i < n; i++ {
		p[i] = i
		r[i] = 1
	}

	find = func(x int) int {
		if p[x] != x {
			// 压缩路径
			p[x] = find(p[x])
		}
		return p[x]
	}
	union = func(x, y int) {
		if fx, fy := find(x), find(y); fx != fy {
			count++
			// 按秩合并
			if r[fx] <= r[fy] {
				p[fx] = fy
			} else {
				p[fy] = fx
			}
			if r[fx] == r[fy] {
				r[fx]++
			}
		}
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i][0] < logs[j][0]
	})
	for _, log := range logs {
		union(log[1], log[2])
		if count == n-1 {
			return log[0]
		}
	}
	return -1
}
