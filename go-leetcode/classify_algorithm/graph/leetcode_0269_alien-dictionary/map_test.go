package leetcode_0269_alien_dictionary

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	graph := map[byte][]byte{}
	graph['a'] = append(graph['a'], 'b')
	graph['a'] = append(graph['a'], 'c')
	graph['a'] = append(graph['a'], 'd')
	fmt.Println(graph)
}
