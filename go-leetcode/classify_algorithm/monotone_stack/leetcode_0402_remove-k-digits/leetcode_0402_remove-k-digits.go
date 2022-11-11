package leetcode_0402_remove_k_digits

import "strings"

// 0402. 移掉 K 位数字
// https://leetcode.cn/problems/remove-k-digits/

// removeKdigits 单调递增（非严格）栈
// 时间复杂度: O(nk)
// 空间复杂度: O(n)
// 思路:
//	a.丢弃k位数字,相当于保留n-k位数字
//	b.对于两个相同长度的数字序列，最左边不同的数字决定了这两个数字的大小，例如，对于 A = 1axxx，B = 1bxxx，如果 a > b， 则 A > B。
//	  基于此，我们可以知道，若要使得剩下的数字最小，需要保证靠前的数字尽可能小。
// 具体方案：
//	a.考虑从左往右增量的构造最后的答案。我们可以用一个栈维护当前的答案序列，栈中的元素代表截止到当前位置，
//	  删除不超过 k 次个数字后，所能得到的最小整数。根据之前的讨论：在使用 k 个删除次数之前，栈中的序列从栈底到栈顶单调不降。
//	  因此，对于每个数字，如果该数字小于栈顶元素，我们就不断地弹出栈顶元素，直到
//		i)栈为空
//		ii)或者新的栈顶元素不大于当前数字
//		iii)或者我们已经删除了 k 位数字
//	b.上述步骤结束后我们还需要针对一些情况做额外的处理：
//		i)如果我们删除了 m 个数字且 m<k，这种情况下我们需要从序列尾部删除额外的 k-m 个数字。
//		ii)如果最终的数字序列存在前导零，我们要删去前导零。
//		iii)如果最终数字序列为空，我们应该返回 0。
func removeKdigits(num string, k int) string {
	var stack []byte
	for _, digit := range []byte(num) {
		for k > 0 && len(stack) > 0 && digit < stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
			k--
		}
		stack = append(stack, digit)
	}
	stack = stack[:len(stack)-k]
	ans := strings.TrimLeft(string(stack), "0")
	if ans == "" {
		ans = "0"
	}
	return ans
}
