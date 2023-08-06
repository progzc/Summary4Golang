package designmode_1_6_singleton

import (
	"fmt"
	"testing"
)

func Test_main(t *testing.T) {
	for i := 0; i < 30; i++ {
		go getInstance()
	}

	fmt.Scanln()
}
