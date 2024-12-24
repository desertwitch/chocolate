package chocolate

import (
	"reflect"

	"github.com/mfulz/chocolate/internal/tree"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type ChocolateCustomUpdateHandlerFct func(msg tea.Msg) (tea.Cmd, bool)

type Chocolate struct {
	// Key mappings
	KeyMap KeyMap

	// bar tree
	tree *tree.Tree[string, ChocolateBar]

	// navigation
	selectables []string
	selectIdx   int
	selected    ChocolateBar
	focused     bool
	selector    bool

	registedUpdater map[reflect.Type][]ChocolateCustomUpdateHandlerFct

	preUpdateHandler ChocolateCustomUpdateHandlerFct
}

func (c *Chocolate) AddBar(pid string, bar ChocolateBar) error {
	// TODO: error handling
	bar.setBarChocolate(c)
	c.tree.Add(bar.GetID(), pid, bar)

	c.selectables = []string{}
	for selectable := range c.tree.FindAllBy(
		func(bar ChocolateBar) bool { return bar.IsSelectable() }, false,
	) {
		c.selectables = append(c.selectables, selectable.GetData().GetID())
	}

	return nil
}

func (c *Chocolate) AddOverlayRoot(bar ChocolateBar) error {
	// TODO: error handling
	bar.setBarChocolate(c)
	bar.setOverlay()
	c.tree.Add(bar.GetID(), c.tree.Root().GetID(), bar)

	c.selectables = []string{}
	for selectable := range c.tree.FindAllBy(
		func(bar ChocolateBar) bool { return bar.IsSelectable() }, false,
	) {
		c.selectables = append(c.selectables, selectable.GetData().GetID())
	}

	return nil
}

func (c Chocolate) IsRoot(bar BarSelector) bool {
	return c.tree.Root().GetData().GetID() == bar.GetID()
}

func (c Chocolate) GetParent(bar BarSelector) BarParent {
	node, ok := c.tree.Find(bar.GetID())
	if !ok {
		return nil
	}
	if node.GetParent() != nil {
		return node.GetParent().GetData()
	}
	return nil
}

func (c *Chocolate) GetChildren(bar BarSelector) []BarChild {
	children := []BarChild{}

	node, ok := c.tree.Find(bar.GetID())
	if !ok {
		return children
	}
	nchildren := node.GetChildren()
	for _, child := range nchildren {
		bar := child.GetData()
		if bar == nil {
			continue
		}
		children = append(children, bar)
	}

	return children
}

func (c Chocolate) IsSelected(bar BarSelector) bool {
	if c.selected == nil {
		return false
	}
	return c.selected.GetID() == bar.GetID()
}

func (c Chocolate) GetSelected() ChocolateBar {
	return c.selected
}

func (c *Chocolate) Select(bar ChocolateBar) {
	if !c.selector {
		c.ForceSelect(bar)
		return
	}

	if bar == nil {
		return
	}
	for i, id := range c.selectables {
		if id == bar.GetID() {
			if c.focused {
				if !bar.IsFocusable() {
					c.focused = false
				}
			}
			c.selectIdx = i
			c.selected = bar
			return
		}
	}
}

func (c *Chocolate) ForceSelect(bar ChocolateBar) {
	c.selected = bar
	if !c.selector {
		c.focused = true
	}
}

func (c Chocolate) IsFocused(bar BarSelector) bool {
	return c.IsSelected(bar) && c.focused
}

func (c *Chocolate) Next() {
	if len(c.selectables) == 0 {
		return
	}
	c.selectIdx++
	if c.selectIdx >= len(c.selectables) {
		c.selectIdx = 0
	}
	if selected, ok := c.tree.Find(c.selectables[c.selectIdx]); !ok {
		c.selectIdx--
		return
	} else {
		if c.focused {
			if !selected.GetData().IsFocusable() {
				c.focused = false
			}
		}
		c.selected = selected.GetData()
	}
}

func (c Chocolate) GetByID(id string) ChocolateBar {
	node, ok := c.tree.Find(id)
	if ok {
		return node.GetData()
	}
	return nil
}

func (c *Chocolate) Prev() {
	if len(c.selectables) == 0 {
		return
	}
	c.selectIdx--
	if c.selectIdx < 0 {
		c.selectIdx = len(c.selectables) - 1
	}
	if selected, ok := c.tree.Find(c.selectables[c.selectIdx]); !ok {
		c.selectIdx++
		return
	} else {
		if c.focused {
			if !selected.GetData().IsFocusable() {
				c.focused = false
			}
		}
		c.selected = selected.GetData()
	}
}

func (c *Chocolate) Focus(bar BarSelector) {
	if bar.IsFocusable() {
		c.focused = true
	}
}

func (c *Chocolate) UnFocus() {
	c.focused = false
}

func (c *Chocolate) handleResize(msg tea.WindowSizeMsg) {
	for bar := range c.tree.FindAllBy(
		func(bar ChocolateBar) bool { return true }, true,
	) {
		bar.GetData().Resize(msg.Width, msg.Height)
	}
}

func (c *Chocolate) handleNavigation(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	b := c.selected
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.KeyMap.Quit):
			return tea.Quit
		case key.Matches(msg, c.KeyMap.NextBar):
			c.Next()
		case key.Matches(msg, c.KeyMap.PrevBar):
			c.Prev()
		case key.Matches(msg, c.KeyMap.Focus):
			if c.IsSelected(b) {
				c.Focus(b)
			}
		}
	}

	return tea.Batch(cmds...)
}

