package go_leetcode

import (
	"testing"
)

func Test_sword2offer_0081_combinationSum(t *testing.T) {
	type params struct {
		candidates []int
		target     int
	}
	tests := []struct {
		p    params
		want [][]int
	}{
		{
			p: params{
				candidates: []int{2, 3, 6, 7},
				target:     7,
			},
			want: [][]int{{7}, {2, 2, 3}},
		},
		{
			p: params{
				candidates: []int{2, 3, 5},
				target:     8,
			},
			want: [][]int{{2, 2, 2, 2}, {2, 3, 3}, {3, 5}},
		},
		{
			p: params{
				candidates: []int{2},
				target:     1,
			},
			want: [][]int{},
		},
		{
			p: params{
				candidates: []int{1},
				target:     1,
			},
			want: [][]int{{1}},
		},
		{
			p: params{
				candidates: []int{1},
				target:     2,
			},
			want: [][]int{{1, 1}},
		},
	}
	for _, test := range tests {
		fact := sword2offer_0081_combinationSum(test.p.candidates, test.p.target)
		t.Logf("params=%v,want=%v,fact=%v", test.p, test.want, fact)
	}
}

// sword2offer_0081_combinationSum 搜索回溯（不加剪枝）
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(target)
// 解题思路：根据一个数选和不选画出树形图
func sword2offer_0081_combinationSum(candidates []int, target int) [][]int {
	var comb []int
	var ans [][]int
	var dfs func(target, idx int)
	dfs = func(target, idx int) {
		if idx == len(candidates) {
			return
		}
		if target == 0 {
			ans = append(ans, append([]int(nil), comb...))
			return
		}
		// 直接跳过
		dfs(target, idx+1)
		// 直接选择当前数
		if target-candidates[idx] >= 0 {
			comb = append(comb, candidates[idx])
			dfs(target-candidates[idx], idx)
			comb = comb[:len(comb)-1]
		}
	}
	dfs(target, 0)
	return ans
}
