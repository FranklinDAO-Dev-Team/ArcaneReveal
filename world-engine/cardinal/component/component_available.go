package component

type Available struct {
	IsAvailable bool
}

func (Available) Name() string {
	return "IsAvailable"
}

func (a *Available) SetIsAvailable(b bool) {
	a.IsAvailable = b
}

func (a *Available) FlipIsAvailable() {
	a.IsAvailable = !a.IsAvailable
}
