package leetcode_0293_flip_game

// 0293. 翻转游戏
// https://leetcode.cn/problems/flip-game/

// generatePossibleNextMoves dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func generatePossibleNextMoves(currentState string) []string {
	var ans []string
	ss := []byte(currentState)
	for i := 0; i < len(ss)-1; i++ {
		if ss[i] == '+' && ss[i+1] == '+' {
			ss[i] = '-'
			ss[i+1] = '-'
			ans = append(ans, string(ss))
			ss[i] = '+'
			ss[i+1] = '+'
		}
	}
	return ans
}
