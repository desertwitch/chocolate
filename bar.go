package chocolate

import (
	"github.com/google/uuid"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LayoutType int

const (
	LIST LayoutType = iota
	LINEAR
)

// Scaling defines how the ChocolateBar will be scaled
// FIXED will be a fixed number of cells
// RELATIVE will grow or shrink on screen resizes relative to the set values
type Scaling int

const (
	PARENT Scaling = iota
	DYNAMIC
	FIXED
)

// Alignment defines how the inner content will be aligned
type Alignment int

const (
	NONE   Alignment = iota // No alignment
	START                   // Will be either left or top (horizontal / vertical)
	END                     // Will be either right or bottom (horizontal / vertical)
	CENTER                  // Will be centered
)

type Sizing struct {
	Scaling   Scaling
	Value     int
	Alignment Alignment
}

type Bar struct {
	// Key mappings
	KeyMap KeyMap

	id          string
	name        string
	selectable  bool
	selected    bool
	flavourType FlavourType
	flavour     Flavour

	// map to hold all bars
	// used to easy select specific
	// bars which can make up dynamic
	// switching between different views
	barsReferences map[string]*Bar

	// bars in order for the layout
	bars []*Bar

	// layout parameters
	layoutType   LayoutType
	widthSizing  Sizing
	heightSizing Sizing

	// sizing information
	width      int
	height     int
	widthPart  int
	heightPart int

	// sizing
	maxWidth    int
	maxHeight   int
	parentParts int
}

func (b *Bar) hhandleResize(width, height int) {
	frameWidth, frameHeight := b.flavour.GetFrameSize()
	b.maxWidth = width - frameWidth
	b.maxHeight = height - frameHeight
	for _, c := range b.bars {
		c.hhandleResize(b.maxWidth, b.maxHeight)
		switch b.layoutType {
		case LIST:
			if c.heightSizing.Scaling == PARENT {
				b.parentParts += c.heightSizing.Value
			}
		case LINEAR:
			if c.widthSizing.Scaling == PARENT {
				b.parentParts += c.widthSizing.Value
			}
		}
	}
	if b.parentParts == 0 {
		b.parentParts = 1
	}
}

func (b *Bar) handleResize(width, height int) {
	var sizingFct func(c *Bar) Sizing
	var resizeFct func(c *Bar)
	var setPartFct func(c *Bar, fixed, parts int)

	frameWidth, frameHeight := b.flavour.GetFrameSize()
	b.width = width - frameWidth
	b.height = height - frameHeight

	// select direction
	switch b.layoutType {
	case LINEAR:
		sizingFct = func(c *Bar) Sizing { return c.widthSizing }
		setPartFct = func(c *Bar, fixed, parts int) {
			c.widthPart = (c.width - fixed) / parts
		}
		resizeFct = func(c *Bar) {
			switch c.widthSizing.Scaling {
			case FIXED:
				c.handleResize(c.widthSizing.Value, b.height)
			case DYNAMIC, PARENT:
				c.handleResize((c.widthSizing.Value * b.widthPart), b.height)
			}
		}
		b.heightPart = 0
	case LIST:
		sizingFct = func(c *Bar) Sizing { return c.heightSizing }
		setPartFct = func(c *Bar, fixed, parts int) {
			c.heightPart = (c.height - fixed) / parts
		}
		resizeFct = func(c *Bar) {
			switch c.heightSizing.Scaling {
			case FIXED:
				c.handleResize(b.width, c.heightSizing.Value)
			case DYNAMIC, PARENT:
				c.handleResize(b.width, (c.heightSizing.Value * b.heightPart))
			}
		}
		b.widthPart = 0
	}

	if len(b.bars) == 0 || b.bars == nil {
		return
	}

	var fixedTotal int
	var totalParts int
	for _, c := range b.bars {
		sizing := sizingFct(c)
		switch sizing.Scaling {
		case FIXED:
			fixedTotal += sizing.Value
		case DYNAMIC, PARENT:
			totalParts += sizing.Value
		}
	}

	if totalParts == 0 {
		totalParts = 1 // TODO: verify if this is working
	}
	setPartFct(b, fixedTotal, totalParts)

	for _, c := range b.bars {
		resizeFct(c)
	}
}

func (b Bar) ID() string {
	return b.id
}

func (b Bar) FlavourType() FlavourType {
	return b.flavourType
}

func (b Bar) Selectable() bool {
	return b.selectable
}

func (b Bar) Selected() bool {
	return b.selected
}

func (b *Bar) SetFlavourType(v FlavourType) {
	b.flavourType = v
}

func (b Bar) Init() tea.Cmd {
	return nil
}

func (b Bar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.hhandleResize(msg.Width, msg.Height)
		return b, nil
	case tea.KeyMsg:
		if msg.String() == "q" {
			return b, tea.Quit
		}
	}

	return b, nil
}

func (b Bar) contentView() string {
	return b.name
}

