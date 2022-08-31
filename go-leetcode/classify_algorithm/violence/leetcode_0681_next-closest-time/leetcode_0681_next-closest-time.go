package leetcode_0681_next_closest_time

import (
	"fmt"
	"strconv"
)

// 0681. 最近时刻
// https://leetcode.cn/problems/next-closest-time/

// nextClosestTime 模拟时钟前进
// 时间复杂度: O(24*60)
// 空间复杂度: O(1)
func nextClosestTime(time string) string {
	hour, _ := strconv.ParseInt(time[:2], 10, 64)
	minute, _ := strconv.ParseInt(time[3:], 10, 64)
	t := hour*60 + minute

	// 使用mask来判断出现的时间是否是由已经出现的数字组成的
	mask := ^(1<<(hour/10) | 1<<(hour%10) | 1<<(minute/10) | 1<<(minute%10))

	for i := 0; i < 24*60; i++ {
		t += 1
		h, m := t/60, t%60
		h = h % 24

		// 判断出现的时间是否合理
		if (1<<(h/10)|1<<(h%10)|1<<(m/10)|1<<(m%10))&mask == 0 {
			return fmt.Sprintf("%02d:%02d", h, m)
		}
	}
	return ""
}

// nextClosestTime_2 从允许的数字生成(对方法一的进一步优化)
// 时间复杂度: O(4^4)
// 空间复杂度: O(1)
func nextClosestTime_2(time string) string {
	hour, _ := strconv.ParseInt(time[:2], 10, 64)
	minute, _ := strconv.ParseInt(time[3:], 10, 64)
	t := hour*60 + minute

	// 枚举的数字集合
	ans, r1, r2 := int64(24*60), hour, minute
	cands := []int64{hour / 10, hour % 10, minute / 10, minute % 10}
	for _, i := range cands {
		for _, j := range cands {
			for _, x := range cands {
				for _, y := range cands {
					h := 10*i + j
					m := 10*x + y
					if h < 24 && m < 60 {
						if nt := h*60 + m; nt != t {
							dt := (nt - t) % (24 * 60)
							// 注意golang的负数取余
							if dt < 0 {
								dt += 24 * 60
							}
							if dt < ans {
								ans = dt
								r1, r2 = h, m
							}
						}
					}
				}
			}
		}
	}
	return fmt.Sprintf("%02d:%02d", r1, r2)
}
