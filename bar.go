package chocolate

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type LayoutType int

const (
	LIST LayoutType = iota
	LINEAR
)

// ScalingType defines how the ChocolateBar will be scaled
// FIXED will be a fixed number of cells
// RELATIVE will grow or shrink on screen resizes relative to the set values
type ScalingType int

const (
	PARENT ScalingType = iota
	DYNAMIC
	FIXED
)

type scaling struct {
	t ScalingType
	v int
}

type Scaler interface {
	Get() (ScalingType, int)
	GetValue() int
	Set(ScalingType, int)
	Is(ScalingType) bool
	IsParent() bool
	IsDynamic() bool
	IsFixed() bool
}

type scaler struct {
	t ScalingType
	v int
}

func (s scaler) Get() (ScalingType, int) { return s.t, s.v }
func (s scaler) GetValue() int           { return s.v }

func (s *scaler) Set(t ScalingType, v int) {
	switch t {
	case DYNAMIC:
		s.t = t
		s.v = 0
	default:
		s.t = t
		s.v = v
		if v <= 0 {
			s.v = 1
		}
	}
}

func (s scaler) Is(t ScalingType) bool { return s.t == t }
func (s scaler) IsParent() bool        { return s.Is(PARENT) }
func (s scaler) IsDynamic() bool       { return s.Is(DYNAMIC) }
func (s scaler) IsFixed() bool         { return s.Is(FIXED) }

func NewScaler(t ScalingType, v int) Scaler {
	ret := &scaler{}
	ret.Set(t, v)
	return ret
}

func NewParentScaler(v int) Scaler {
	return NewScaler(PARENT, v)
}

func NewDynamicScaler() Scaler {
	return NewScaler(DYNAMIC, 0)
}

func NewFixedScaler(v int) Scaler {
	return NewScaler(FIXED, v)
}

type Scaling struct {
	X Scaler
	Y Scaler
}

type ChocolateBar struct {
	Scaling
	id string

	// bars in order for the layout
	bars []*ChocolateBar
	// backref to the parent bar
	// this is used to propagate the
	// dynamic sizing back to the parent
	// as well to also adjust depending
	// on the parent layout
	parent *ChocolateBar

	// layout parameters
	layoutType LayoutType

	// possible maximum content size
	// This is used to have a maximum for the content
	// the real size will be calculated during the
	// view rendering as it is the only possible
	// place to handle the scaling which depends on
	// possible dynmic content
	maxWidth  int
	maxHeight int

	// content size after calculation of the whole
	// layout
	width  int
	height int

	// model of the entry this can only be set, if
	// there are no sub bars and is the final leaf
	// of the whole tree and provides the real content
	model tea.Model

	// pre rendered view with maximum content sizes
	// this is used to get the correct sizes of the
	// models view to be used for dynamic scaling
	// and the following calculations
	preRendered   bool
	preView       string
	contentWidth  int
	contentHeight int

	// rendered
	view     string
	rendered bool

	// flavour
	flavour Flavour

	// flavourPrefs generation function
	// this can be used to override the default
	// flavour preferences
	FlavourPrefsFct func() FlavourPrefs

	// if the bar is hidden
	// hidden bars are removed from the layout
	// rendering and the space is used for the
	// other bars
	hidden bool

	// if this bar can be selected
	selectable bool

	// if this bar is selected
	selected bool
	// if this bar has input focus
	focus bool
}

func (b *ChocolateBar) defaultFlavourPrefs() FlavourPrefs {
	ret := NewFlavourPrefs()
	if len(b.bars) == 0 {
		ret = ret.BorderType(b.flavour.GetBorderType())
	}
	if b.selected {
		ret = ret.ForegroundBorder(FOREGROUND_HIGHLIGHT_PRIMARY)
	}
	if b.focus {
		ret = ret.Foreground(FOREGROUND_HIGHLIGHT_PRIMARY)
		ret = ret.Background(BACKGROUND_HIGHLIGHT_PRIMARY)
		ret = ret.ForegroundBorder(FOREGROUND_HIGHLIGHT_PRIMARY)
		// ret = ret.BackgroundBorder(BACKGROUND_HIGHLIGHT_PRIMARY)
	}

	return ret
}

func (b ChocolateBar) GetStyle() lipgloss.Style {
	return b.flavour.GetStyle(b.FlavourPrefsFct())
}

func (b *ChocolateBar) Select(v bool) {
	b.selected = v
}

func (b *ChocolateBar) Focus(v bool) {
	b.focus = v
}

func (b *ChocolateBar) Resize(w, h int) {
	// if there is a frame set for the bar
	// this has to be removed from the available
	// content size
	width := w - b.GetStyle().GetHorizontalFrameSize()
	height := h - b.GetStyle().GetVerticalFrameSize()

	// if this is a fixed scaling than we don't have
	// to calculate anything
	if b.X.IsFixed() {
		width = b.X.GetValue()
	}
	if b.Y.IsFixed() {
		height = b.X.GetValue()
	}

	if width <= 0 || height <= 0 {
		// TODO: error handling
		return
	}

	b.maxWidth = width
	b.maxHeight = height
	if b.model != nil {
		b.model, _ = b.model.Update(tea.WindowSizeMsg{Width: width, Height: height})
	} else {
		for _, c := range b.bars {
			c.Resize(width, height)
		}
	}
}