func (b Bar) View() string {
	var ret string
	bars := make([]string, len(b.bars))
	heightLeft := b.maxHeight
	widthLeft := b.maxWidth
	if b.parentParts == 0 {
		b.parentParts = 1
	}

	var s lipgloss.Style
	for i, c := range b.bars {
		s = lipgloss.NewStyle()
		if c.widthSizing.Scaling == FIXED {
			s = s.Width(c.widthSizing.Value)
		}
		if c.heightSizing.Scaling == FIXED {
			s = s.Height(c.heightSizing.Value)
		}
		s = s.MaxWidth(c.maxWidth)
		s = s.MaxHeight(c.maxHeight)
		switch b.layoutType {
		case LIST:
			if c.heightSizing.Scaling != PARENT {
				v := c.View()
				heightLeft -= lipgloss.Height(v)
				bars[i] = s.Render(v)
			}
		case LINEAR:
			if c.widthSizing.Scaling != PARENT {
				v := c.View()
				widthLeft -= lipgloss.Width(v)
				bars[i] = s.Render(v)
			}
		}
	}

	for i, c := range b.bars {
		s = lipgloss.NewStyle()
		if c.widthSizing.Scaling == FIXED {
			s = s.Width(c.widthSizing.Value)
		}
		if c.heightSizing.Scaling == FIXED {
			s = s.Height(c.heightSizing.Value)
		}
		s = s.MaxWidth(c.maxWidth)
		s = s.MaxHeight(c.maxHeight)
		switch b.layoutType {
		case LIST:
			if c.heightSizing.Scaling == PARENT {
				s = s.Height(c.heightSizing.Value * (heightLeft / b.parentParts))
				v := c.View()
				bars[i] = s.Render(v)
			}
		case LINEAR:
			if c.widthSizing.Scaling == PARENT {
				s = s.Width(c.widthSizing.Value * (widthLeft / b.parentParts))
				v := c.View()
				bars[i] = s.Render(v)
			}
		}
	}

	s = lipgloss.NewStyle().
		// Width(b.width).
		// Height(b.height).
		// MaxWidth(b.width).
		// MaxHeight(b.height).
		Border(lipgloss.RoundedBorder())
	switch b.widthSizing.Scaling {
	case FIXED:
		s = s.Width(b.maxWidth)
	default:
		s = s.MaxWidth(b.maxWidth)
	}

	switch b.heightSizing.Scaling {
	case FIXED:
		s = s.Height(b.maxHeight)
	default:
		s = s.MaxHeight(b.maxHeight)
	}

	if len(b.bars) == 0 {
		ret = b.contentView()
	} else {
		switch b.layoutType {
		case LIST:
			ret = lipgloss.JoinVertical(0, bars...)
		case LINEAR:
			ret = lipgloss.JoinHorizontal(0, bars...)
		}
	}
	return s.Render(ret)
}

func (b Bar) view() string {
	var ret string
	var bars []string

	s := lipgloss.NewStyle().
		// Width(b.width).
		// Height(b.height).
		// MaxWidth(b.width).
		// MaxHeight(b.height).
		Border(lipgloss.RoundedBorder())

	switch b.widthSizing.Scaling {
	case FIXED, PARENT:
		s = s.Width(b.width)
	default:
		s = s.MaxWidth(b.width)
	}

	switch b.heightSizing.Scaling {
	case FIXED, PARENT:
		s = s.Height(b.height)
	default:
		s = s.MaxHeight(b.height)
	}

	for _, b := range b.bars {
		bars = append(bars, b.View())
	}

	if len(b.bars) == 0 {
		ret = b.contentView()
	} else {
		switch b.layoutType {
		case LIST:
			ret = lipgloss.JoinVertical(0, bars...)
		case LINEAR:
			ret = lipgloss.JoinHorizontal(0, bars...)
		}
	}
	return s.Render(ret)
}

type chocolateBarOptions func(*Bar)

func ID(v string) chocolateBarOptions {
	return func(b *Bar) {
		b.id = v
	}
}

func Name(v string) chocolateBarOptions {
	return func(b *Bar) {
		b.name = v
	}
}

func WithSizing(width Sizing, height Sizing) chocolateBarOptions {
	return func(b *Bar) {
		b.widthSizing = width
		b.heightSizing = height
	}
}

func Selectable() chocolateBarOptions {
	return func(b *Bar) {
		b.selectable = true
	}
}

func WithFlavourType(v FlavourType) chocolateBarOptions {
	return func(b *Bar) {
		b.flavourType = v
	}
}

func WithLayoutType(v LayoutType) chocolateBarOptions {
	return func(b *Bar) {
		b.layoutType = v
	}
}

func NewChocolateBar(bars []*Bar, opts ...chocolateBarOptions) *Bar {
	ret := &Bar{
		KeyMap:         DefaultKeyMap(),
		id:             uuid.NewString(),
		flavour:        NewFlavour(),
		widthSizing:    Sizing{PARENT, 1, CENTER},
		heightSizing:   Sizing{PARENT, 1, CENTER},
		selectable:     true,
		selected:       false,
		flavourType:    FLAVOUR_PRIMARY,
		barsReferences: make(map[string]*Bar),
		bars:           bars,
	}

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
