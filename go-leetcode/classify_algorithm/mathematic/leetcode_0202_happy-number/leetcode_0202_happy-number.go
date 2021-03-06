package leetcode_0202_happy_number

// 0202.快乐数
// https://leetcode-cn.com/problems/happy-number/

// isHappy 哈希表+环
// 时间复杂度: O(log(n))
// 空间复杂度: O(log(n))
// 思路：会出现环 表示会进行无限循环
func isHappy(n int) bool {
	r := map[int]bool{}
	for {
		if n == 1 {
			return true
		}
		if r[n] {
			return false
		} else {
			r[n] = true
		}
		n = next(n)
	}
}

// isHappy_2 快慢指针+环
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
// 思路：会出现环 表示会进行无限循环
func isHappy_2(n int) bool {
	slow, fast := n, next(n)
	for fast != 1 && slow != fast {
		slow = next(slow)
		fast = next(next(fast))
	}
	return fast == 1
}

func next(n int) int {
	sum := 0
	for n > 0 {
		sum += (n % 10) * (n % 10)
		n /= 10
	}
	return sum
}
