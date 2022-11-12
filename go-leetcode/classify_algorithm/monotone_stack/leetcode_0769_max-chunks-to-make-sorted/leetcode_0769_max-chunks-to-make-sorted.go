package leetcode_0769_max_chunks_to_make_sorted

// 0769. 最多能完成排序的块
// https://leetcode.cn/problems/max-chunks-to-make-sorted/

// maxChunksToSorted 单调栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
// 	根据题目，我们可以发现，从左到右，每个分块都有一个最大值，并且这些分块的最大值呈单调递增。
// 	我们可以用一个栈来存储这些分块的最大值。最后得到的栈的大小，也就是题目所求的最多能完成排序的块。
// 注意: 若arr中有重复元素,那么单调栈仍适用,而贪心不再适用.
//		故本题答案仍适用于 【768. 最多能完成排序的块 II】
func maxChunksToSorted(arr []int) int {
	var stack []int
	for _, num := range arr {
		if len(stack) == 0 || num >= stack[len(stack)-1] {
			stack = append(stack, num)
		} else {
			mx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for len(stack) > 0 && num < stack[len(stack)-1] {
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, mx)
		}
	}
	return len(stack)
}
