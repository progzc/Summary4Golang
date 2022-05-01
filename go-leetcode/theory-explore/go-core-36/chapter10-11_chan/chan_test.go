package chapter10

import (
	"fmt"
	"testing"
	"time"
)

type user struct {
	name string
	age  int8
}

var u = user{name: "Ankur", age: 25}
var g = &u

func modifyUser(pu *user) {
	fmt.Println("modifyUser Received Value", pu) // modifyUser Received Value &{Ankur Anand 100}
	pu.name = "Anand"
}

func printUser(u <-chan *user) {
	time.Sleep(2 * time.Second)
	fmt.Println("printUser GoRoutine called", <-u) // printUser GoRoutine called &{Ankur 25}
}

func Test_chan_1(t *testing.T) {
	c := make(chan *user, 5)
	c <- g
	fmt.Println(g) // &{Ankur 25}
	g = &user{name: "Ankur Anand", age: 100}
	go printUser(c)
	go modifyUser(g)
	time.Sleep(5 * time.Second)
	fmt.Println(g) // &{Anand 100}
}
