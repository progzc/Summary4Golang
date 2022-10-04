package leetcode_1057_campus_bikes

import "sort"

// 1057. 校园自行车分配
// https://leetcode.cn/problems/campus-bikes/

// assignBikes 有序的map+排序
// 时间复杂度: O(m*n)
// 空间复杂度: O(m+n)
// 思路：
//	a.使用一个有序的map结构：key为worker到bike的距离；value为{工人序号,自行车序号}的vector，
//	  遍历计算所有工人到所有自行车的距离，若有相同距离自然按照序号从小到大放到vector里，
//	  最后遍历这个有序的map，得到答案。
//	b.其中使用了2个bool数组标记已使用过的工人和自行车。
func assignBikes(workers [][]int, bikes [][]int) []int {
	w, b := len(workers), len(bikes)
	ans := make([]int, w)
	dm, ks := make(map[int][][2]int), make([]int, 0)
	wUsed := make([]bool, w)
	bUsed := make([]bool, b)

	for i := 0; i < w; i++ {
		for j := 0; j < b; j++ {
			distance := abs(workers[i][0]-bikes[j][0]) + abs(workers[i][1]-bikes[j][1])
			if _, ok := dm[distance]; !ok {
				ks = append(ks, distance)
			}
			dm[distance] = append(dm[distance], [2]int{i, j})
		}
	}

	sort.SliceStable(ks, func(i, j int) bool {
		return ks[i] < ks[j]
	})

	for _, distance := range ks {
		v := dm[distance]
		for _, pair := range v {
			if wUsed[pair[0]] || bUsed[pair[1]] {
				continue
			}
			ans[pair[0]] = pair[1]
			wUsed[pair[0]], bUsed[pair[1]] = true, true
		}
	}
	return ans
}

type Node struct {
	Distance int
	Worker   int
	Bike     int
}

type Pair []Node

func (p Pair) Len() int {
	return len(p)
}

func (p Pair) Less(i, j int) bool {
	if p[i].Distance != p[j].Distance {
		return p[i].Distance < p[j].Distance
	}
	if p[i].Worker != p[j].Worker {
		return p[i].Worker < p[j].Worker
	}
	return p[i].Bike < p[j].Bike
}

func (p Pair) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// assignBikes_2 排序(效率比"有序的map+排序"低许多)
// 时间复杂度: O(m*n)
// 空间复杂度: O(m+n)
// 思路：
//	a.使用一个有序的map结构：key为worker到bike的距离；value为{工人序号,自行车序号}的vector，
//	  遍历计算所有工人到所有自行车的距离，若有相同距离自然按照序号从小到大放到vector里，
//	  最后遍历这个有序的map，得到答案。
//	b.其中使用了2个bool数组标记已使用过的工人和自行车。
func assignBikes_2(workers [][]int, bikes [][]int) []int {
	w, b := len(workers), len(bikes)
	wUsed := make([]bool, w)
	bUsed := make([]bool, b)
	records := Pair{}
	ans := make([]int, w)

	for i := 0; i < w; i++ {
		for j := 0; j < b; j++ {
			node := Node{}
			node.Distance = abs(workers[i][0]-bikes[j][0]) + abs(workers[i][1]-bikes[j][1])
			node.Worker = i
			node.Bike = j
			records = append(records, node)
		}
	}

	sort.Sort(records)
	for i := 0; i < len(records); i++ {
		if wUsed[records[i].Worker] || bUsed[records[i].Bike] {
			continue
		}
		ans[records[i].Worker] = records[i].Bike
		wUsed[records[i].Worker], bUsed[records[i].Bike] = true, true
	}
	return ans
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}
