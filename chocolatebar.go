package chocolate

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CChocolateBar interface {
	GetID() string
	IsRoot() bool
	GetParent() CChocolateBar
	GetLayout() LayoutType
	CanFocus() bool
	CanSelect() bool
	GetSelectables() []string
	InputOnSelect() bool
	SetChocolate(*Chocolate)
	GetChocolate() *Chocolate
	Resize(int, int)
	GetStyle() lipgloss.Style
	GetScaling() Scaling
	SetScaling(Scaling)
	GetModel() tea.Model
	SelectModel(string)
	Render() string
	GetBars() []CChocolateBar
	Hide(bool)
	HandleUpdate(tea.Msg) tea.Cmd
}

type ChocolateBar interface {
	GetID() string
	IsRoot() bool
	GetParent() CChocolateBar
	CanSelect() bool
	InputOnSelect() bool
	SetChocolate(*Chocolate)
	GetChocolate() *Chocolate
	Resize(int, int)
	GetStyle() lipgloss.Style
	GetScaling() Scaling
	SetScaling(Scaling)
	Render() string
	GetBars() []ChocolateBar
	Hide(bool)
	HandleUpdate(tea.Msg) tea.Cmd
}

type RootBar interface {
	HandleUpdate(tea.Msg) tea.Cmd
}

type LayoutBar interface {
	GetLayout() LayoutType
	SetLayout(LayoutType)
	GetBars() []ChocolateBar
}

type ModelBar interface {
	GetModel() tea.Model
	SelectModel(string)
	CanFocus() bool
}
