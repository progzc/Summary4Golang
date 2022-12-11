package interview_guide_17_09_get_kth_magic_number_lcci

// 面试题 17.09. 第 k 个数
// https://leetcode.cn/problems/get-kth-magic-number-lcci/

// getKthMagicNumber
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	概念: 把符合题目要求的这些数叫做丑数。
//	规律: 一个丑数总是由前面的某一个丑数 x3 / x5 / x7 得到
//	推导:
//	 若把丑数数列叫做 ugly[i]，那么考虑一下三个数列:
//		1. ugly[0]*3, ugly[1]*3, ugly[2]*3, ugly[3]*3, ugly[4]*3, ugly[5]*3……
//		2. ugly[0]*5, ugly[1]*5, ugly[2]*5, ugly[3]*5, ugly[4]*5, ugly[5]*5……
//		3. ugly[0]*7, ugly[1]*7, ugly[2]*7, ugly[3]*7, ugly[4]*7, ugly[5]*7……
//	 上面这个三个数列合在一起就形成了新的、更长的丑数数列.
//	 如果合在一起呢？这其实就是一个合并有序线性表的问题。
func getKthMagicNumber(k int) int {
	ugly := make([]int, k)
	p1, p2, p3 := 0, 0, 0
	ugly[0] = 1
	for i := 1; i < k; i++ {
		ugly[i] = min(min(ugly[p1]*3, ugly[p2]*5), ugly[p3]*7)
		if ugly[i] == ugly[p1]*3 {
			p1++
		}
		if ugly[i] == ugly[p2]*5 {
			p2++
		}
		if ugly[i] == ugly[p3]*7 {
			p3++
		}
	}
	return ugly[k-1]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
