package leetcode_0294_flip_game_ii

// 0294. 翻转游戏 II
// https://leetcode.cn/problems/flip-game-ii/

// ACM-ICPC: https://edwiv.com/wp-content/uploads/2019/08/ACM-ICPC_Templates_201805.pdf
// 本质是博弈论的Sprague-Grundy定理: https://zhuanlan.zhihu.com/p/20611132

// canWin dfs
// 时间复杂度: O(n^(n/2))
// 空间复杂度: O(n)
func canWin(currentState string) bool {
	n := len(currentState)
	ss := []byte(currentState)
	// 遍历字符串, 找到两个相邻的'+', 若找不到则必败, 返回false
	for i := 0; i < n-1; i++ {
		// 找到两个相邻的'+', 翻转他们, 如果对手对于翻转后的字符串找不到相邻的'+',
		// 说明你赢了, 有必胜的方案
		if ss[i] == '+' && ss[i+1] == '+' {
			ss[i] = '-'
			ss[i+1] = '-'
			if canWin(string(ss)) == false {
				return true
			}
			// 回溯, 尝试下个方案
			ss[i] = '+'
			ss[i+1] = '+'
		}
	}
	return false
}
