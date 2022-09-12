package leetcode_0582_kill_process

// 0582. 杀掉进程
// https://leetcode.cn/problems/kill-process/

// killProcess 构造树
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func killProcess(pid []int, ppid []int, kill int) []int {
	type Node struct {
		Val      int
		Children []*Node
	}

	var (
		m              = map[int]*Node{}
		getAllChildren func(parent *Node)
		ans            []int
	)

	for _, v := range pid {
		m[v] = &Node{Val: v}
	}
	for i, v := range ppid {
		if v > 0 {
			if node, ok := m[v]; ok {
				node.Children = append(node.Children, m[pid[i]])
			}
		}
	}

	getAllChildren = func(parent *Node) {
		for _, child := range parent.Children {
			ans = append(ans, child.Val)
			getAllChildren(child)
		}
	}

	if kill == 0 {
		return ans
	}
	ans = append(ans, kill)
	getAllChildren(m[kill])
	return ans
}

// killProcess_2 哈希表+dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func killProcess_2(pid []int, ppid []int, kill int) []int {
	var (
		m              = map[int][]int{}
		getAllChildren func(child []int)
		ans            []int
	)
	for i, v := range ppid {
		if v > 0 {
			m[v] = append(m[v], pid[i])
		}
	}

	getAllChildren = func(child []int) {
		for _, v := range child {
			ans = append(ans, v)
			getAllChildren(m[v])
		}
	}
	if kill == 0 {
		return ans
	}
	ans = append(ans, kill)
	getAllChildren(m[kill])
	return ans
}

// killProcess_3 哈希表+bfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func killProcess_3(pid []int, ppid []int, kill int) []int {
	var (
		m     = map[int][]int{}
		stack []int
		ans   []int
	)
	for i, v := range ppid {
		if v > 0 {
			m[v] = append(m[v], pid[i])
		}
	}

	if kill == 0 {
		return ans
	}
	ans = append(ans, kill)
	stack = append(stack, m[kill]...)
	for len(stack) > 0 {
		t := stack[0]
		stack = stack[1:]
		ans = append(ans, t)
		stack = append(stack, m[t]...)
	}
	return ans
}
