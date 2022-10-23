package leetcode_0567_permutation_in_string

// 0567. 字符串的排列
// https://leetcode.cn/problems/permutation-in-string/

// checkInclusion 固定滑动窗口
// 时间复杂度: O(n)
// 空间复杂度：O(A),A=26
// 思路：
//	由于排列不会改变字符串中每个字符的个数，所以只有当两个字符串每个字符的个数均相等时，一个字符串才是另一个字符串的排列。
//	根据这一性质，记 s1 的长度为 n，我们可以遍历 s2 中的每个长度为 n 的子串，判断子串和 s1 中每个字符的个数是否相等，
//	若相等则说明该子串是 s1 的一个排列。
func checkInclusion(s1 string, s2 string) bool {
	n, m := len(s1), len(s2)
	if n > m {
		return false
	}

	cnt1, cnt2 := [26]int{}, [26]int{}
	for i, ch := range s1 {
		cnt1[ch-'a']++
		cnt2[s2[i]-'a']++
	}
	// 注意：数组可以直接判等。但是这里效率太低，可以优化（见 固定滑动窗口(优化)）
	if cnt1 == cnt2 {
		return true
	}

	for i := n; i < m; i++ {
		cnt2[s2[i]-'a']++
		cnt2[s2[i-n]-'a']--
		// 注意：数组可以直接判等。但是这里效率太低，可以优化（见 固定滑动窗口(优化)）
		if cnt1 == cnt2 {
			return true
		}
	}
	return false
}

// checkInclusion_2 固定滑动窗口(优化)
// 时间复杂度: O(n)
// 空间复杂度：O(A),A=26
// 思路：
//	由于排列不会改变字符串中每个字符的个数，所以只有当两个字符串每个字符的个数均相等时，一个字符串才是另一个字符串的排列。
//	根据这一性质，记 s1 的长度为 n，我们可以遍历 s2 中的每个长度为 n 的子串，判断子串和 s1 中每个字符的个数是否相等，
//	若相等则说明该子串是 s1 的一个排列。
// 优化：
//	注意到每次窗口滑动时，只统计了一进一出两个字符，却比较了整个 cnt1 和 cnt2 数组。 从这个角度出发，
//	我们可以用一个变量 diff 来记录 cnt1 与 cnt2 的不同值的个数，这样判断 cnt1 和 cnt2 是否相等就转换成了判断 diff 是否为 0。
func checkInclusion_2(s1 string, s2 string) bool {
	n, m := len(s1), len(s2)
	if n > m {
		return false
	}

	cnt, diff := [26]int{}, 0
	for i, ch := range s1 {
		// 注意：这里--/++不能颠倒
		cnt[ch-'a']--
		cnt[s2[i]-'a']++
	}
	for _, c := range cnt[:] {
		if c != 0 {
			diff++
		}
	}
	if diff == 0 {
		return true
	}

	for i := n; i < m; i++ {
		x, y := s2[i]-'a', s2[i-n]-'a'
		if x == y {
			continue
		}

		// 进之前若cnt[x]==0，则diff++，进之后若cnt[x]==0，则diff--
		if cnt[x] == 0 {
			diff++
		}
		// 这里表示 进
		cnt[x]++
		if cnt[x] == 0 {
			diff--
		}

		// 出之前若cnt[y]==0，则diff++，进之后若cnt[y]==0，则diff--
		if cnt[y] == 0 {
			diff++
		}
		// 这里表示 出
		cnt[y]--
		if cnt[y] == 0 {
			diff--
		}

		if diff == 0 {
			return true
		}
	}
	return false
}

// checkInclusion_3 双指针(很棒的思路，推荐)
// 时间复杂度: O(n)
// 空间复杂度：O(A),A=26
func checkInclusion_3(s1 string, s2 string) bool {
	n, m := len(s1), len(s2)
	if n > m {
		return false
	}

	cnt := [26]int{}
	for _, ch := range s1 {
		cnt[ch-'a']--
	}

	for l, r := 0, 0; r < m; r++ {
		x := s2[r] - 'a'
		cnt[x]++
		for cnt[x] > 0 {
			cnt[s2[l]-'a']--
			l++
		}
		if r-l+1 == n {
			return true
		}
	}
	return false
}
