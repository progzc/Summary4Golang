package leetcode_1152_analyze_user_website_visit_pattern

import "sort"

// 1152. 用户网站访问行为分析
// https://leetcode.cn/problems/analyze-user-website-visit-pattern/

type Node struct {
	username  string
	timestamp int
	website   string
}

// mostVisitedPattern 工程题
// 时间复杂度: O(n^3)
// 空间复杂度: O(n^3)
// 注意事项：
//	a.至少按某种次序访问过一次，这句话告诉我们需要根据时间排序。
//	b.用户可能不是连续访问这三个路径的。假设一个用户在 1,2,3,4 这 4 个时间点访问了 a,b,c,d 这 4 个网站，那么 a,c,d 也是一个合法的路径。
//	c.题目要求的是最多的用户访问的路径，所以一个用户访问很多次，也只能算一次（这点要特别注意）。
// 思路：
//	a.首先我们需要将 username、timestamp 和 website 这 3 个东西绑定起来。最直观的办法就是使用结构体。将 3 个数组使用结构体数组关联起来。
//		struct Node { username, timestamp, website }
//	b.对结构体数组按照 timestamp 排序，保证每个用户的访问次序。
//	c.使用哈希表存储每一个用户的访问的网站，哈希表的键是用户名 name，值是一个字符串数组。数组的值就是对应的 website。因为第二步已经排过序了，所以数组是有序的，可以直接使用
//	d.三重遍历每个用户的 website，获得所有的访问路径。再一次用哈希表存储所有的访问路径，对应的值就是用户的数量。
//	e.最后通过遍历哈希表获得最多用户访问且字典序排列最小的那个值。
func mostVisitedPattern(username []string, timestamp []int, website []string) []string {
	// a.首先我们需要将 username、timestamp 和 website 这 3 个东西绑定起来。最直观的办法就是使用结构体。将 3 个数组使用结构体数组关联起来。
	n := len(username)
	nodes := make([]Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = Node{username[i], timestamp[i], website[i]}
	}

	// b.对结构体数组按照 timestamp 排序，保证每个用户的访问次序。
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].timestamp < nodes[j].timestamp
	})

	// c.使用哈希表存储每一个用户的访问的网站，哈希表的键是用户名 name，值是一个字符串数组。数组的值就是对应的 website。因为第二步已经排过序了，所以数组是有序的，可以直接使用
	m := make(map[string][]Node)
	for i := 0; i < n; i++ {
		m[nodes[i].username] = append(m[nodes[i].username], nodes[i])
	}

	// d.三重遍历每个用户的 website，获得所有的访问路径。再一次用哈希表存储所有的访问路径，对应的值就是用户的数量。
	route := make(map[[3]string]int)
	for _, v := range m {
		tmp := make(map[[3]string]int)
		for i := 0; i < len(v); i++ {
			for j := i + 1; j < len(v); j++ {
				for k := j + 1; k < len(v); k++ {
					tmp[[3]string{v[i].website, v[j].website, v[k].website}] = 1
				}
			}
		}

		for k1, v1 := range tmp {
			route[k1] += v1
		}
	}

	// e.最后通过遍历哈希表获得最多用户访问且字典序排列最小的那个值。
	max := -1
	ans := [3]string{}
	for k, v := range route {
		if v > max {
			ans = k
			max = v
		} else if v == max {
			if k[0] < ans[0] || (k[0] == ans[0] && k[1] < ans[1]) ||
				(k[0] == ans[0] && k[1] == ans[1] && k[2] < ans[2]) {
				ans = k
			}
		}
	}
	return []string{ans[0], ans[1], ans[2]}
}
