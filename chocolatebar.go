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

type CheckAttributes int

const (
	CA_LAYOUT CheckAttributes = iota
	CA_X_SCALING
	CA_Y_SCALING
	CA_CAN_SELECT
	CA_CAN_FOCUS
	CA_INPUT_ON_SELECT
	CA_IS_MODELBAR
)

type checkAttributes struct {
	checks map[CheckAttributes]interface{}
}

func (c checkAttributes) CheckLayout(v LayoutType) checkAttributes {
	c.checks[CA_LAYOUT] = v
	return c
}

func (c checkAttributes) CheckXScaling(v ScalingType) checkAttributes {
	c.checks[CA_X_SCALING] = v
	return c
}

func (c checkAttributes) CheckYScaling(v ScalingType) checkAttributes {
	c.checks[CA_Y_SCALING] = v
	return c
}

func (c checkAttributes) CheckCanSelect() checkAttributes {
	c.checks[CA_CAN_SELECT] = true
	return c
}

func (c checkAttributes) CheckCanFocus() checkAttributes {
	c.checks[CA_CAN_FOCUS] = true
	return c
}

func (c checkAttributes) CheckInputOnSelect() checkAttributes {
	c.checks[CA_INPUT_ON_SELECT] = true
	return c
}

func (c checkAttributes) CheckIsModelBar() checkAttributes {
	c.checks[CA_IS_MODELBAR] = true
	return c
}

func newCheckAtrributes() checkAttributes {
	return checkAttributes{
		checks: make(map[CheckAttributes]interface{}),
	}
}

type ChocolateBar interface {
	Has(interface{}) bool
	GetID() string
	GetStyle() lipgloss.Style
	// IsRoot() bool
	// GetParent() CChocolateBar
	// CanSelect() bool
	// InputOnSelect() bool
	// SetChocolate(*Chocolate)
	// GetChocolate() *Chocolate
	// Resize(int, int)
	// GetScaling() Scaling
	// SetScaling(Scaling)
	// Render() string
	// GetBars() []ChocolateBar
	Hide(bool)
	// HandleUpdate(tea.Msg) tea.Cmd

	// new
	getChocolate() *NChocolate
	setChocolate(*NChocolate)

	setID(string)
	setXScaler(Scaler)
	setYScaler(Scaler)
	setSelectable()
	setInputOnSelect()

	getContentSize() (int, int)
	setContentSize(int, int)
	preRender(ChocolateBar) bool
}

type LayoutBar interface {
	ChocolateBar
	GetLayout() LayoutType
	SetLayout(LayoutType)
	getTotalParts() int
	setTotalParts(int)
}

type ModelBar interface {
	ChocolateBar
	GetModel() tea.Model
	SelectModel(string)
}
