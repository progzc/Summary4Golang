package designmode_3_1_chain

type Department interface {
	execute(*Patient)
	setNext(Department)
}

type Patient struct {
	name              string
	registrationDone  bool
	doctorCheckUpDone bool
	medicineDone      bool
	paymentDone       bool
}
