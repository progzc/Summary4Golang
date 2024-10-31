package leetcode_0043_multiply_strings

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// 0043.字符串相乘
// https://leetcode.cn/problems/multiply-strings

func TestMultiStr(t *testing.T) {
	s1 := "98972"
	s2 := "2345"
	fmt.Println(multiStr(s1, s2)) // 232089340

	s1 = "98972"
	s2 = "10"
	fmt.Println(multiStr(s1, s2)) // 989720

	s1 = "98972"
	s2 = "0"
	fmt.Println(multiStr(s1, s2)) // 0
}

func TestAddStr(t *testing.T) {
	s1 := "98972"
	s2 := "2345"
	fmt.Println(addStr(s1, s2)) // 101317

	s1 = "98972"
	s2 = "0"
	fmt.Println(addStr(s1, s2)) // 98972

	s1 = "93"
	s2 = "7"
	fmt.Println(addStr(s1, s2)) // 100
}

// multiStr 两个字符串相乘
func multiStr(num1, num2 string) string {
	var ans string
	m, n := len(num1), len(num2)
	for j := n - 1; j >= 0; j-- {
		var tmp string
		left := 0
		for i := m - 1; i >= 0; i-- {
			x := int(num1[i] - '0')
			y := int(num2[j] - '0')
			mod := (left + x*y) % 10
			left = (left + x*y) / 10
			tmp = strconv.Itoa(mod) + tmp
		}
		if left > 0 {
			tmp = strconv.Itoa(left) + tmp
		}
		tmp = tmp + strings.Repeat("0", n-1-j)
		if len(tmp) > 0 && tmp[0] == '0' {
			tmp = "0"
		}
		ans = addStr(ans, tmp)
	}
	return ans
}

// addStr 两个字符串相加
func addStr(num1, num2 string) string {
	var ans string
	m, n := len(num1), len(num2)
	i, j := m-1, n-1
	left := 0
	for i >= 0 && j >= 0 {
		x, y := int(num1[i]-'0'), int(num2[j]-'0')
		mod := (left + x + y) % 10
		left = (left + x + y) / 10
		ans = strconv.Itoa(mod) + ans
		i--
		j--
	}
	for ; i >= 0; i-- {
		x := int(num1[i] - '0')
		mod := (left + x) % 10
		left = (left + x) / 10
		ans = strconv.Itoa(mod) + ans
	}

	for ; j >= 0; j-- {
		x := int(num2[j] - '0')
		mod := (left + x) % 10
		left = (left + x) / 10
		ans = strconv.Itoa(mod) + ans
	}

	if left > 0 {
		ans = strconv.Itoa(left) + ans
	}
	return ans
}
