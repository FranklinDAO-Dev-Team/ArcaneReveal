package component

type Health struct {
	MaxHealth  int
	CurrHealth int
}

func (Health) Name() string {
	return "Health"
}
