package designmode_1_1_simple_factory

// IGunFactory 枪支工厂接口
type IGunFactory interface {
	GetGun() IGun
}

// Ak47Factory Ak47工厂
type Ak47Factory struct{}

func (f *Ak47Factory) GetGun() IGun {
	return &Ak47{}
}

// MusketFactory Musket工厂
type MusketFactory struct{}

func (f *MusketFactory) GetGun() IGun {
	return &Musket{}
}

// NewGunFactory 工厂方法
func NewGunFactory(s string) IGunFactory {
	switch s {
	case "Ak47":
		return &Ak47Factory{}
	case "Musket":
		return &MusketFactory{}
	}
	return nil
}
