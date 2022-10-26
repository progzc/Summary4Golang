package leetcode_1215_stepping_numbers

import "fmt"

// 1215. 步进数
// https://leetcode.cn/problems/stepping-numbers/

// countSteppingNumbers 暴力法(超时)
// 时间复杂度: O((high-low)*len(high))
// 空间复杂度: O(1)
func countSteppingNumbers(low int, high int) []int {
	var ans []int
	for num := low; num <= high; num++ {
		if isSteppingNumbers(num) {
			ans = append(ans, num)
		}
	}
	return ans
}

// countSteppingNumbers_2 bfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func countSteppingNumbers_2(low int, high int) []int {
	var (
		ans   []int
		queue []int
	)
	if low == 0 {
		ans = append(ans, 0)
	}
	// 注意：这里不能加入0
	for i := 1; i <= 9; i++ {
		queue = append(queue, i)
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur > high {
			return ans
		}
		if cur >= low && cur <= high {
			ans = append(ans, cur)
		}
		lastDigit := cur % 10
		// 注意：这里不能==0，不然会发生进位
		if lastDigit > 0 {
			queue = append(queue, cur*10+lastDigit-1)
		}
		// 注意：这里不能==9，不然会发生进位
		if lastDigit < 9 {
			queue = append(queue, cur*10+lastDigit+1)
		}
	}
	return ans
}

func isSteppingNumbers(num int) bool {
	str := fmt.Sprintf("%d", num)
	if len(str) == 1 {
		return true
	}
	for i := 1; i < len(str); i++ {
		if abs(int(str[i])-int(str[i-1])) != 1 {
			return false
		}
	}
	return true
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}
