package leetcode_0017_letter_combinations_of_a_phone_number

// 0017.电话号码的字母组合
// https://leetcode-cn.com/problems/letter-combinations-of-a-phone-number/

// letterCombinations 回溯
// 时间复杂度: O(3^m*4^n)
// 空间复杂度: O(m+n)
//	其中m是输入中对应3个字母的数字个数（包括数字2、3、4、5、6、8）;
//	n是输入中对应4个字母的数字个数（包括数字7、9）;
//	m+n是输入数字的总个数。
func letterCombinations(digits string) []string {
	m := map[byte]string{
		2: "abc",
		3: "def",
		4: "ghi",
		5: "jkl",
		6: "mno",
		7: "pqrs",
		8: "tuv",
		9: "wxyz",
	}
	var (
		ans []string
		dfs func(idx int, item string)
	)
	dfs = func(idx int, item string) {
		if idx == len(digits) {
			ans = append(ans, item)
			return
		}
		letters := m[digits[idx]]
		for i := 0; i < len(letters); i++ {
			dfs(idx+1, item+string(letters[i]))
		}
	}
	if len(digits) == 0 {
		return ans
	}
	dfs(0, "")
	return ans
}
