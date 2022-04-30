package leetcode_0202_happy_number

// 0202.快乐数
// https://leetcode-cn.com/problems/happy-number/

// isHappy 哈希表+环
// 时间复杂度: O(log(n))
// 空间复杂度: O(log(n))
// 思路：会出现环 表示会进行无限循环
func isHappy(n int) bool {
	r := map[int]bool{}
	for n != 1 && !r[n] {
		n, r[n] = next(n), true
	}
	return n == 1
}

func next(n int) int {
	sum := 0
	for n > 0 {
		sum += (n % 10) * (n % 10)
		n /= 10
	}
	return sum
}
