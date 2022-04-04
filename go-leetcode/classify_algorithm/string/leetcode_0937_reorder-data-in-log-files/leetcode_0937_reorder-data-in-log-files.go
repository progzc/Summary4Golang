package leetcode_0937_reorder_data_in_log_files

import (
	"sort"
	"strings"
)

// 0937.重新排列日志文件
// https://leetcode-cn.com/problems/reorder-data-in-log-files/

// reorderLogFiles 自定义排序
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
func reorderLogFiles(logs []string) []string {
	var nums []string
	var strs []string
	for _, log := range logs {
		ss := strings.Split(log, " ")
		if ss[1][0] >= '0' && ss[1][0] <= '9' {
			nums = append(nums, log)
		} else {
			strs = append(strs, log)
		}
	}

	sort.Slice(strs, func(i, j int) bool {
		x := strings.Index(strs[i], " ")
		y := strings.Index(strs[j], " ")
		sx := strs[i][x+1:]
		sy := strs[j][y+1:]

		if sx != sy {
			return sx < sy
		} else {
			return strs[i][:x+1] < strs[j][:y+1]
		}
	})
	strs = append(strs, nums...)
	return strs
}

// reorderLogFiles_2 自定义排序（优化）
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
func reorderLogFiles_2(logs []string) []string {
	// 注意事项：下面这句换成sort.Slice则会出错
	// sort.SliceStable：稳定排序
	// sort.Slice：不稳定排序
	sort.SliceStable(logs, func(i, j int) bool {
		x := strings.SplitN(logs[i], " ", 2)
		y := strings.SplitN(logs[j], " ", 2)

		isDigitX := x[1][0] >= '0' && x[1][0] <= '9'
		isDigitY := y[1][0] >= '0' && y[1][0] <= '9'

		// 1.字母日志按字母数字顺序排列，先按内容排序，再按标识符排序
		if !isDigitX && !isDigitY {
			if x[1] != y[1] {
				// 1.1先按内容排序
				return x[1] < y[1]
			} else {
				// 1.2再按标识符排序
				return x[0] < y[0]
			}
		} else if isDigitX && isDigitY {
			// 2.数字日志的顺序保持不变
			return i < j
		} else {
			// 3.有字母日志及数字日志，则数字日志排在后面
			if isDigitX {
				return false
			} else {
				return true
			}
		}
	})
	return logs
}
