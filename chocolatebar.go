package chocolate

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// A LayoutType defines the base direction of the bar
type LayoutType int

// Layout types
const (
	LIST   LayoutType = iota // will define a vertical arranged layout
	LINEAR                   // will define a horizontal arranged layout
)

// A ScalingType defines how the ChocolateBar will be scaled
type ScalingType int

// Scaling types
const (
	PARENT  ScalingType = iota // will fill up the available size
	DYNAMIC                    // will grow as big as the content is
	FIXED                      // is a fixed size
)

type ScalingAxis int

const (
	X ScalingAxis = iota
	Y
)

type BarStyler interface {
	GetStyle() lipgloss.Style
}

type BarSelector interface {
	GetID() string
	SetID(id string)
	IsHidden() bool
	IsSelectable() bool
}

type BarScaler interface {
	GetScaler(axis ScalingAxis) (ScalingType, int)
	SetScaler(axis ScalingAxis, scalingType ScalingType, value int)
}

func IsXFixed(barScaler BarScaler) bool {
	t, _ := barScaler.GetScaler(X)
	return t == FIXED
}

func IsXParent(barScaler BarScaler) bool {
	t, _ := barScaler.GetScaler(X)
	return t == PARENT
}

func GetXValue(barScaler BarScaler) int {
	_, v := barScaler.GetScaler(X)
	return v
}

func IsYFixed(barScaler BarScaler) bool {
	t, _ := barScaler.GetScaler(Y)
	return t == FIXED
}

func IsYParent(barScaler BarScaler) bool {
	t, _ := barScaler.GetScaler(Y)
	return t == PARENT
}

func GetYValue(barScaler BarScaler) int {
	_, v := barScaler.GetScaler(Y)
	return v
}

type BarController interface {
	Hide(value bool)
	Selectable(value bool)
}

type BarRenderer interface {
	Resize(width, height int)
	PreRender() bool
	Render()
	GetView() string
}

type BarSizer interface {
	GetSize() (width, height int)
	SetSize(width, height int)
}

func SetWidth(barSizer BarSizer, width int) {
	barSizer.SetSize(width, -1)
}

func SetHeight(barSizer BarSizer, height int) {
	barSizer.SetSize(-1, height)
}

type BarMaxSizer interface {
	GetMaxSize() (width, height int)
}

type BarContentSizer interface {
	GetContentSize() (width, height int)
}

type BarParent interface {
	BarSizer
	BarMaxSizer
}

type BarChild interface {
	BarScaler
	BarSelector
	BarSizer
	BarContentSizer
	GetView() string
}

type BarLayout interface {
	GetLayout() LayoutType
	SetLayout(LayoutType)
}

type BarModel interface {
	GetModel() tea.Model
	SelectModel(string)
}

type ChocolateSelector interface {
	IsSelected(barSelector BarSelector) bool
	IsRoot(barSelector BarSelector) bool
	IsFocused(barSelector BarSelector) bool
	GetParent(barSelector BarSelector) BarParent
	GetChildren(barSelector BarSelector) []BarChild
}

type NChocolateBar interface {
	BarScaler
	BarSelector
	BarRenderer
	BarSizer
	BarMaxSizer
	BarContentSizer
	SetBarChocolate(barChocolate ChocolateSelector)
}
