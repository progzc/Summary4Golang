package leetcode_0727_minimum_window_subsequence

import "math"

// 0727. 最小窗口子序列
// https://leetcode.cn/problems/minimum-window-subsequence/

// minWindow_2 动态规划(前缀递推)
// 思路：先计算包含 s2 的前缀子串的窗口，再根据前缀子串窗口不断拓展，找到包含整个字符串的窗口。
func minWindow_2(s1 string, s2 string) string {
	return ""
}

// minWindow 双指针(推荐)
// 思路：严格说这题并不是通俗的滑动窗口思路，因为这题对顺序有要求，所以单纯用哈希表统计字符个数的方法是行不通的，
// 不过依然还是使用双指针的思路：先正向遍历找到包含T的字串，正向遍历的终点就是一个可行的右边界，
// 然后再从右边界反向遍历回来寻找离右边界最近的左边界，这样就得到一个可行的序列，记录长度，
// 又开始从左边界的下一位开始正向遍历...
func minWindow(s1 string, s2 string) string {
	m, n, minLen := len(s1), len(s2), math.MaxInt32
	i, j := 0, 0
	var ans string
	for i < m {
		if s1[i] == s2[j] {
			if j == n-1 { // 找到了一个包含s2的子串（右边界为i）
				right := i // 此时i指向右边界
				// 再顺着右边界反向遍历，直到找到（离右边界最近的）左边界
				for j >= 0 {
					if s1[i] == s2[j] {
						i--
						j--
					} else {
						i--
					}
				}
				// 因为上面的反向遍历结束时i指向了左边界的前一位，i++之后此时i就指向左边界
				i++
				if minLen > right-i+1 {
					minLen = right - i + 1
					ans = s1[i : i+minLen]
				}
			}
			j++
		}
		i++
	}
	return ans
}
