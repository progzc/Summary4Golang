package leetcode_0084_largest_rectangle_in_histogram

import (
	"fmt"
	"testing"
)

func Test_largestRectangleArea_2(t *testing.T) {
	heights := []int{2, 1, 5, 6, 2, 3}
	fmt.Println(largestRectangleArea_2(heights))
}
