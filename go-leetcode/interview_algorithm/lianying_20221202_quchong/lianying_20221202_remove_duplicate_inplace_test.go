package lianying_20221202_quchong

import (
	"reflect"
	"testing"
)

func Test_rep(t *testing.T) {
	type Example struct {
		Name         string
		Input        []int
		ExpectOutput []int
		ExpectLen    int
	}

	tests := []Example{
		{
			Name:         "用例-1",
			Input:        []int{1, 2, 2, 3, 4, 4, 4, 5},
			ExpectOutput: []int{1, 2, 3, 4, 5, 4, 4, 5},
			ExpectLen:    5,
		},
		{
			Name:         "用例-2",
			Input:        nil,
			ExpectOutput: nil,
			ExpectLen:    0,
		},
		{
			Name:         "用例-3",
			Input:        []int{1, 1, 1, 1, 1},
			ExpectOutput: []int{1, 1, 1, 1, 1},
			ExpectLen:    1,
		},
		{
			Name:         "用例-4",
			Input:        []int{1, 2, 3, 4, 5},
			ExpectOutput: []int{1, 2, 3, 4, 5},
			ExpectLen:    5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := removeDuplicateInplace(tt.Input)
			if !reflect.DeepEqual(tt.Input, tt.ExpectOutput) || got != tt.ExpectLen {
				t.Errorf("Input = %v, ExpectOutput = %v, got = %v,  ExpectLen = %v", tt.Input, tt.ExpectOutput, got, tt.ExpectLen)
			}
		})
	}
}
