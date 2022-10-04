package leetcode_0755_pour_water

// 0755. 倒水
// https://leetcode.cn/problems/pour-water/

// pourWater 模拟
// 时间复杂度: O(v*n)
// 空间复杂度: O(1)
// 不足：代码太繁琐了
func pourWater(heights []int, volume int, k int) []int {
	for ; volume > 0; volume-- {
		find := false
		// 检查左边
		for left := k - 1; left >= 0; left-- {
			if heights[left] > heights[k] {
				break
			} else if heights[left] < heights[k] {
				find = true
				real := left
				for next := left - 1; next >= 0; next-- {
					if heights[next] > heights[next+1] {
						break
					} else if heights[next] < heights[next+1] {
						real = next
						continue
					} else {
						continue
					}
				}
				heights[real]++
				break
			} else {
				continue
			}
		}

		if !find {
			// 检查右边
			for right := k + 1; right < len(heights); right++ {
				if heights[right] > heights[k] {
					break
				} else if heights[right] < heights[k] {
					find = true
					real := right
					// 还需要检查右边有没有更小的
					for next := right + 1; next < len(heights); next++ {
						if heights[next] > heights[next-1] {
							break
						} else if heights[next] < heights[next-1] {
							real = next
							continue
						} else {
							continue
						}
					}
					heights[real]++
					break
				} else {
					continue
				}
			}
		}

		// 左边右边都没找到
		if !find {
			heights[k]++
		}
	}
	return heights
}

// pourWater 模拟
// 时间复杂度: O(v*n)
// 空间复杂度: O(1)
// 优点：代码更加简洁
func pourWater_2(heights []int, volume int, k int) []int {
	for ; volume > 0; volume-- {
		for d := -1; d <= 1; d += 2 {
			i, best := k, k
			for i+d >= 0 && i+d < len(heights) && heights[i+d] <= heights[i] {
				if heights[i+d] < heights[i] {
					best = i + d
				}
				i += d
			}
			if heights[best] < heights[k] {
				heights[best]++
				goto next
			}
		}
		heights[k]++
	next:
	}
	return heights
}
