package leetcode_0742_closest_leaf_in_a_binary_tree

import (
	"testing"
)

func Test_map(t *testing.T) {
	g := make(map[*TreeNode][]*TreeNode)
	g[nil] = []*TreeNode{
		{
			Val: 1,
		},
	}
	node := &TreeNode{
		Val: 2,
	}
	g[node] = nil

	// len(g): 2, g: map[<nil>:[0xc000004540] 0xc000004560:[]]
	t.Logf("len(g): %d, g: %v\n", len(g), g)
}
