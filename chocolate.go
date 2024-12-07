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
	if s.selected.CanFocus() {
		s.focused = s.selected
	}
}

func (s *selector) unfocus() {
	s.focused = nil
}

func (s *selector) forceSelect(v *ChocolateBar) {
	s.selected = v
	s.focused = v
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
	// flavour *Flavour

	// autofocus is used to tell the Chocolate
	// to directly hand over input focus to the
	// selected bar
	// This is usefull for something like a menu
	// which has focus on start and will load other
	// models on selection and changes selected
	// When autofocus is enabled the default leave
	// keyMap will be disabled and the handling
	// of selecting, focusing, etc. will be handed
	// over to the bar and it's models
	autofocus bool
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

	b := c.GetFocused()
	if b != nil {
		cmds = append(cmds, b.HandleUpdate(msg))
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.handleResize(msg)
		return c, nil
	case tea.KeyMsg:
		if b != nil {
			if key.Matches(msg, c.KeyMap.Release) && !c.autofocus {
				c.barctl.unfocus()
			} else {
				cmds = append(cmds, b.HandleUpdate(msg))
			}
		} else {
			cmds = append(cmds, c.handleNavigation(msg))
		}
	default:
		if b != nil {
			cmds = append(cmds, b.HandleUpdate(msg))
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

func (c *Chocolate) ForceSelect(v *ChocolateBar) {
	if !c.autofocus {
		return
	}
	c.barctl.forceSelect(v)
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

func buildDefaultSelector(v *ChocolateBar) (*selector, error) {
	selectables := getSelectables(v)

	if len(selectables) == 0 {
		// if nothing is selectable
		// at least the root bar must be
		v.selectable = true
		selectables = append(selectables, v.id)
	}

	barMap, err := initBarMap(v)
	if err != nil {
		return nil, err
	}

	ret := &selector{
		barMap:      barMap,
		selectables: selectables,
		selectedIdx: -1,
	}
	ret.next()

	return ret, nil
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

func WithAutofocus(v *ChocolateBar) func(*Chocolate) {
	return func(c *Chocolate) {
		c.autofocus = true
		c.ForceSelect(v)
	}
}

func NewChocolate(bar *ChocolateBar, opts ...chocolateOptions) (*Chocolate, error) {
	ret := &Chocolate{
		KeyMap:    DefaultKeyMap(),
		autofocus: false,
	}

	// bar initializing
	var err error
	if err = ret.initBar(bar); err != nil {
		return nil, err
	}
	if ret.barctl, err = buildDefaultSelector(bar); err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(ret)
	}

	return ret, nil
}
