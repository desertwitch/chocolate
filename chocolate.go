package chocolate

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type selector struct {
	barMap      map[string]*ChocolateBar
	selectables []string
	selectedIdx int
	selected    *ChocolateBar
	focused     *ChocolateBar
}

func (s *selector) next() {
	s.selectedIdx++
	if s.selectedIdx >= len(s.selectables) {
		s.selectedIdx = 0
	}
	s.selected = s.getByID(s.selectables[s.selectedIdx])
}

func (s *selector) prev() {
	s.selectedIdx--
	if s.selectedIdx < 0 {
		s.selectedIdx = len(s.selectables) - 1
	}
	s.selected = s.getByID(s.selectables[s.selectedIdx])
}

func (s selector) getByID(v string) *ChocolateBar {
	if b, ok := s.barMap[v]; ok {
		return b
	}
	return nil
}

func (s selector) hasFocus(v *ChocolateBar) bool {
	return s.focused == v
}

func (s selector) isSelected(v *ChocolateBar) bool {
	return s.selected == v
}

func (s selector) getSelected() *ChocolateBar {
	return s.selected
}

func (s selector) getFocused() *ChocolateBar {
	return s.focused
}

func (s *selector) focus() {
	s.focused = s.selected
}

func (s *selector) unfocus() {
	s.focused = nil
}

type Chocolate struct {
	// Key mappings
	KeyMap KeyMap

	// root bar
	bar *ChocolateBar

	// bar selector used for easy access
	// and tracking of input focus and
	// selecting
	barctl *selector

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
		if b := c.GetFocused(); b != nil {
			if key.Matches(msg, c.KeyMap.Release) {
				c.barctl.unfocus()
			} else {
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
			c.barctl.next()
		case key.Matches(msg, c.KeyMap.PrevBar):
			c.barctl.prev()
		case key.Matches(msg, c.KeyMap.Focus):
			c.barctl.focus()
		}
	}

	return tea.Batch(cmds...)
}

func (c Chocolate) GetBarByID(v string) *ChocolateBar {
	return c.barctl.getByID(v)
}

func (c Chocolate) GetFocused() *ChocolateBar {
	return c.barctl.focused
}

func (c Chocolate) GetSelected() *ChocolateBar {
	return c.barctl.selected
}

func (c Chocolate) IsSelected(v *ChocolateBar) bool {
	return c.barctl.isSelected(v)
}

func (c Chocolate) IsFocused(v *ChocolateBar) bool {
	return c.barctl.hasFocus(v)
}

func (c Chocolate) GetFlavour() Flavour {
	return c.flavour
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

	barMap, err := initBarMap(v)
	if err != nil {
		return nil
	}

	ret := &selector{
		barMap:      barMap,
		selectables: selectables,
		selectedIdx: -1,
	}
	ret.next()

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

func (c *Chocolate) initBar(v *ChocolateBar) error {
	if v == nil || !v.IsRoot() {
		return fmt.Errorf("Not a root bar")
	}

	// set the chocolate
	v.SetChocolate(c)

	// add the var to the chocolate
	c.bar = v

	return nil
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
			return
		}
		if s < 0 || s >= len(v) {
			s = 0
		}

		c.barctl.selectables = v
		c.barctl.selectedIdx = s - 1
		c.barctl.next()
	}
}

func NewChocolate(bar *ChocolateBar, opts ...chocolateOptions) *Chocolate {
	ret := &Chocolate{
		KeyMap:  DefaultKeyMap(),
		flavour: NewFlavour(),
	}

	// bar initializing
	if err := ret.initBar(bar); err != nil {
		// TODO: error handling
		return nil
	}
	if ret.barctl = buildDefaultSelector(bar); ret.barctl == nil {
		// TODO: error handling
		return nil
	}

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
