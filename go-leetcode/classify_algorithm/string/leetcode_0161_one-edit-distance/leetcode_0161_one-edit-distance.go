package leetcode_0161_one_edit_distance

// 0161.相隔为 1 的编辑距离
// https://leetcode.cn/problems/one-edit-distance/

// isOneEditDistance 指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func isOneEditDistance(s string, t string) bool {
	ns, nt := len(s), len(t)

	// 假设s总是比t短或者长度长度，否则交换比较
	// 把考虑 增加/删除 变成了只需要考虑 增加 的情况
	if ns > nt {
		return isOneEditDistance(t, s)
	}

	// 1.如果两个字符串的长度差大于1，则编辑距离不为1
	if nt-ns > 1 {
		return false
	}

	for i := 0; i < ns; i++ {
		if s[i] != t[i] {
			if ns == nt {
				// 想一想：这里为什么不会发生索引越界呢?
				return s[i+1:] == t[i+1:]
			} else {
				return s[i:] == t[i+1:]
			}
		}
	}
	return ns+1 == nt
}
