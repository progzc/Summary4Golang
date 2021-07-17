package go_concurrency

import (
	"crypto/rand"
	"fmt"
	"testing"
)

// TestGenerateRandNum 生成指定位数的随机数
func TestGenerateRandNum(t *testing.T) {
	fmt.Println(RandToken(5))
}

// 生成num*2位的字符串
func RandToken(num int) string {
	b := make([]byte, num)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
