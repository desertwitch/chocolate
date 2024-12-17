package chocolate

import (
	"fmt"

	"github.com/mfulz/chocolate/internal/tree"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type selector struct {
	barMap      map[string]CChocolateBar
	selectables []string
	selectedIdx int
	selected    CChocolateBar
	focused     CChocolateBar
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

func (s selector) getByID(v string) CChocolateBar {
	if b, ok := s.barMap[v]; ok {
		return b
	}
	return nil
}

func (s selector) hasFocus(v CChocolateBar) bool {
	return s.focused == v
}

func (s selector) isSelected(v CChocolateBar) bool {
	return s.selected == v
}

func (s selector) getSelected() CChocolateBar {
	return s.selected
}

func (s selector) getFocused() CChocolateBar {
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

func (s *selector) selectBar(v CChocolateBar) {
	s.selected = v
}

func (s *selector) forceSelect(v CChocolateBar) {
	s.selected = v
	s.focused = v
}

type rootBar struct{}

// Chocolate is the main entry point and acts as
// a control handler to work with the layouts
// provided by ChocolateBar.
// It implements the tea.Model interface and can
// be used as program entry point.
// Further it provides some rough default key bindings
// to be directly usable as ui.
type Chocolate struct {
	// Key mappings
	KeyMap KeyMap

	// root bar
	bar CChocolateBar

	// bar selector used for easy access
	// and tracking of input focus and
	// selecting
	barctl *selector

	// theme
	// flavour *Flavour

	// disableSelector is used to tell the Chocolate
	// to directly hand over input focus to the
	// selected bar
	// This is usefull for something like a menu
	// which has focus on start and will load other
	// models on selection and changes selected
	// When disableSelector is enabled the default leave
	// keyMap will be disabled and the handling
	// of selecting, focusing, etc. will be handed
	// over to the bar and it's models
	disableSelector bool
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.handleResize(msg)
		return c, nil
	case ForceSelectMsg:
		if bar := c.GetBarByID(string(msg)); bar != nil {
			c.ForceSelect(bar)
		}
		return c, nil
	case SelectMsg:
		if bar := c.GetBarByID(string(msg)); bar != nil {
			c.Select(bar)
		}
		return c, nil
	case BarHideMsg:
		if bar := c.GetBarByID(msg.Id); bar != nil {
			if c.IsSelected(bar) && msg.Value {
				c.barctl.next()
			}
			bar.Hide(msg.Value)
		}
		return c, nil
	case ModelChangeMsg:
		if bar := c.GetBarByID(msg.Id); bar != nil {
			cmds = append(cmds, bar.HandleUpdate(msg))
		}
	case tea.KeyMsg:
		if b != nil {
			if key.Matches(msg, c.KeyMap.Release) && !c.disableSelector {
				c.barctl.unfocus()
			} else {
				cmds = append(cmds, b.HandleUpdate(msg))
			}
			if c.barctl.focused == nil {
				cmds = append(cmds, c.handleNavigation(msg))
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

func (c Chocolate) GetBarByID(v string) CChocolateBar {
	return c.barctl.getByID(v)
}

func (c Chocolate) GetFocused() CChocolateBar {
	ret := c.barctl.focused
	if ret != nil {
		return ret
	}

	ret = c.barctl.selected
	if ret != nil && ret.InputOnSelect() {
		return ret
	}
	return nil
}

func (c Chocolate) GetSelected() CChocolateBar {
	return c.barctl.selected
}

func (c Chocolate) IsSelected(v CChocolateBar) bool {
	return c.barctl.isSelected(v)
}

func (c Chocolate) IsFocused(v CChocolateBar) bool {
	return c.barctl.hasFocus(v)
}

func (c *Chocolate) ForceSelect(v CChocolateBar) {
	if !c.disableSelector {
		return
	}
	c.barctl.forceSelect(v)
}

func (c *Chocolate) Select(v CChocolateBar) {
	c.barctl.selectBar(v)
}

func (c Chocolate) View() string {
	return c.bar.Render()
}

func buildDefaultSelector(v CChocolateBar) (*selector, error) {
	selectables := v.GetSelectables()

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

func initBarMap(v CChocolateBar) (map[string]CChocolateBar, error) {
	ret := make(map[string]CChocolateBar)

	if _, ok := ret[v.GetID()]; ok {
		return nil, fmt.Errorf("ID: %s already exists", v.GetID())
	}
	ret[v.GetID()] = v

	for _, b := range v.GetBars() {
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

func (c *Chocolate) initBar(v CChocolateBar) error {
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

func WithAutofocus(v CChocolateBar) func(*Chocolate) {
	return func(c *Chocolate) {
		c.disableSelector = true
		c.ForceSelect(v)
	}
}

func NewChocolate(bar CChocolateBar, opts ...chocolateOptions) (*Chocolate, error) {
	ret := &Chocolate{
		KeyMap:          DefaultKeyMap(),
		disableSelector: false,
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

type NChocolate struct {
	// Key mappings
	KeyMap KeyMap

	// bar tree
	tree *tree.Tree[string, ChocolateBar]

	// navigation
	selectables []string
	selectIdx   int
	selected    ChocolateBar
	focused     bool
}

func (c *NChocolate) AddBar(pid string, bar ChocolateBar) error {
	// TODO: error handling
	bar.setChocolate(c)
	c.tree.Add(bar.GetID(), pid, bar)

	c.selectables = []string{}
	selectedCheck := newCheckAtrributes().CheckCanSelect(true)
	for selectable := range c.tree.FindAllBy(selectedCheck) {
		c.selectables = append(c.selectables, selectable.GetData().GetID())
	}

	return nil
}

func (c NChocolate) IsRoot(bar ChocolateBar) bool {
	return c.tree.Root().GetData().GetID() == bar.GetID()
}

func (c NChocolate) IsSelected(bar ChocolateBar) bool {
	if c.selected == nil {
		return false
	}
	return c.selected.GetID() == bar.GetID()
}

func (c NChocolate) HasFocus(bar ChocolateBar) bool {
	return c.IsSelected(bar) && c.focused
}

func (c *NChocolate) Next() {
	c.selectIdx++
	if c.selectIdx >= len(c.selectables) {
		c.selectIdx = 0
	}
	if selected, ok := c.tree.Find(c.selectables[c.selectIdx]); !ok {
		c.selectIdx--
		return
	} else {
		c.selected = selected.GetData()
	}
}

func (c *NChocolate) Prev() {
	c.selectIdx--
	if c.selectIdx < 0 {
		c.selectIdx = len(c.selectables) - 1
	}
	if selected, ok := c.tree.Find(c.selectables[c.selectIdx]); !ok {
		c.selectIdx++
		return
	} else {
		c.selected = selected.GetData()
	}
}

type chocolateOption func(*NChocolate)

func SetLayout(v LayoutType) chocolateOption {
	return func(c *NChocolate) {
		c.tree.Root().GetData().(LayoutBar).SetLayout(v)
	}
}

func NewNChocolate(opts ...chocolateOption) (*NChocolate, error) {
	ret := &NChocolate{
		KeyMap: DefaultKeyMap(),
		tree:   tree.NewTree[string, ChocolateBar](),
	}

	rootBar := NewLayoutBar(LIST,
		SetID("root"),
	)
	rootBar.setChocolate(ret)

	ret.tree.Add("root", "", rootBar)

	for _, opt := range opts {
		opt(ret)
	}

	return ret, nil
}
