package designmode_3_6_observer

import "fmt"

// Observer 观察者接口
type Observer interface {
	update(string)
	getID() string
}

type Customer struct {
	id string
}

func (c *Customer) update(itemName string) {
	fmt.Printf("Sending email to customer %s for item %s\n", c.id, itemName)
}

func (c *Customer) getID() string {
	return c.id
}
