package leetcode_1256_encode_number

import (
	"fmt"
)

// 1256. 加密数字
// https://leetcode.cn/problems/encode-number/

// encode 找规律
// 时间复杂度: O(1)
// 空间复杂度: O(1)
// 思路: 将 num+1 再去掉最高位即可。
func encode(num int) string {
	return fmt.Sprintf("%b", num+1)[1:]
}
