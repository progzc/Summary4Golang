package leetcode_0249_group_shifted_strings

// 0249. 移位字符串分组
// https://leetcode.cn/problems/group-shifted-strings/

// groupStrings 哈希算法+map
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：设计一个哈希函数
func groupStrings(strings []string) [][]string {
	var (
		ans [][]string
		m   = make(map[string][]string)
	)

	for _, str := range strings {
		key := hash(str)
		if len(m[key]) == 0 {
			m[key] = make([]string, 0)
		}
		m[key] = append(m[key], str)
	}

	for _, v := range m {
		ans = append(ans, v)
	}

	return ans
}

func hash(str string) string {
	var (
		offset = str[0] - 'a'
		ans    = []byte(str)
	)
	for i := 0; i < len(ans); i++ {
		if ans[i]-'a' < offset {
			ans[i] += 26 - offset
		} else {
			ans[i] -= offset
		}
	}
	return string(ans)
}
