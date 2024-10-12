package weipai_20241012_three_sum

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

//找出下列数组中，三个数之和加起来为0的数组的去重后列表
//[3,-3,0,-1,-2,2,1,1,-1]，输入一维数组，输出二维数组
//
//[[-3,0,3],[-2,0,2]...]

func TestSum(t *testing.T) {
	nums := []int{3, -3, 0, -1, -2, 2, 1, 1, -1}
	fmt.Println(Sum(nums)) // [[-3 0 3] [-3 1 2] [-2 -1 3] [-2 0 2] [-2 1 1] [-1 -1 2] [-1 0 1]]

	nums = []int{1, 2, -2, -1}
	fmt.Println(Sum(nums)) // []

	fmt.Println(strings.Repeat("-", 100))

	nums = []int{3, -3, 0, -1, -2, 2, 1, 1, -1}
	fmt.Println(Sum_2(nums)) // [[-3 0 3] [-3 1 2] [-2 -1 3] [-2 0 2] [-2 1 1] [-1 -1 2] [-1 0 1]]

	nums = []int{1, 2, -2, -1}
	fmt.Println(Sum_2(nums)) // []
}

// Sum 会超时
func Sum(nums []int) [][]int {
	var ans [][]int
	if len(nums) < 3 {
		return ans
	}

	sort.Ints(nums)
	fmt.Println(nums)

	n := len(nums)
	for i := 0; i < n-2; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		for j := i + 1; j < n-1; j++ {
			if j > i+1 && nums[j] == nums[j-1] {
				continue
			}
			k := n - 1
			for k > j && nums[i]+nums[j]+nums[k] > 0 {
				k--
			}
			if k > j && nums[k]+nums[j]+nums[i] == 0 {
				ans = append(ans, []int{nums[i], nums[j], nums[k]})
			}
		}
	}
	return ans
}

// Sum_2 优化
func Sum_2(nums []int) [][]int {
	var ans [][]int
	if len(nums) < 3 {
		return ans
	}

	sort.Ints(nums)
	fmt.Println(nums)

	n := len(nums)
	for i := 0; i < n-2; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		k := n - 1
		for j := i + 1; j < n-1; j++ {
			if j > i+1 && nums[j] == nums[j-1] {
				continue
			}
			for k > j && nums[i]+nums[j]+nums[k] > 0 {
				k--
			}
			if j == k {
				break
			}
			if nums[k]+nums[j]+nums[i] == 0 {
				ans = append(ans, []int{nums[i], nums[j], nums[k]})
			}
		}
	}
	return ans
}

func TestQueue(t *testing.T) {
	Queue(50)
}

func Queue(target int) {
	nums := make([]int, target)
	for i := 0; i < target; i++ {
		nums[i] = i + 1
	}
	fmt.Printf("src: %v\n", nums)

	count := 0
	for len(nums) > 1 {
		newNums := make([]int, 0, len(nums)/2)
		for i := 0; i < len(nums); i++ {
			if (i+1)%2 == 0 {
				newNums = append(newNums, nums[i])
			}
		}
		count++
		nums = newNums
		fmt.Printf("第 %d 次后, nums: %v\n", count, nums)
	}
	//  src: [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50]
	//	第 1 次后, nums: [2 4 6 8 10 12 14 16 18 20 22 24 26 28 30 32 34 36 38 40 42 44 46 48 50]
	//	第 2 次后, nums: [4 8 12 16 20 24 28 32 36 40 44 48]
	//	第 3 次后, nums: [8 16 24 32 40 48]
	//	第 4 次后, nums: [16 32 48]
	//	第 5 次后, nums: [32]
}
