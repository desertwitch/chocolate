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
	NONE
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
	XAXIS ScalingAxis = iota
	YAXIS
)

type BarStyler interface {
	GetStyle() lipgloss.Style
}

type BarSelector interface {
	GetID() string
	SetID(id string)
	IsHidden() bool
	IsSelectable() bool
	IsFocusable() bool
	isOverlay() bool
}

type BarScaler interface {
	GetScaler(axis ScalingAxis) (ScalingType, int)
	SetScaler(axis ScalingAxis, scalingType ScalingType, value int)
}

func IsXFixed(barScaler BarScaler) bool {
	t, _ := barScaler.GetScaler(XAXIS)
	return t == FIXED
}

func IsXParent(barScaler BarScaler) bool {
	t, _ := barScaler.GetScaler(XAXIS)
	return t == PARENT
}

func GetXValue(barScaler BarScaler) int {
	_, v := barScaler.GetScaler(XAXIS)
	return v
}

func IsYFixed(barScaler BarScaler) bool {
	t, _ := barScaler.GetScaler(YAXIS)
	return t == FIXED
}

func IsYParent(barScaler BarScaler) bool {
	t, _ := barScaler.GetScaler(YAXIS)
	return t == PARENT
}

func GetYValue(barScaler BarScaler) int {
	_, v := barScaler.GetScaler(YAXIS)
	return v
}

func GetScalerValue(axis ScalingAxis, barScaler BarScaler) int {
	_, v := barScaler.GetScaler(axis)
	return v
}

func SetScalerValue(axis ScalingAxis, barScaler BarScaler, value int) {
	t, _ := barScaler.GetScaler(axis)
	barScaler.SetScaler(axis, t, value)
}

type BarController interface {
	Hide(value bool)
	Selectable(value bool)
	Focusable(value bool)
	setOverlay()
}

type PlacementType int

const (
	START PlacementType = iota
	CENTER
	END
	POSITION
)

type BarPlacer interface {
	GetPlacement(axis ScalingAxis) (PlacementType, int)
	SetPlacement(axis ScalingAxis, placementType PlacementType, value int)
}

func calcCenterPos(sizeView, sizeParentView int) int {
	pos := (sizeParentView / 2) - (sizeView / 2)
	return pos
}

func calcEndPos(sizeView, sizeParentView int) int {
	pos := sizeParentView - sizeView
	return pos
}

func calcPlacerXPos(bar BarPlacer, view, parentView string) int {
	p, v := bar.GetPlacement(XAXIS)
	vx, _ := lipgloss.Size(view)
	px, _ := lipgloss.Size(parentView)

	switch p {
	case START:
		return 0
	case POSITION:
		return v
	case CENTER:
		return calcCenterPos(vx, px)
	case END:
		return calcEndPos(vx, px)
	}
	return 0
}

func calcPlacerYPos(bar BarPlacer, view, parentView string) int {
	p, v := bar.GetPlacement(YAXIS)
	_, vy := lipgloss.Size(view)
	_, py := lipgloss.Size(parentView)

	switch p {
	case START:
		return 0
	case POSITION:
		return v
	case CENTER:
		return calcCenterPos(vy, py)
	case END:
		return calcEndPos(vy, py)
	}
	return 0
}

type BarOverlay interface {
	GetZindex() int
	SetZindex(index int)
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
	BarLayouter
	BarSizer
	BarMaxSizer
}

type BarChild interface {
	BarScaler
	BarPlacer
	BarSelector
	BarSizer
	BarContentSizer
	GetView() string
}

type BarLayouter interface {
	GetLayout() LayoutType
	SetLayout(LayoutType)
}

type BarModeler interface {
	GetModel() tea.Model
	SelectModel(string)
}

type ChocolateSelector interface {
	IsSelected(barSelector BarSelector) bool
	IsRoot(barSelector BarSelector) bool
	IsFocused(barSelector BarSelector) bool
	GetParent(barSelector BarSelector) BarParent
	GetChildren(barSelector BarSelector) []BarChild
	Select(bar ChocolateBar)
	ForceSelect(bar ChocolateBar)
	Focus(barSelector BarSelector)
	GetByID(id string) ChocolateBar
}

type BarUpdater interface {
	HandleUpdate(msg tea.Msg) tea.Cmd
}

type ChocolateBar interface {
	BarScaler
	BarPlacer
	BarSelector
	BarController
	BarRenderer
	BarLayouter
	BarModeler
	BarSizer
	BarMaxSizer
	BarContentSizer
	ChocolateSelector
	BarUpdater
	setBarStyler(barStyler BarStyler)
	setBarPlacer(barPlacer BarPlacer)
	setBarScaler(barScaler BarScaler)
	setBarSelector(barSelector BarSelector)
	setBarController(barController BarController)
	setBarChocolate(barChocolate ChocolateSelector)
	setStyleCustomizeHandler(styler BaseBarStyleCustomizeHanleFct)
}

func SetLayoutSize(nChocolateBar ChocolateBar, value int) {
	p := nChocolateBar.GetParent(nChocolateBar)

	var axis ScalingAxis
	switch p.GetLayout() {
	case LINEAR:
		axis = XAXIS
	case LIST:
		axis = YAXIS
	}

	nt, _ := nChocolateBar.GetScaler(axis)
	if nt == DYNAMIC {
		return
	}
	SetScalerValue(axis, nChocolateBar, value)
}

func IncLayoutSize(nChocolateBar ChocolateBar) int {
	p := nChocolateBar.GetParent(nChocolateBar)

	var axis ScalingAxis
	switch p.GetLayout() {
	case LINEAR:
		axis = XAXIS
	case LIST:
		axis = YAXIS
	}

	nt, value := nChocolateBar.GetScaler(axis)
	if nt == DYNAMIC {
		return 0
	}
	value++
	SetScalerValue(axis, nChocolateBar, value)
	return value
}

func DecLayoutSize(nChocolateBar ChocolateBar) int {
	p := nChocolateBar.GetParent(nChocolateBar)

	var axis ScalingAxis
	switch p.GetLayout() {
	case LINEAR:
		axis = XAXIS
	case LIST:
		axis = YAXIS
	}

	nt, value := nChocolateBar.GetScaler(axis)
	if nt == DYNAMIC {
		return 0
	}
	value--
	SetScalerValue(axis, nChocolateBar, value)
	return value
}
