package leetcode_0013_roman_to_integer

// 0013.罗马数字转整数
// https://leetcode.cn/problems/roman-to-integer/

// romanToInt 散列表
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func romanToInt(s string) int {
	single := map[byte]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}

	recomb := map[string]int{
		"IV": 4,
		"IX": 9,
		"XL": 40,
		"XC": 90,
		"CD": 400,
		"CM": 900,
	}

	sum := 0
	i, n := 0, len(s)

	for i < n-1 {
		v, ok := recomb[s[i:i+2]]
		if ok {
			sum += v
			i += 2
		} else {
			sum += single[s[i]]
			i += 1
		}
	}
	if i == n-1 {
		sum += single[s[i]]
	}

	return sum
}

// romanToInt_2 散列表（优化）
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func romanToInt_2(s string) int {
	single := map[byte]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}

	sum, n := 0, len(s)
	for i := range s {
		v := single[s[i]]
		if i < n-1 && v < single[s[i+1]] {
			sum -= v
		} else {
			sum += v
		}
	}
	return sum
}