// pre render all models with their actual sizes
// this is a temporary task that has to be done
// so that it is possible to calculate the dynamic
// sizes
// TODO: Is there a better way to avoid calling models view?
func (b *ChocolateBar) preRender() {
	// skip hidden bars
	if b.hidden {
		return
	}

	if b.model != nil {
		if !b.preRendered {
			b.preView = b.model.View()
			b.contentWidth, b.contentHeight = lipgloss.Size(b.preView)

			if b.contentWidth > b.maxWidth || b.contentHeight > b.maxHeight {
				// TODO: error handling
				return
			}
			b.preRendered = true

			if b.parent != nil {
				t, v := b.X.Get()
				switch t {
				case DYNAMIC:
					b.parent.contentWidth += b.contentWidth + b.GetStyle().GetHorizontalFrameSize()
					b.width = b.contentWidth //+ b.getStyle().GetHorizontalFrameSize()
				case FIXED:
					b.parent.contentWidth += v + b.GetStyle().GetHorizontalFrameSize()
					b.width = v
				}
				t, v = b.Y.Get()
				switch t {
				case DYNAMIC:
					b.parent.contentHeight += b.contentHeight + b.GetStyle().GetVerticalFrameSize()
					b.height = b.contentHeight //+ b.getStyle().GetVerticalFrameSize()
				case FIXED:
					b.parent.contentHeight += v + b.GetStyle().GetVerticalFrameSize()
					b.height = v
				}
			}
		}
		return
	}

	// must be a bar without model
	// so go recursive to generate
	// all preViews of models
	for _, c := range b.bars {
		c.preRender()
	}

	// all sub bars of this model are now pre rendered
	// we can build up the used sizes of the fixed
	// and dynamic sub bars
	if b.parent != nil {
		t, v := b.X.Get()
		switch t {
		case DYNAMIC:
			b.parent.contentWidth += b.contentWidth + b.GetStyle().GetHorizontalFrameSize()
		case FIXED:
			b.parent.contentHeight += v + b.GetStyle().GetHorizontalFrameSize()
		}
		t, v = b.Y.Get()
		switch t {
		case DYNAMIC:
			b.parent.contentHeight += b.contentHeight + b.GetStyle().GetVerticalFrameSize()
		case FIXED:
			b.parent.contentHeight += v + b.GetStyle().GetVerticalFrameSize()
		}
	}
}

func (b *ChocolateBar) recalcSizes() {
	// skip hidden bars
	if b.hidden {
		return
	}

	// already done so just return
	if b.preRendered {
		return
	}

	switch b.layoutType {
	case LIST:
		b.recalcVerticalSizes()
	case LINEAR:
		b.recalcHorizontalSizes()
	}
}

func (b *ChocolateBar) recalcVerticalSizes() {
	// after pre render all leafs with models
	// this must be a bar holding subs
	// so go recursive till we reach the last
	// layers
	for _, c := range b.bars {
		c.recalcSizes()
	}

	// go over again and start calculation
	totalParts := 0
	totalParents := 0
	for _, c := range b.bars {
		if c.Y.IsParent() && !c.hidden {
			totalParts += c.Y.GetValue()
			totalParents++
		}
	}

	if totalParts > 0 {
		partSize := (b.maxHeight - b.contentHeight) / totalParts
		partLast := (b.maxHeight - b.contentHeight) % totalParts

		for _, c := range b.bars {
			if c.Y.IsParent() && !c.hidden {
				totalParents--
				height := c.Y.GetValue() * partSize
				if totalParents == 0 {
					height += partLast
				}
				c.height = height - c.GetStyle().GetVerticalFrameSize()
				b.contentHeight += c.height
			}
		}
	}

	b.height = b.contentHeight
	if b.height > b.maxHeight {
		// TODO: error handling
		return
	}

	if b.parent != nil {
		b.parent.contentHeight += b.height
	}
	b.preRendered = true
}

func (b *ChocolateBar) recalcHorizontalSizes() {
	// after pre render all leafs with models
	// this must be a bar holding subs
	// so go recursive till we reach the last
	// layers
	for _, c := range b.bars {
		c.recalcSizes()
	}

	// go over again and start calculation
	totalParts := 0
	totalParents := 0
	for _, c := range b.bars {
		if c.X.IsParent() && !c.hidden {
			totalParts += c.X.GetValue()
			totalParents++
		}
	}

	if totalParts > 0 {
		partSize := (b.maxWidth - b.contentWidth) / totalParts
		partLast := (b.maxWidth - b.contentWidth) % totalParts

		for _, c := range b.bars {
			if c.X.IsParent() && !c.hidden {
				totalParents--
				width := c.X.GetValue() * partSize
				if totalParents == 0 {
					width += partLast
				}
				c.width = width - c.GetStyle().GetHorizontalFrameSize()
				b.contentWidth += c.width
			}
		}
	}

	b.width = b.contentWidth
	if b.width > b.maxWidth {
		// TODO: error handling
		return
	}

	if b.parent != nil {
		b.parent.contentWidth += b.width
	}
	b.preRendered = true
}

