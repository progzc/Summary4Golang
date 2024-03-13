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
		'2': "abc",
		'3': "def",
		'4': "ghi",
		'5': "jkl",
		'6': "mno",
		'7': "pqrs",
		'8': "tuv",
		'9': "wxyz",
	}
	var (
		ans []string
		dfs func(idx int, item string)
	)
	dfs = func(idx int, item string) {
		if idx == len(digits) {
			// 由于ans相当于全局变量，所以这里不需要使用切片指针。
			// 但如果通过函数参数传递ans，则需要使用切片指针
			ans = append(ans, item)
			return
		}
		// 注意事项：下面换成letters := m[digits[idx]-'0']会出错
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

// letterCombinations_3 回溯
// 时间复杂度: O(3^m*4^n)
// 空间复杂度: O(m+n)
//	其中m是输入中对应3个字母的数字个数（包括数字2、3、4、5、6、8）;
//	n是输入中对应4个字母的数字个数（包括数字7、9）;
//	m+n是输入数字的总个数。
func letterCombinations_3(digits string) []string {
	pair := map[byte]string{
		'2': "abc",
		'3': "def",
		'4': "ghi",
		'5': "jkl",
		'6': "mno",
		'7': "pqrs",
		'8': "tuv",
		'9': "wxyz",
	}
	var dfs func(idx int, cur string)
	var ans []string
	if len(digits) == 0 {
		return ans
	}
	dfs = func(idx int, cur string) {
		if idx == len(digits) {
			ans = append(ans, cur)
			return
		}
		v, ok := pair[digits[idx]]
		if !ok {
			return
		}
		for i := 0; i < len(v); i++ {
			dfs(idx+1, cur+string(v[i]))
		}
	}
	dfs(0, "")
	return ans
}

// letterCombinations 回溯（切片指针的使用）
// 时间复杂度: O(3^m*4^n)
// 空间复杂度: O(m+n)
//	其中m是输入中对应3个字母的数字个数（包括数字2、3、4、5、6、8）;
//	n是输入中对应4个字母的数字个数（包括数字7、9）;
//	m+n是输入数字的总个数。
func letterCombinations_2(digits string) []string {
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
	)
	if len(digits) == 0 {
		return ans
	}
	backtrack(&ans, 0, digits, "", m)
	return ans
}

func backtrack(ans *[]string, idx int, digits, item string, m map[byte]string) {
	if idx == len(digits) {
		// 但如果通过函数参数传递ans，则需要使用切片指针
		*ans = append(*ans, item)
		return
	}
	// 注意事项：下面换成letters := m[digits[idx]]会出错
	letters := m[digits[idx]-'0']
	for i := 0; i < len(letters); i++ {
		backtrack(ans, idx+1, digits, item+string(letters[i]), m)
	}
}
