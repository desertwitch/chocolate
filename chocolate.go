package chocolate

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type selector struct {
	selectables []string
	selected    int
}

func (s *selector) Next() {
	s.selected++
	if s.selected >= len(s.selectables) {
		s.selected = 0
	}
}

func (s *selector) Prev() {
	s.selected--
	if s.selected < 0 {
		s.selected = len(s.selectables) - 1
	}
}

func (s selector) Get() string {
	return s.selectables[s.selected]
}

type Chocolate struct {
	// Key mappings
	KeyMap KeyMap

	// root bar
	bar *ChocolateBar

	// bar selector
	activeBar  *selector
	inputFocus bool

	// bar map for easy access
	bars map[string]*ChocolateBar

	// theme
	flavour Flavour
}

func (c *Chocolate) handleResize(size tea.WindowSizeMsg) {
	if c.bar != nil {
		c.bar.Resize(size.Width, size.Height)
	}
}

func (c Chocolate) Init() tea.Cmd {
	return nil
}

func (c Chocolate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.handleResize(msg)
		return c, nil
	case tea.KeyMsg:
		if c.inputFocus {
			if key.Matches(msg, c.KeyMap.Release) {
				c.inputFocus = false
				c.focusBar(false)
			}
			if b := c.getFocusedBar(); b != nil {
				cmds = append(cmds, b.HandleUpdate(msg))
			}
		} else {
			cmds = append(cmds, c.handleNavigation(msg))
		}
	}

	return c, tea.Batch(cmds...)
}

func (c *Chocolate) handleNavigation(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.KeyMap.Quit):
			return tea.Quit
		case key.Matches(msg, c.KeyMap.NextBar):
			c.next()
		case key.Matches(msg, c.KeyMap.PrevBar):
			c.prev()
		case key.Matches(msg, c.KeyMap.Focus):
			c.inputFocus = true
			c.focusBar(true)
		}
	}

	return tea.Batch(cmds...)
}

func (c Chocolate) GetBarByID(v string) *ChocolateBar {
	if b, ok := c.bars[v]; ok {
		return b
	}
	return nil
}

func (c Chocolate) getFocusedBar() *ChocolateBar {
	if b, ok := c.bars[c.activeBar.Get()]; ok {
		return b
	}
	return nil
}

func (c *Chocolate) selectBar(v bool) {
	if b, ok := c.bars[c.activeBar.Get()]; ok {
		b.Select(v)
	}
}

func (c *Chocolate) focusBar(v bool) {
	if b, ok := c.bars[c.activeBar.Get()]; ok {
		b.Focus(v)
	}
}

func (c *Chocolate) next() {
	c.selectBar(false)
	c.activeBar.Next()
	c.selectBar(true)
}

func (c *Chocolate) prev() {
	c.selectBar(false)
	c.activeBar.Prev()
	c.selectBar(true)
}

func (c Chocolate) View() string {
	return c.bar.Render()
}

func getSelectables(v *ChocolateBar) []string {
	ret := []string{}

	if v.selectable {
		ret = append(ret, v.id)
	}

	for _, b := range v.bars {
		ret = append(ret, getSelectables(b)...)
	}

	return ret
}

func buildDefaultSelector(v *ChocolateBar) *selector {
	selectables := getSelectables(v)

	if len(selectables) == 0 {
		// if nothing is selectable
		// at least the root bar must be
		v.selectable = true
		selectables = append(selectables, v.id)
	}

	ret := &selector{
		selectables: selectables,
		selected:    0,
	}

	return ret
}

func initBarMap(v *ChocolateBar) (map[string]*ChocolateBar, error) {
	ret := make(map[string]*ChocolateBar)

	if _, ok := ret[v.id]; ok {
		return nil, fmt.Errorf("ID: %s already exists", v.id)
	}
	ret[v.id] = v

	for _, b := range v.bars {
		if submap, err := initBarMap(b); err != nil {
			return nil, err
		} else {
			for k, b := range submap {
				if _, ok := ret[k]; ok {
					return nil, fmt.Errorf("ID: %s already exists", k)
				} else {
					ret[k] = b
				}
			}
		}
	}
	return ret, nil
}

type chocolateOptions func(*Chocolate)

func WithFlavor(v Flavour) func(*Chocolate) {
	return func(c *Chocolate) {
		c.flavour = v
	}
}

func WithSelector(v []string, s int) func(*Chocolate) {
	return func(c *Chocolate) {
		if len(v) == 0 {
			c.activeBar = nil
		}
		c.activeBar = &selector{
			selectables: v,
			selected:    s,
		}
	}
}

func NewChocolate(bar *ChocolateBar, opts ...chocolateOptions) *Chocolate {
	ret := &Chocolate{
		KeyMap:     DefaultKeyMap(),
		flavour:    NewFlavour(),
		bar:        bar,
		inputFocus: false,
	}

	if bar == nil {
		// TODO: error handling
		return nil
	}

	var err error
	if ret.bars, err = initBarMap(bar); err != nil {
		// TODO error handling
		return nil
	}

	for _, opt := range opts {
		opt(ret)
	}

	if ret.activeBar == nil {
		ret.activeBar = buildDefaultSelector(bar)
		ret.selectBar(true)
	}

	return ret
}
