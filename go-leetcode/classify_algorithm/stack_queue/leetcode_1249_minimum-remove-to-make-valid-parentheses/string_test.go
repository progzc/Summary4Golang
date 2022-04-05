package leetcode_1249_minimum_remove_to_make_valid_parentheses

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	str := "abcd"
	k := len(str) - 1
	//  注意：str[k+1:]不会出现索引越界
	fmt.Println(string(append([]byte(str[:k]), str[k+1:]...)))
}
