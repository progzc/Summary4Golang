package order_leetcode

// link:https://stackoverflow.com/questions/36000487/check-for-equality-on-slices-without-order
func sameIntSlice(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}
	diff := make(map[int]int, len(x))
	for _, _x := range x {
		diff[_x]++
	}
	for _, _y := range y {
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}
