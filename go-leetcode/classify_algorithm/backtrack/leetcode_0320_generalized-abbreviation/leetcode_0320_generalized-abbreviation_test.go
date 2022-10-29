package leetcode_0320_generalized_abbreviation

import (
	"fmt"
	"testing"
)

func Test_generateAbbreviations(t *testing.T) {
	word := "word"
	// [4 3d 2r1 2rd 1o2 1o1d 1or1 1ord w3 w2d w1r1 w1rd wo2 wo1d wor1 word]
	fmt.Println(generateAbbreviations(word))
}
