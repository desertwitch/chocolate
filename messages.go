package chocolate

type BarHideMsg struct {
	Id    string
	Value bool
}

type ModelChangeMsg struct {
	Id    string
	Model string
}

type ForceSelectMsg string

type ErrorMsg error