func (c *Chocolate) handleFocused(msg tea.Msg, bar BarUpdater) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.KeyMap.Release):
			c.UnFocus()
			return nil
		}
	}

	cmds = append(cmds, bar.HandleUpdate(msg))
	return tea.Batch(cmds...)
}

func (c *Chocolate) Init() tea.Cmd {
	return nil
}

func (c *Chocolate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if c.preUpdateHandler != nil {
		cmd, stop := c.preUpdateHandler(msg)
		cmds = append(cmds, cmd)
		if stop {
			return c, tea.Batch(cmds...)
		}
	}

	if fcts := c.getRegisteredUpdateFcts(msg); fcts != nil {
		for _, fct := range fcts {
			cmd, stop := fct(msg)
			cmds = append(cmds, cmd)
			if stop {
				return c, tea.Batch(cmds...)
			}
		}
	}

	b := c.selected
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.handleResize(msg)
		return c, nil
	case SelectMsg:
		c.Select(c.GetByID(string(msg)))
		return c, nil
	case ForceSelectMsg:
		c.ForceSelect(c.GetByID(string(msg)))
		return c, nil
	case BarHideMsg:
		bar := c.GetByID(msg.Id)
		if bar != nil {
			bar.Hide(msg.Value)
		}
		return c, nil
	case ModelChangeMsg:
		bar := c.GetByID(msg.Id)
		if bar != nil {
			cmds = append(cmds, b.HandleUpdate(msg))
		}
	case tea.KeyMsg:
		if c.IsFocused(b) && c.selector {
			cmds = append(cmds, c.handleFocused(msg, b))
		} else if c.IsSelected(b) {
			cmds = append(cmds, b.HandleUpdate(msg))
			cmds = append(cmds, c.handleNavigation(msg))
		} else {
			cmds = append(cmds, c.handleNavigation(msg))
		}
	}

	return c, tea.Batch(cmds...)
}

func (c *Chocolate) View() string {
	for bar := range c.tree.FindAllBy(
		func(bar ChocolateBar) bool { return true }, false,
	) {
		bar.GetData().PreRender()
	}

	for bar := range c.tree.FindAllBy(
		func(bar ChocolateBar) bool { return true }, false,
	) {
		bar.GetData().Render()
	}

	return c.tree.Root().GetData().GetView()
}

func (c *Chocolate) RegisterUpdateFor(msg tea.Msg, fct ChocolateCustomUpdateHandlerFct) {
	if fct == nil || msg == nil {
		return
	}
	c.registedUpdater[reflect.TypeOf(msg)] = append(c.registedUpdater[reflect.TypeOf(msg)], fct)
}

func (c Chocolate) getRegisteredUpdateFcts(msg interface{}) []ChocolateCustomUpdateHandlerFct {
	if fcts, ok := c.registedUpdater[reflect.TypeOf(msg)]; ok {
		return fcts
	}
	return nil
}

type chocolateOption func(*Chocolate)

func SetLayout(v LayoutType) chocolateOption {
	return func(c *Chocolate) {
		c.tree.Root().GetData().(BarLayouter).SetLayout(v)
	}
}

func WithoutSelector() chocolateOption {
	return func(c *Chocolate) {
		c.selector = false
	}
}

func WithPreUpdateHandler(v ChocolateCustomUpdateHandlerFct) chocolateOption {
	return func(c *Chocolate) {
		c.preUpdateHandler = v
	}
}

func NewNChocolate(opts ...chocolateOption) (*Chocolate, error) {
	ret := &Chocolate{
		KeyMap:          DefaultKeyMap(),
		tree:            tree.NewTree[string, ChocolateBar](),
		selector:        true,
		registedUpdater: make(map[reflect.Type][]ChocolateCustomUpdateHandlerFct),
	}

	rootBar := NewLayoutBar(LIST,
		WithBarID("root"),
	)
	rootBar.setBarChocolate(ret)

	ret.tree.Add("root", "", rootBar)

	for _, opt := range opts {
		opt(ret)
	}

	return ret, nil
}
