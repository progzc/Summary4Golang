package leetcode_1136_parallel_courses

// 1136. 并行课程
// https://leetcode.cn/problems/parallel-courses/

// minimumSemesters 拓扑排序
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：我们知道只有当一门课程没有前置课程的时候，才能上这门课。我们可以把这门课程需要上的前置课程的数量表示为入度。
//		比如上图中 0 的入度为 0，没有前置课程需要上，而 3 的入度为 2，有 1 和 2 需要上。
//		所以我们可以每一次遍历把所有入度为 0 的课程上完，并且把他们的后置课程的入度减 1，重复这个过程直到没有课程可以上。
func minimumSemesters(n int, relations [][]int) int {
	// 所有课程的入度
	pre := map[int]int{}
	// 后置课程
	next := make([][]int, n)
	// 初始化：所有课程的入度为0
	for i := 0; i < n; i++ {
		pre[i] = 0
	}
	for _, r := range relations {
		// 计算初始入度
		pre[r[1]-1]++
		// 统计后置课程
		next[r[0]-1] = append(next[r[0]-1], r[1]-1)
	}

	term := 0
	for len(pre) > 0 {
		term++
		var learn []int
		for class, count := range pre {
			if count == 0 {
				// 入度为0,可以学习
				learn = append(learn, class)
			}
		}

		if len(learn) == 0 {
			// 如果没有课程可以学习了,说明有循环
			return -1
		}

		// 根据学过的课程更新pre
		for _, l := range learn {
			delete(pre, l)
			for _, class := range next[l] {
				pre[class]--
			}
		}
	}

	return term
}
