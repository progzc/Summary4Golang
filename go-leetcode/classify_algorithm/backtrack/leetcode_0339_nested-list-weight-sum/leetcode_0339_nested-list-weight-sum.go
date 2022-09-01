package leetcode_0339_nested_list_weight_sum

// 0339. 嵌套列表权重和
// https://leetcode.cn/problems/nested-list-weight-sum/

type NestedInteger struct {
	Num int
	Ns  []*NestedInteger
}

func (n NestedInteger) IsInteger() bool {
	return n.Ns == nil
}

func (n NestedInteger) GetInteger() int {
	return n.Num
}

func (n *NestedInteger) SetInteger(value int) {
	n.Num = value
}

func (n *NestedInteger) Add(elem NestedInteger) {
	n.Ns = append(n.Ns, &elem)
}

func (n NestedInteger) GetList() []*NestedInteger {
	return n.Ns
}

// depthSum 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func depthSum(nestedList []*NestedInteger) int {
	var dfs func(list []*NestedInteger, depth int) int
	dfs = func(list []*NestedInteger, depth int) int {
		sum := 0
		for _, x := range list {
			if x.IsInteger() {
				sum += x.GetInteger() * depth
			} else {
				sum += dfs(x.GetList(), depth+1)
			}
		}
		return sum
	}
	return dfs(nestedList, 1)
}
