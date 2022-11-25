package haomo_20221124_character_string_add

import (
	"fmt"
	"strconv"
)

// 毫末笔试题:
// 两个字符串相加

func Add(s1, s2 string) string {
	n1, n2 := len(s1)-1, len(s2)-1
	carry := 0
	var ans string
	for n1 >= 0 || n2 >= 0 {
		var num1, num2 int
		if n1 >= 0 {
			num1, _ = strconv.Atoi(string(s1[n1]))
			n1--
		}
		if n2 >= 0 {
			num2, _ = strconv.Atoi(string(s2[n2]))
			n2--
		}
		tmp := num1 + num2 + carry
		cur := tmp % 10
		carry = tmp / 10
		ans = fmt.Sprintf("%d", cur) + ans
	}
	if carry > 0 {
		ans = fmt.Sprintf("%d", carry) + ans
	}
	return ans
}
