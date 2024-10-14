package leetcode_0763_partition_labels

// 0763. 划分字母区间🌟
// https://leetcode.cn/problems/partition-labels

// partitionLabels
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func partitionLabels(s string) []int {
	var ans []int
	n := len(s)
	if n == 0 {
		return ans
	}
	if n == 1 {
		ans = []int{1}
		return ans
	}

	m := make(map[byte]int)
	for i := 0; i < n; i++ {
		m[s[i]] = i
	}

	i := 0
	for i < n {
		end := m[s[i]]
		j := i
		for ; j <= end; j++ {
			if m[s[j]] > end {
				end = m[s[j]]
			}
		}
		ans = append(ans, j-i)
		i = j
	}
	return ans
}
