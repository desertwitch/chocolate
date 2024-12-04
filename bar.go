package chocolate

import (
	"log"

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

type Scaling interface {
	SetParent(int)
	SetDynamic()
	SetFixed(int)
	GetScaling() (ScalingType, int)
	GetDynamic() (bool, int)
	GetParent() (bool, int)
	GetFixed() (bool, int)
	IsParent() bool
	IsDynamic() bool
	IsFixed() bool
}

type scaling struct {
	t ScalingType
	v int
}

func (s *scaling) SetParent(v int) {
	s.t = PARENT
	s.v = v
}

func (s *scaling) SetDynamic() {
	s.t = DYNAMIC
	s.v = 0
}

func (s *scaling) SetFixed(v int) {
	s.t = FIXED
	s.v = v
}

func (s scaling) GetScaling() (ScalingType, int) {
	return s.t, s.v
}

func (s scaling) GetDynamic() (bool, int) {
	if s.t == DYNAMIC {
		return true, s.v
	}
	return false, 0
}

func (s scaling) GetParent() (bool, int) {
	if s.t == PARENT {
		return true, s.v
	}
	return false, 0
}

func (s scaling) GetFixed() (bool, int) {
	if s.t == FIXED {
		return true, s.v
	}
	return false, 0
}

func (s scaling) IsParent() bool {
	return s.t == PARENT
}

func (s scaling) IsDynamic() bool {
	return s.t == DYNAMIC
}

func (s scaling) IsFixed() bool {
	return s.t == FIXED
}

type Bar struct {
	scaling
	id string

	// bars in order for the layout
	bars []*Bar

	// layout parameters
	layoutType LayoutType

	// possible maximum content size
	// This is used to have a maximum for the content
	// the real size will be calculated during the
	// view rendering as it is the only possible
	// place to handle the scaling which depends on
	// possible dynmic content
	width  int
	height int

	// fixed sizing already used by dynamic and fixed
	// scaling bars
	fixedSize int

	// for dynamic scaling this is used to hold
	// the actual view from the model, which will
	// be used for the size calculation of the
	// bar
	rendered bool
	view     string

	// flavour
	flavour Flavour

	// flavourPrefs generation function
	// this can be used to override the default
	// flavour preferences
	flavourPrefsFct func() FlavourPrefs
}

func (b *Bar) defaultFlavourPrefs() FlavourPrefs {
	ret := NewFlavourPrefs()
	if len(b.bars) == 0 {
		ret = ret.BorderType(b.flavour.GetBorderType())
	}

	return ret
}

func (b Bar) getStyle() lipgloss.Style {
	return b.flavour.GetStyle(b.flavourPrefsFct())
}

func (b *Bar) Resize(size tea.WindowSizeMsg, models map[string]tea.Model, parent *Bar) {
	// calculate available size for bars
	w, h := b.getStyle().GetFrameSize()
	size.Width -= w
	size.Height -= h

	b.width = size.Width
	b.height = size.Height
	log.Printf("Bar: %s w=%d h=%d\n", b.id, b.width, b.height)

	// check size
	if size.Width <= 0 || size.Height <= 0 {
		// TODO: Error handling
		return
	}

	if models[b.id] != nil {
		// set layout to parent if model
		b.layoutType = parent.layoutType
	}

	// handle fixed bars as they're not changing
	// sizes dynamic
	if ok, v := b.GetFixed(); ok {
		switch b.layoutType {
		case LIST:
			if v > size.Height {
				// TODO Error handling
				return
			}
			size.Height = v
		case LINEAR:
			if v > size.Width {
				// TODO Error handling
				return
			}
			size.Width = v
		}
		b.width = size.Width
		b.height = size.Height
	}

	// if there is a model set then update
	// it's size and return
	if models[b.id] != nil {
		log.Printf("Model: %s w=%d h=%d\n", b.id, size.Width, size.Height)
		models[b.id].Update(size)
		return
	}

	// set the possible maximum sizes of the bars
	for _, c := range b.bars {
		c.Resize(size, models, b)
	}
}

func (b *Bar) resetRender() {
	if b.IsDynamic() {
		b.v = 0
	}

	b.rendered = false
	b.view = ""
	b.fixedSize = 0

	for _, c := range b.bars {
		c.resetRender()
	}
}

func (b *Bar) renderDynamic(models map[string]tea.Model, layout LayoutType) {
	if models[b.id] != nil {
		if b.IsDynamic() && !b.rendered {
			log.Printf("Dynamic: %s\n", b.id)
			view := models[b.id].View()
			// update sizes by real c.ontent
			width := lipgloss.Width(view)   // - b.flavour.GetHorizontalFrameSize()
			height := lipgloss.Height(view) // - b.flavour.GetVerticalFrameSize()
			if width > b.width || height > b.height {
				// TODO: error handling
				return
			}
			b.rendered = true
			b.width = width
			b.height = height

			b.view = b.getStyle().
				Width(b.width).
				Height(b.height).
				Render(models[b.id].View())
			log.Printf("Dynamic: %s %d %d\n", b.id, b.width, b.height)

			switch layout {
			case LIST:
				b.v = height + b.getStyle().GetHorizontalFrameSize()
			case LINEAR:
				b.v = width + b.getStyle().GetVerticalFrameSize()
			}
		}
	}

	for _, c := range b.bars {
		c.renderDynamic(models, b.layoutType)
	}
	for _, c := range b.bars {
		if c.IsDynamic() || c.IsFixed() {
			switch b.layoutType {
			case LIST:
				log.Printf("Fixed: %s %d %d\n", c.id, c.width, c.height)
				b.fixedSize += c.height + c.getStyle().GetVerticalFrameSize()
			case LINEAR:
				log.Printf("Fixed: %s %d %d\n", c.id, c.width, c.height)
				b.fixedSize += c.width + c.getStyle().GetHorizontalFrameSize()
			}
		}
	}
}

func (b *Bar) render(models map[string]tea.Model) bool {
	// pre render all dynamic models and set
	// sizes
	b.renderDynamic(models, b.layoutType)
	b.resize(models)

	if models[b.id] != nil && !b.rendered {
		log.Printf("%s: w=%d h=%d\n", b.id, b.width, b.height)
		// render the model as this is now real content
		// the size calculation has to be already done
		// dynamic models already rendered via preRender
		b.view = b.getStyle().
			Width(b.width).
			Height(b.height).
			Render(models[b.id].View())
		b.rendered = true
		return true
	}

	switch b.layoutType {
	case LIST:
		b.renderVertical(models)
	case LINEAR:
		b.renderHorizontal(models)
	}

	return false
}

func (b *Bar) resize(models map[string]tea.Model) {
	var check int
	var sizeTotal int
	var sizeMsgFct func(v int) tea.WindowSizeMsg

	switch b.layoutType {
	case LIST:
		check = b.height     //- b.flavour.GetVerticalFrameSize()
		sizeTotal = b.height //- b.flavour.GetVerticalFrameSize()
		sizeMsgFct = func(v int) tea.WindowSizeMsg {
			return tea.WindowSizeMsg{
				Width:  b.width, // - b.flavour.GetHorizontalFrameSize(),
				Height: v,
			}
		}
	case LINEAR:
		check = b.width     //- b.flavour.GetHorizontalFrameSize()
		sizeTotal = b.width //- b.flavour.GetHorizontalFrameSize()
		sizeMsgFct = func(v int) tea.WindowSizeMsg {
			return tea.WindowSizeMsg{
				Width:  v,
				Height: b.height, // - b.flavour.GetVerticalFrameSize(),
			}
		}
	}

	// check sizes
	if b.fixedSize >= check {
		log.Printf("ID: %s fixed=%d\n", b.id, b.fixedSize)
		// TODO: error handling
		return
	}

	totalParts := 0
	totalParents := 0
	for _, c := range b.bars {
		if ok, v := c.GetParent(); ok {
			totalParts += v
			totalParents++
		}
	}
	log.Printf("ID: %s totalParts=%d totalParents=%d fixed=%d\n", b.id, totalParts, totalParents, b.fixedSize)

	if totalParts > 0 {
		partSize := (sizeTotal - b.fixedSize) / totalParts
		partLast := (sizeTotal - b.fixedSize) % totalParts

		parentNum := 0
		for i, c := range b.bars {
			if ok, v := c.GetParent(); ok {
				parentNum++
				size := v * partSize
				if parentNum == totalParents {
					size += partLast
				}
				sizeMsg := sizeMsgFct(size)
				// sizeMsg.Width -= c.flavour.GetHorizontalFrameSize()
				// sizeMsg.Height -= c.flavour.GetVerticalFrameSize()
				log.Printf("Resize\n")
				log.Printf("%s: w=%d h=%d\n", c.id, sizeMsg.Width, sizeMsg.Height)
				c.Resize(sizeMsgFct(size), models, b)
				c.resize(models)
				// b.bars[i] = c
				log.Printf("%s: w=%d h=%d\n", b.bars[i].id, b.bars[i].width, b.bars[i].height)
			} else {
				c.resize(models)
			}
		}
	}
}

func (b *Bar) renderVertical(models map[string]tea.Model) {
	if len(b.bars) == 0 {
		// TODO: error handling
		return
	}

	totalHeight := 0
	for _, c := range b.bars {
		c.render(models)
		totalHeight += c.height
	}

	if totalHeight > b.height {
		// TODO: error handling
		return
	}
}

func (b *Bar) renderHorizontal(models map[string]tea.Model) {
	if len(b.bars) == 0 {
		// TODO: error handling
		return
	}

	totalWidth := 0
	for _, c := range b.bars {
		c.render(models)
		totalWidth += c.width
	}

	if totalWidth > b.width {
		// TODO: error handling
		return
	}
}

func (b *Bar) joinBars() {
	var bars []string
	if !b.rendered {
		for _, c := range b.bars {
			c.joinBars()
			bars = append(bars, c.view)
		}
		switch b.layoutType {
		case LIST:
			b.view = b.getStyle().
				// Width(b.width).
				// Height(b.height).
				Render(lipgloss.JoinVertical(0, bars...))
		case LINEAR:
			b.view = b.getStyle().
				// Width(b.width).
				// Height(b.height).
				Render(lipgloss.JoinHorizontal(0, bars...))
		}
		b.rendered = true
	}
}

func (b Bar) GetID() string {
	return b.id
}

type barOptions func(*Bar)

func WithID(v string) func(*Bar) {
	return func(b *Bar) {
		b.id = v
	}
}

func WithParent(v int) func(*Bar) {
	return func(b *Bar) {
		b.SetParent(v)
	}
}

func WithDynamic() func(*Bar) {
	return func(b *Bar) {
		b.SetDynamic()
	}
}

func WithFixed(v int) func(*Bar) {
	return func(b *Bar) {
		b.SetFixed(v)
	}
}

func WithLayout(v LayoutType) func(*Bar) {
	return func(b *Bar) {
		b.layoutType = v
	}
}

func NewBar(bars []*Bar, opts ...barOptions) *Bar {
	ret := &Bar{
		id:         uuid.NewString(),
		bars:       bars,
		layoutType: LIST,
		flavour:    NewFlavour(),
	}
	ret.SetParent(1)
	ret.flavourPrefsFct = ret.defaultFlavourPrefs

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
