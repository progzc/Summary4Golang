package leetcode_0251_flatten_2d_vector

// 0251. 展开二维向量
// https://leetcode.cn/problems/flatten-2d-vector/

// Vector2D 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
type Vector2D struct {
	nums [][]int
	i, j int
}

func Constructor(vec [][]int) Vector2D {
	return Vector2D{
		nums: vec,
	}
}

func (this *Vector2D) Next() int {
	v := this.nums[this.i][this.j]
	if this.j+1 < len(this.nums[this.i]) {
		this.j++
	} else {
		this.i++
		this.j = 0
	}
	return v
}

func (this *Vector2D) HasNext() bool {
	if this.i < len(this.nums) {
		if this.j < len(this.nums[this.i]) {
			return true
		} else {
			this.i++
			this.j = 0
			return this.HasNext()
		}
	}
	return false
}