func (b *ChocolateBar) finalizeSizing() {
	// skip hidden bars
	if b.hidden {
		return
	}

	for _, c := range b.bars {
		c.finalizeSizing()
	}

	width := b.maxWidth
	height := b.maxHeight
	if b.parent != nil {
		width = b.parent.width
		height = b.parent.height
		if width <= 0 {
			width = b.parent.maxWidth
		}
		if height <= 0 {
			height = b.parent.maxHeight
		}
	}
	if b.width <= 0 {
		b.width = width - b.GetStyle().GetHorizontalFrameSize()
	}
	if b.height <= 0 {
		b.height = height - b.GetStyle().GetVerticalFrameSize()
	}

	if b.model != nil {
		b.model, _ = b.model.Update(tea.WindowSizeMsg{Width: b.width, Height: b.height})
	}
}

func (b *ChocolateBar) render() {
	b.preRender()
	b.recalcSizes()
	b.finalizeSizing()

	// skip hidden bars
	if b.hidden {
		return
	}

	if b.model != nil {
		b.view = b.GetStyle().
			Width(b.width).
			Height(b.height).
			Render(b.model.View())
		b.rendered = true
		return
	}

	for _, c := range b.bars {
		c.render()
	}
}

func (b *ChocolateBar) joinBars() {
	// skip hidden bars
	if b.hidden {
		return
	}

	var bars []string
	if !b.rendered {
		b.rendered = true
		for _, c := range b.bars {
			c.joinBars()
			if c.hidden {
				continue
			}
			w, h := lipgloss.Size(c.view)
			switch b.layoutType {
			case LIST:
				if w < b.width {
					s := b.GetStyle().
						BorderTop(false).
						BorderBottom(false).
						BorderLeft(false).
						BorderRight(false).
						Width(b.width + b.GetStyle().GetHorizontalFrameSize())
					bars = append(bars, s.Render(c.view))
				} else {
					bars = append(bars, c.view)
				}
			case LINEAR:
				if h < b.height {
					s := b.GetStyle().
						BorderTop(false).
						BorderBottom(false).
						BorderLeft(false).
						BorderRight(false).
						Height(b.height - b.GetStyle().GetVerticalFrameSize())
					bars = append(bars, s.Render(c.view))
				} else {
					bars = append(bars, c.view)
				}
			default:
				bars = append(bars, c.view)
			}
		}
		switch b.layoutType {
		case LIST:
			b.view = b.GetStyle().
				Render(lipgloss.JoinVertical(0, bars...))
		case LINEAR:
			b.view = b.GetStyle().
				Render(lipgloss.JoinHorizontal(0, bars...))
		}
	}
}

func (b *ChocolateBar) resetRender() {
	for _, c := range b.bars {
		c.resetRender()
	}

	b.preRendered = false
	b.contentHeight = 0
	b.contentWidth = 0
	b.width = 0
	b.height = 0
	b.preView = ""
	b.rendered = false
	b.view = ""
}

func (b *ChocolateBar) Render() string {
	defer b.resetRender()

	b.resetRender()
	b.render()
	b.joinBars()
	return b.view
}

func (b *ChocolateBar) Hide(v bool) {
	b.hidden = v
}

func (b *ChocolateBar) HandleUpdate(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if b.model != nil {
		b.model, cmd = b.model.Update(msg)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

type chocolateBarOptions func(*ChocolateBar)

func WithLayout(v LayoutType) func(*ChocolateBar) {
	return func(b *ChocolateBar) {
		b.layoutType = v
	}
}

func WithID(v string) func(*ChocolateBar) {
	return func(b *ChocolateBar) {
		b.id = v
	}
}

func WithModel(v tea.Model) func(*ChocolateBar) {
	return func(b *ChocolateBar) {
		b.model = v
		b.bars = nil
	}
}

func WithBarFlavor(v Flavour) func(*ChocolateBar) {
	return func(b *ChocolateBar) {
		b.flavour = v
	}
}

func WithSelectable() func(*ChocolateBar) {
	return func(b *ChocolateBar) {
		b.selectable = true
	}
}

func NewChocolateBar(bars []*ChocolateBar, opts ...chocolateBarOptions) *ChocolateBar {
	ret := &ChocolateBar{
		id:            uuid.NewString(),
		bars:          bars,
		layoutType:    LIST,
		flavour:       NewFlavour(),
		preRendered:   false,
		preView:       "",
		rendered:      false,
		view:          "",
		width:         0,
		height:        0,
		contentWidth:  0,
		contentHeight: 0,
		hidden:        false,
		selectable:    false,
		focus:         false,
	}
	ret.X = NewParentScaler(1)
	ret.Y = NewParentScaler(1)

	for _, c := range bars {
		c.parent = ret
	}

	ret.FlavourPrefsFct = ret.defaultFlavourPrefs

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
