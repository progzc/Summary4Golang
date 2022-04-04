package leetcode_0049_group_anagrams

import "sort"

// 0049.字母异位词分组
// https://leetcode-cn.com/problems/group-anagrams/

// groupAnagrams 排序+哈希表
// 时间复杂度: O(n*k*log(k))
// 空间复杂度: O(nk)
func groupAnagrams(strs []string) [][]string {
	m := map[string][]string{}
	for _, str := range strs {
		// 先将str转为[]byte（若有中文字符则转化为[]rune）,再排序
		s := []byte(str)
		sort.SliceStable(s, func(i, j int) bool {
			return s[i] < s[j]
		})
		m[string(s)] = append(m[string(s)], str)
	}

	ans := make([][]string, 0, len(m))
	for _, v := range m {
		ans = append(ans, v)
	}
	return ans
}

// groupAnagrams_2 哈希表+计数
// 时间复杂度: O(n*k)
// 空间复杂度: O(n*k)
// 思路：使用[26]int作为map的key，比较巧妙
func groupAnagrams_2(strs []string) [][]string {
	m := map[[26]int][]string{}
	for _, str := range strs {
		cnt := [26]int{}
		for _, b := range str {
			cnt[b-'a']++
		}
		m[cnt] = append(m[cnt], str)
	}
	ans := make([][]string, 0, len(m))
	for _, v := range m {
		ans = append(ans, v)
	}
	return ans
}
