package leetcode_1258_synonymous_sentences

import (
	"fmt"
	"testing"
)

func Test_generateSentences(t *testing.T) {
	synonyms := [][]string{{"happy", "joy"}, {"sad", "sorrow"}, {"joy", "cheerful"}}
	text := "I am happy today but was sad yesterday"
	results := generateSentences(synonyms, text)
	for _, result := range results {
		fmt.Println(result)
	}
}

func Test_generateSentences2(t *testing.T) {
	synonyms := [][]string{{"a", "b"}, {"b", "c"}, {"d", "e"}, {"c", "d"}}
	text := "a b"
	results := generateSentences(synonyms, text)
	for _, result := range results {
		fmt.Println(result)
	}
}
