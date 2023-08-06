package designmode_1_1_simple_factory

import (
	"fmt"
	"testing"
)

func Test_main(t *testing.T) {
	ak47 := NewGun("Ak47")
	musket := NewGun("Musket")

	fmt.Println(ak47.Name())
	fmt.Println(musket.Name())
}
