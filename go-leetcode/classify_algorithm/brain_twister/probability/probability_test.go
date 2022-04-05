package probability

import "math/rand"

// 脑筋急转弯：有三个奖品，每个有各自的概率，设计一个函数返回抽到的奖品。

func GetThings() string {
	t1 := 0.1 // 物品A抽到的概率为0.1
	t2 := 0.2 // 物品B抽到的概率为0.2
	t3 := 0.7 // 物品C抽到的概率为0.7

	x := rand.Float64()
	if x >= 0 && x < t1 {
		return "A"
	}
	if x >= t1 && x < (t1+t2) {
		return "B"
	}
	if x >= (t1+t2) && x < (t1+t2+t3) {
		return "C"
	}
	return "A"
}
