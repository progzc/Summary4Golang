package leetcode_0210_course_schedule_ii

// 210. 课程表 II
// https://leetcode.cn/problems/course-schedule-ii/

// 与下面题的解法类似:
// 207. 课程表
// https://leetcode.cn/problems/course-schedule/

// findOrder bfs（使用入度来判断有向图是否有环）
// 时间复杂度: O(n+m)
// 空间复杂度: O(n+m)
// 概念:
//	拓扑排序: 把一个 有向无环图 转成 线性的排序。
//	入度和出度: 如果存在一条有向边 A --> B，则这条边给 A 增加了 1 个出度，给 B 增加了 1 个入度。
// 思路：使用拓扑排序来解决判断有向图是否有环的问题。
//	具体步骤如下:
//	1.根据依赖关系，构建邻接表、入度数组。
//	2.选取入度为 0 的数据，根据邻接表，减小依赖它的数据的入度。
//	3.找出入度变为 0 的数据，重复第 2 步。
//	4.直至所有数据的入度为 0，得到排序，如果还有数据的入度不为 0，说明图中存在环。
func findOrder(numCourses int, prerequisites [][]int) []int {
	// 1.课号和对应的入度
	inDegree := make(map[int]int)
	// 2.将所有的课程先放入
	for i := 0; i < numCourses; i++ {
		inDegree[i] = 0
	}
	// 3.依赖关系(构建邻接表): 依赖当前课程的后续课程
	adj := make(map[int][]int)
	// 4.初始化入度和依赖关系
	for _, relate := range prerequisites {
		cur, next := relate[1], relate[0]
		inDegree[next] += 1
		adj[cur] = append(adj[cur], next)
	}
	// 5.将入度为0的课程放入队列，队列中的课程就是没有先修可以先学的课程
	var (
		queue []int
		ans   []int
	)
	for k, v := range inDegree {
		if v == 0 {
			queue = append(queue, k)
			ans = append(ans, k)
		}
	}
	// 6.遍历当前邻接表, 更新其入度; 更新之后查看入度, 如果为0, 加入到队列
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		// 遍历当前课程的邻接表, 更新后继节点的入度
		if _, ok := adj[cur]; !ok {
			continue
		}
		for _, k := range adj[cur] {
			inDegree[k] -= 1
			if inDegree[k] == 0 {
				queue = append(queue, k)
				ans = append(ans, k)
			}
		}
	}
	// 7.遍历入队, 如果还有课程的入度不为0, 返回false
	for _, v := range inDegree {
		if v != 0 {
			return []int{}
		}
	}
	return ans
}
