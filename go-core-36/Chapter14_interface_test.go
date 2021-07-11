package go_core_36

import (
	"fmt"
	"testing"
)

type Pet interface {
	SetName(name string)
	Name() string
	Category() string
}

type Dog struct {
	name string // 名字。
}

func (dog *Dog) SetName(name string) {
	dog.name = name
}

func (dog Dog) Name() string {
	return dog.name
}

func (dog Dog) Category() string {
	return "dog"
}
func TestInterface(t *testing.T) {
	// 示例1。
	dog := Dog{"little pig"}
	_, ok := interface{}(dog).(Pet)
	fmt.Printf("Dog implements interface Pet: %v\n", ok)
	_, ok = interface{}(&dog).(Pet)
	fmt.Printf("*Dog implements interface Pet: %v\n", ok)
	fmt.Println()

	// 示例2。
	var pet Pet = &dog
	fmt.Printf("This pet is a %s, the name is %q.\n",
		pet.Category(), pet.Name())
}

func TestInterface2(t *testing.T) {
	var dog1 *Dog
	fmt.Println("The first dog is nil. [wrap1]")
	dog2 := dog1
	fmt.Println("The second dog is nil. [wrap1]")
	var pet Pet = dog2
	if pet == nil {
		fmt.Println("The pet is nil. [wrap1]")
	} else {
		fmt.Println("The pet is not nil. [wrap1]")
	}
}
