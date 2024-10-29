package autowise_20241028_power_exp

import (
	"fmt"
	"testing"
)

// 仙途智能（一面）
// 2的幂次方表示
// https://www.luogu.com.cn/problem/P1010

func TestPowerExp(t *testing.T) {
	fmt.Println(PowerExp(137))  // 2(2(2)+2+2(0))+2(2+2(0))+2(0)
	fmt.Println(PowerExp(1315)) // 2(2(2+2(0))+2)+2(2(2+2(0)))+2(2(2)+2(0))+2+2(0)
}

func PowerExp(n int) string {
	var ans string
	addFlag := false
	for i := 32; i >= 0; i-- {
		if (n>>i)&1 == 1 {
			if addFlag {
				ans = ans + "+"
			} else {
				addFlag = true
			}
			if i == 0 {
				ans = ans + "2(0)"
			} else if i == 1 {
				ans = ans + "2"
			} else {
				ans = ans + fmt.Sprintf("2(%s)", PowerExp(i))
			}
		}
	}
	return ans
}
