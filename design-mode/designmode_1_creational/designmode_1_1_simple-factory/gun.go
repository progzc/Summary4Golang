package designmode_1_1_simple_factory

// IGun 枪支接口
type IGun interface {
	Name() string
}

// Ak47 Ak47枪支
type Ak47 struct{}

func (a *Ak47) Name() string {
	return "this is Ak47"
}

// Musket Musket枪支
type Musket struct{}

func (m *Musket) Name() string {
	return "this is Musket"
}

// NewGun 简单工厂
func NewGun(s string) IGun {
	switch s {
	case "Ak47":
		return &Ak47{}
	case "Musket":
		return &Musket{}
	}
	return nil
}
