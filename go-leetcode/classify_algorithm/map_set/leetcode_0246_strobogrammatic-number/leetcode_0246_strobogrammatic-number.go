package leetcode_0246_strobogrammatic_number

// 0246. 中心对称数
// https://leetcode.cn/problems/strobogrammatic-number/

// isStrobogrammatic 哈希表
// 时间复杂度：O(n)
// 空间复杂度: O(1)
func isStrobogrammatic(num string) bool {
	lookup := map[byte]byte{
		'0': '0',
		'1': '1',
		'6': '9',
		'8': '8',
		'9': '6',
	}
	l, r := 0, len(num)-1
	for l <= r {
		if lookup[num[l]] != num[r] {
			return false
		} else {
			l++
			r--
		}
	}
	return true
}
