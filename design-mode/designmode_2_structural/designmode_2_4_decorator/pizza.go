package designmode_2_4_decorator

type IPizza interface {
	getPrice() int
}

type VeggeMania struct {
}

func (p *VeggeMania) getPrice() int {
	return 15
}
