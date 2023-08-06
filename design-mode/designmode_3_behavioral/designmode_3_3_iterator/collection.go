package designmode_3_3_iterator

type Collection interface {
	createIterator() Iterator
}

type UserCollection struct {
	users []*User
}

type User struct {
	name string
	age  int
}

func (u *UserCollection) createIterator() Iterator {
	return &UserIterator{
		users: u.users,
	}
}
