package leetcode_0320_generalized_abbreviation

import "strconv"

// 0320. 列举单词的全部缩写
// https://leetcode.cn/problems/generalized-abbreviation/

// generateAbbreviations dfs
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(n)
func generateAbbreviations(word string) []string {
	var (
		n = len(word)
		// i为当前遍历到字符串的第几个字符
		// k为连续缩写的字符个数
		dfs func(i, k int, path string)
		ans []string
	)

	dfs = func(i, k int, path string) {
		if i == n {
			if k > 0 {
				// 若计数大于0,则转为字符串,加到末尾
				path += strconv.Itoa(k)
			}
			ans = append(ans, path)
			return
		}

		// a.每个字符有两种选择,一种是缩写
		dfs(i+1, k+1, path)
		// b.一种是不缩写。不缩写的时候，要把之前的k转换为string加到path后面，然后再添加字符本身
		if k > 0 {
			path += strconv.Itoa(k)
		}
		path += string(word[i])
		dfs(i+1, 0, path)
	}
	dfs(0, 0, "")
	return ans
}
