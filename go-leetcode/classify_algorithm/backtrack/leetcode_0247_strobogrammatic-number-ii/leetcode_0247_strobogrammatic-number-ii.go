package leetcode_0247_strobogrammatic_number_ii

// 0247. 中心对称数 II
// https://leetcode.cn/problems/strobogrammatic-number-ii/
// 关于中心对称数字: 0->0,1->1,8->8,6->9,9->6

// findStrobogrammatic 找规律迭代
// 时间复杂度：O(n)
// 空间复杂度: O(1)
// 思路：
// a.初始条件：
//	n=1 (3个): "0","1","8"
//	n=2 (4个): "11","69","88","96"
//	n=3 (12个): "111","181","101","888","808","818","906","916","986","609","619","689"
// b.迭代条件
//	当n为奇数时,值集合为上一个偶数值集合的每个元素中心分别加"0"/"1"/"8"得到;
//	当n为偶数时,值集合为上一个偶数值集合的每个元素中心分别加"00"/"11"/"69"/"88"/"96"得到
func findStrobogrammatic(n int) []string {
	odd, even := []string{"0", "1", "8"}, []string{"11", "69", "88", "96"}
	oddPairs, evenPairs := []string{"0", "1", "8"}, []string{"00", "11", "69", "88", "96"}

	if n == 1 {
		return odd
	} else if n == 2 {
		return even
	}

	preEven := even
	var ans []string
	for i := 3; i <= n; i++ {
		ans = make([]string, 0)
		// 偶数情况
		if i%2 == 0 {
			for _, item := range preEven {
				for _, pair := range evenPairs {
					ans = append(ans, item[:len(item)/2]+pair+item[len(item)/2:])
				}
			}
			preEven = ans
		} else {
			// 奇数情况
			for _, item := range preEven {
				for _, pair := range oddPairs {
					ans = append(ans, item[:len(item)/2]+pair+item[len(item)/2:])
				}
			}
		}
	}
	return ans
}
