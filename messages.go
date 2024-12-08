package chocolate

import tea "github.com/charmbracelet/bubbletea"

type BarHideMsg struct {
	Id    string
	Value bool
}

type ModelChangeMsg struct {
	Id    string
	Model string
}

type ForceSelectMsg string

type SelectMsg string

type ErrorMsg error

func NewBarHideMsg(id string, v bool) tea.Cmd {
	return func() tea.Msg {
		return BarHideMsg{
			Id:    id,
			Value: v,
		}
	}
}

func NewModelChangeMsg(id, model string) tea.Cmd {
	return func() tea.Msg {
		return ModelChangeMsg{
			Id:    id,
			Model: model,
		}
	}
}

func NewForceSelectMsg(id string) tea.Cmd {
	return func() tea.Msg {
		return ForceSelectMsg(id)
	}
}
