package leetcode_0076_minimum_window_substring

import "math"

// 0076.最小覆盖子串
// https://leetcode-cn.com/problems/minimum-window-substring/

// minWindow 滑动窗口（双指针+哈希）
// 时间复杂度：O(n)
// 空间复杂度: O(k) k为字符集合
// 思路：
//	need: 始终记录着当前滑动窗口下，我们还需要的元素数量，我们在改变i,j时，需同步维护need。
//	needCnt: 如果每次判断滑动窗口是否包含了T的所有元素，都去遍历need看是否所有元素数量都小于等于0，这个会耗费O(k)的时间复杂度;
//			 可以维护一个额外的变量needCnt来记录所需元素的总数量，当我们碰到一个所需元素c，不仅need[c]的数量减少1，同时needCnt也要减少1，
//			 这样我们通过needCnt就可以知道是否满足条件，而无需遍历字典了。
func minWindow(s string, t string) string {
	var (
		need    = map[byte]int{}
		needCnt = len(t)
	)
	for i := 0; i < len(t); i++ {
		need[t[i]] += 1
	}
	// 记录起始位置
	l, r := 0, 0
	start, end, size := 0, 0, math.MaxInt32
	for r < len(s) {
		// 若需要字符s[r]
		if need[s[r]] > 0 {
			needCnt--
		}
		need[s[r]]--
		// 若窗口中已经包含所有元素
		if needCnt == 0 {
			for l < r && need[s[l]] < 0 {
				// 是否左边移出窗口的字符
				need[s[l]]++
				l++
			}
			if r-l+1 < size {
				size = r - l + 1
				// 记录最小值时候的开始位置
				start = l
				end = r + 1
			}
			// 这个时候s[l]一定是必须的元素了，这里需要释放，重新寻找下一个最短字符串
			// 左边界右移之前需要释放need[s[l]]
			need[s[l]]++
			l++
			needCnt++
		}
		r++
	}
	return s[start:end]
}
