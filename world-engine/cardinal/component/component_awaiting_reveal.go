package component

type AwaitingReveal struct {
	IsAvailable bool
}

func (AwaitingReveal) Name() string {
	return "AwaitingReveal"
}
