package leetcode_0277_find_the_celebrity

// 0277. 搜寻名人
// https://leetcode.cn/problems/find-the-celebrity/

// solution 脑筋急转弯
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：
//	我们不妨先假定ans为0，然后迭代n个人，如果此时ans认识某个人k(0<=k<n),那么令ans为k；
//	a.如何证明若存在名人，则名人必定为ans呢？ 我们知道如果存在名人，那么在迭代的过程必定会遇到名人,并且此时ans认识名人，
//	  不管ans此时是不是名人，所以此时令ans=名人。在接下来的迭代中，由于名人不认识其他人，则必然不会发生ans值的变更。
//	  所以可知若存在名人，则ans必为名人。
//	b.得到ans后，我们需要判断ans是不是名人，这个判断过程很简单，就不细说了
// 总结：两两比较总能淘汰1个，n轮淘汰下来，淘汰n-1个，剩下1个，而题目给出名人最多1个，再检查这个剩下的候选名人即可。
func solution(knows func(a int, b int) bool) func(n int) int {
	return func(n int) int {
		ans := 0
		// 查找名人
		for k := 0; k < n; k++ {
			if knows(ans, k) {
				ans = k
			}
		}

		// 判断名人
		for k := 0; k < n; k++ {
			if k != ans {
				// 若其他人有不认识名人的，则名人不存在；
				// 若名人认识其他人的，则名人不存在
				if !knows(k, ans) || knows(ans, k) {
					return -1
				}
			}
		}
		return ans
	}
}
