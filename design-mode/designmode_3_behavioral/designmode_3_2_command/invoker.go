package designmode_3_2_command

type Button struct {
	command Command
}

func (b *Button) press() {
	b.command.execute()
}
