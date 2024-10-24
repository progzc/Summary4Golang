package leetcode_3175_find_the_first_player_to_win_k_games_in_a_row

import (
	"fmt"
	"testing"
)

// 3175. 找到连续赢 K 场比赛的第一位玩家
// https://leetcode.cn/problems/find-the-first-player-to-win-k-games-in-a-row

func TestFindWinningPlayer(t *testing.T) {
	skills, k := []int{7, 11}, 2
	fmt.Println(findWinningPlayer(skills, k)) // 1

	skills, k = []int{16, 4, 7, 17}, 562084119
	fmt.Println(findWinningPlayer(skills, k)) // 3

	skills, k = []int{4, 2, 6, 3, 9}, 2
	fmt.Println(findWinningPlayer(skills, k)) // 2
}

// findWinningPlayer 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func findWinningPlayer(skills []int, k int) int {
	n := len(skills)
	cnt := 0
	i, lastI := 0, 0
	for i < n {
		j := i + 1
		for j < n && skills[i] > skills[j] && cnt < k {
			j++
			cnt++
		}
		if cnt == k {
			return i
		}
		cnt = 1
		lastI = i
		i = j
	}
	return lastI
}

// findWinningPlayer_2 模拟（会超时）
func findWinningPlayer_2(skills []int, k int) int {
	n := len(skills)
	stack := make([]int, n)
	for i := 0; i < n; i++ {
		stack[i] = i
	}

	if k > n-1 {
		k = n - 1
	}
	count := 0
	for count < k {
		if skills[stack[0]] > skills[stack[1]] {
			count++
			x := stack[1]
			stack = append([]int{stack[0]}, stack[2:]...)
			stack = append(stack, x)
		} else {
			x := stack[0]
			stack = append(stack[1:], x)
			count = 1
		}
	}
	return stack[0]
}
