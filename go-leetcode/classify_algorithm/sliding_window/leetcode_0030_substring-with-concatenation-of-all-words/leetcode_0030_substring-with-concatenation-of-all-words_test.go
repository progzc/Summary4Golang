package leetcode_0030_substring_with_concatenation_of_all_words

import (
	"fmt"
	"testing"
)

func permutation(words []string) []string {
	var (
		dfs   func(pos int, cur string)
		has   = map[string]bool{}
		perms []string
	)
	dfs = func(i int, cur string) {
		if i == len(words) {
			if _, ok := has[cur]; !ok {
				perms = append(perms, cur)
				has[cur] = true
			}
			return
		}
		for j := i; j < len(words); j++ {
			swap(words, i, j)
			temp := cur + words[i]
			dfs(i+1, temp)
			swap(words, i, j)
		}
		return
	}
	dfs(0, "")
	return perms
}

func Test_permutation(t *testing.T) {
	words := []string{"ab", "cd", "ef"}
	fmt.Println(permutation(words))

	words = []string{"word", "good", "best", "good"}
	fmt.Println(permutation(words))

	words = []string{"foo", "bar"}
	fmt.Println(permutation(words))
}

func Test_findSubstring(t *testing.T) {
	s := "wordgoodgoodgoodbestword"
	words := []string{"word", "good", "best", "good"}
	fmt.Printf("findSubstring: %v\n", findSubstring(s, words))
}
