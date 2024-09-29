package common_mistake

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	var s1 []string         // good
	s2 := []string{}        // bad
	s3 := make([]string, 0) // bad
	b1, _ := json.Marshal(s1)
	b2, _ := json.Marshal(s2)
	b3, _ := json.Marshal(s3)
	fmt.Println(s1, s1 == nil, string(b1))
	fmt.Println(s2, s2 == nil, string(b2))
	fmt.Println(s3, s3 == nil, string(b3))
}

func TestSlice2(t *testing.T) {
	ss := []string{"a", "b", "c"}
	for s := range ss {
		fmt.Println(s) // 0 1 2
	}

	for _, s := range ss {
		fmt.Println(s) // a b c
	}
}
