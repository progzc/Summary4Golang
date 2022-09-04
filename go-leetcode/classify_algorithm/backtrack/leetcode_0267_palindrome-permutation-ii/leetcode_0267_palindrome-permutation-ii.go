package leetcode_0267_palindrome_permutation_ii

import "sort"

// 0267. 回文排列 II
// https://leetcode.cn/problems/palindrome-permutation-ii/

// generatePalindromes 排序组合去重
func generatePalindromes(s string) []string {
	var (
		ans []string

		half   []byte           // 只考虑一半字符串的组合
		odd    string           // 奇数字符最多只能存在一个
		oddCnt int              // 奇数字符的数量
		m      = map[byte]int{} // 计算每个字符的数量
	)

	// 快速判断
	if len(s) == 0 {
		return ans
	} else if len(s) == 1 {
		ans = append(ans, s)
		return ans
	}

	// 取一半的字符
	for _, b := range []byte(s) {
		m[b]++
	}
	for k, v := range m {
		if v%2 != 0 {
			odd = string(k)
			oddCnt++
			if oddCnt > 1 {
				return ans
			}
		}
		v /= 2
		for i := 0; i < v; i++ {
			half = append(half, k)
		}
	}

	// 排列组合(包含去重)
	bs := permuteUnique(half)
	for _, v := range bs {
		ans = append(ans, string(v)+odd+reverse(string(v)))
	}

	return ans
}

func reverse(s string) string {
	str := []byte(s)
	left, right := 0, len(str)-1
	for left < right {
		str[left], str[right] = str[right], str[left]
		left++
		right--
	}
	return string(str)
}

// permuteUnique 深度遍历
// 时间复杂度：O(n*n!)
// 空间复杂度：O(2n)
func permuteUnique(nums []byte) [][]byte {
	if len(nums) == 0 {
		return [][]byte{}
	}
	// 必需先排序
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	var (
		ans    [][]byte
		used   = make([]bool, len(nums))
		output []byte
		dfs    func(idx int, output []byte)
	)
	// dfs 表示从左往右填到第idx个位置,当前排列为output
	// 其中第0~idx-1的元素均已经填充过了，而idx~n-1的元素还未填充过
	dfs = func(idx int, output []byte) {
		// 终止条件
		if idx == len(nums) {
			// 注意事项：ans = append(ans,output)这种写法是错误的
			ans = append(ans, append([]byte(nil), output...))
		}
		for i := 0; i < len(nums); i++ {
			// 剪枝：只可以选择未选择过的数
			if used[i] {
				continue
			}
			// 剪枝：只可以选择之前未选择过的不同的数
			if i > 0 && !used[i-1] && nums[i] == nums[i-1] {
				continue
			}

			// 选择：选择第i个位置作为第idx个数
			output = append(output, nums[i])
			used[i] = true
			// 递归：填下一个位置
			dfs(idx+1, output)
			// 回溯：在下一次选择之前，必须回撤销上一次的选择
			used[i] = false
			output = output[:len(output)-1]
		}
	}
	dfs(0, output)
	return ans
}
