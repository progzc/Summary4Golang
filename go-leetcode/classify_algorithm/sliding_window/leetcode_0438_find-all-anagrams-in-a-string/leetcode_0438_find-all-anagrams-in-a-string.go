package leetcode_0438_find_all_anagrams_in_a_string

// 0438. 找到字符串中所有字母异位词
// https://leetcode.cn/problems/find-all-anagrams-in-a-string/

// 与 https://leetcode.cn/problems/permutation-in-string/ 类似

// findAnagrams 滑动窗口
// 时间复杂度: O(n)
// 空间复杂度: O(A),A=26
func findAnagrams(s string, p string) []int {
	n, m := len(p), len(s)
	if n > m {
		return nil
	}

	var ans []int
	cnt1, cnt2 := [26]byte{}, [26]byte{}
	for i, ch := range p {
		cnt1[ch-'a']++
		cnt2[s[i]-'a']++
	}
	if cnt1 == cnt2 {
		ans = append(ans, 0)
	}

	for i := n; i < m; i++ {
		cnt2[s[i]-'a']++
		cnt2[s[i-n]-'a']--
		if cnt1 == cnt2 {
			ans = append(ans, i-n+1)
		}
	}
	return ans
}

// findAnagrams_2 滑动窗口（优化）
// 时间复杂度: O(n)
// 空间复杂度: O(A),A=26
func findAnagrams_2(s string, p string) []int {
	n, m := len(p), len(s)
	if n > m {
		return nil
	}

	cnt, diff := [26]int{}, 0
	var ans []int
	for i, ch := range p {
		// 注意：这里--/++不能颠倒
		cnt[ch-'a']--
		cnt[s[i]-'a']++
	}
	for _, c := range cnt[:] {
		if c != 0 {
			diff++
		}
	}
	if diff == 0 {
		ans = append(ans, 0)
	}

	for i := n; i < m; i++ {
		x, y := s[i]-'a', s[i-n]-'a'
		// 注意：对比 [0567. 字符串的排列]，这里需要注释掉，不然会出错
		//if x == y {
		//	continue
		//}

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
			ans = append(ans, i-n+1)
		}
	}
	return ans
}

// findAnagrams_3 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(A),A=26
func findAnagrams_3(s string, p string) []int {
	n, m := len(p), len(s)
	if n > m {
		return nil
	}

	cnt := [26]int{}
	for i := 0; i < n; i++ {
		cnt[p[i]-'a']--
	}

	var ans []int
	for l, r := 0, 0; r < m; r++ {
		x := s[r] - 'a'
		cnt[x]++
		for cnt[x] > 0 {
			cnt[s[l]-'a']--
			l++
		}
		if r-l+1 == n {
			ans = append(ans, l)
		}
	}
	return ans
}
