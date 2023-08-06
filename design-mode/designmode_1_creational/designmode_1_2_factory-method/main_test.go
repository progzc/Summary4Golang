package designmode_1_1_simple_factory

import (
	"fmt"
	"testing"
)

func Test_main(t *testing.T) {
	ak47Factory := NewGunFactory("Ak47")
	ak47 := ak47Factory.GetGun()

	musketFactory := NewGunFactory("Musket")
	musket := musketFactory.GetGun()

	fmt.Println(ak47.Name())
	fmt.Println(musket.Name())
}
