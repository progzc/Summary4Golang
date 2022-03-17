package leetcode_0191_number_of_1_bits

// 0191.位1的个数
// https://leetcode-cn.com/problems/number-of-1-bits/

// hammingWeight 位操作
// 时间复杂度: O(k)
// 空间复杂度: O(1)
func hammingWeight(num uint32) int {
	count := 0
	for i := 0; i < 32; i++ {
		if (num>>i)&1 > 0 {
			count++
		}
	}
	return count
}

// hammingWeight_2 位运算的优化
// 时间复杂度: O(logn)
// 空间复杂度: O(1)
// 思路：n&(n−1) 会把最低位的1变为0
func hammingWeight_2(num uint32) int {
	count := 0
	for ; num > 0; num &= num - 1 {
		count++
	}
	return count
}
