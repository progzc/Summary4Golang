package leetcode_3216_lexicographically_smallest_string_after_a_swap

// 3216. 交换后字典序最小的字符串
// https://leetcode.cn/problems/lexicographically-smallest-string-after-a-swap
// 注意：最多交换一次

// getSmallestString 顺序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func getSmallestString(s string) string {
	ans := s
	for i := 0; i < len(s)-1; i++ {
		b := []byte(s)
		x, y := int(s[i]), int(s[i+1])
		if x%2 == y%2 && x > y {
			b[i], b[i+1] = b[i+1], b[i]
			ans = string(b)
			break
		}
	}
	return ans
}

// getSmallestString_2 顺序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func getSmallestString_2(s string) string {
	r := []rune(s)
	for i := 0; i+1 < len(r); i++ {
		if r[i] > r[i+1] && r[i]%2 == r[i+1]%2 {
			r[i], r[i+1] = r[i+1], r[i]
			break
		}
	}
	return string(r)
}
