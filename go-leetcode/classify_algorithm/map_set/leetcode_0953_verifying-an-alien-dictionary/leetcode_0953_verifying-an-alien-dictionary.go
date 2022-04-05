package leetcode_0953_verifying_an_alien_dictionary

// 0953.验证外星语词典
// https://leetcode-cn.com/problems/verifying-an-alien-dictionary/

// isAlienSorted 哈希表
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func isAlienSorted(words []string, order string) bool {
	rank := [26]int{}
	for i, o := range []byte(order) {
		rank[o-'a'] = i + 1
	}

	var compare = func(s1, s2 string) bool {
		l1, l2 := len(s1), len(s2)
		i, j := 0, 0
		for ; i < l1 && j < l2; i, j = i+1, j+1 {
			if rank[s1[i]-'a'] < rank[s2[j]-'a'] {
				return true
			} else if rank[s1[i]-'a'] > rank[s2[j]-'a'] {
				return false
			} else {
				continue
			}
		}
		if i == l1 && j != l2 {
			return true
		}

		if i != l1 && j == l2 {
			return false
		}
		return true
	}

	for i := 0; i < len(words)-1; i++ {
		if !compare(words[i], words[i+1]) {
			return false
		}
	}
	return true
}
