package leetcode_0079_word_search

// 0079.单词搜索
// https://leetcode-cn.com/problems/word-search/

// exist 回溯
// 时间复杂度：O(M*N*3^L)
// 空间复杂度: O(min(M*N,L))
func exist(board [][]byte, word string) bool {
	var (
		m, n  = len(board), len(board[0])
		dirs  = [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		check func(i, j, k int) bool
		visit = make([][]bool, m)
	)
	// 初始化
	for i := 0; i < m; i++ {
		visit[i] = make([]bool, n)
	}
	// check 表示判断以网格的 (i,j)位置出发，能否搜索到单词word[k:]
	check = func(i, j, k int) bool {
		if i < 0 || i >= m || j < 0 || j >= n || visit[i][j] {
			return false
		}
		if board[i][j] != word[k] {
			return false
		}
		if k == len(word)-1 {
			return true
		}
		visit[i][j] = true
		defer func() {
			visit[i][j] = false
		}()
		for _, dir := range dirs {
			if check(dir[0]+i, dir[1]+j, k+1) {
				return true
			}
		}
		return false
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if check(i, j, 0) {
				return true
			}
		}
	}
	return false
}

// exist_2 回溯（优化）
// 时间复杂度：O(M*N*3^L)
// 空间复杂度: O(L)
// 思路：注意到board和word仅由大小写英文字母组成
func exist_2(board [][]byte, word string) bool {
	var (
		m, n  = len(board), len(board[0])
		dirs  = [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		check func(i, j, k int) bool
	)
	// check 表示判断以网格的 (i,j)位置出发，能否搜索到单词word[k:]
	check = func(i, j, k int) bool {
		if i < 0 || i >= m || j < 0 || j >= n || board[i][j] == '#' {
			return false
		}
		if board[i][j] != word[k] {
			return false
		}
		if k == len(word)-1 {
			return true
		}
		temp := board[i][j]
		board[i][j] = '#'
		defer func() {
			board[i][j] = temp
		}()
		for _, dir := range dirs {
			if check(dir[0]+i, dir[1]+j, k+1) {
				return true
			}
		}
		return false
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if check(i, j, 0) {
				return true
			}
		}
	}
	return false
}

func exist_wrong(board [][]byte, word string) bool {
	var (
		m, n  = len(board), len(board[0])
		dirs  = [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		check func(i, j, k int) bool
		visit = make([][]bool, m)
	)
	// 初始化
	for i := 0; i < m; i++ {
		visit[i] = make([]bool, n)
	}
	// check 表示判断以网格的 (i,j)位置出发，能否搜索到单词word[k:]
	check = func(i, j, k int) bool {
		if i < 0 || i >= m || j < 0 || j >= n || visit[i][j] {
			return false
		}
		// 针对输入 [["a"]] "a"，会报错
		// 易错点：
		//	a.下面两段if的顺序必须颠倒
		//	b.下面这句应该是k==len(word)-1，为什么?因为k是从小慢慢变大的,如果单词不存在board中，
		//	则必然在索引越界之前就会出现board[i][j] != word[k]；如果单词存在board中，则在搜完之前
		//	索引都不会越界
		if k == len(word) {
			return true
		}
		if board[i][j] != word[k] {
			return false
		}
		// 易错点：
		//	c.注意下面这行以及defer的位置
		visit[i][j] = true
		defer func() {
			visit[i][j] = false
		}()
		for _, dir := range dirs {
			if check(dir[0]+i, dir[1]+j, k+1) {
				return true
			}
		}
		return false
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if check(i, j, 0) {
				return true
			}
		}
	}
	return false
}
