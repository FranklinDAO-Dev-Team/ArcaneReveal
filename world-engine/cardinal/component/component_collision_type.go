package component

type CollisionType struct {
	Type string `json:"type"`
}

func (CollisionType) Name() string {
	return "CollisionType"
}
