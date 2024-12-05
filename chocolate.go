package chocolate

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Chocolate struct {
	// Key mappings
	KeyMap KeyMap

	// root bar
	bar *Bar

	// theme
	flavour Flavour

	// models
	models map[string]tea.Model
}

func (c *Chocolate) handleResize(size tea.WindowSizeMsg) {
	log.Printf("chocolate: w=%d h=%d\n", size.Width, size.Height)
	if c.bar != nil {
		c.bar.Resize(size, c.models, c.bar)
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
		if msg.String() == "q" {
			return c, tea.Quit
		}
		// cmds = append(cmds, c.handleNavigation(msg))
	}

	return c, tea.Batch(cmds...)
}

func (c *Chocolate) handleNavigation(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.KeyMap.NextBar):
			// c.Next()
		}
	}

	return tea.Batch(cmds...)
}

func (c Chocolate) View() string {
	var ret string
	log.Printf("View called\n")
	c.bar.resetRender()
	c.bar.render(c.models)
	c.bar.joinBars()
	ret = c.bar.view
	c.bar.resetRender()
	return ret
}

func (c *Chocolate) SetModel(id string, model tea.Model) error {
	if _, ok := c.models[id]; !ok {
		return fmt.Errorf("No bar with id: %s", id)
	}
	c.models[id] = model
	return nil
}

type chocolateOptions func(*Chocolate)

func (c Chocolate) initModels(bar *Bar) {
	c.models[bar.id] = nil
	for _, b := range bar.bars {
		c.initModels(b)
	}
}

func NewChocolate(bar *Bar, opts ...chocolateOptions) *Chocolate {
	ret := &Chocolate{
		KeyMap:  DefaultKeyMap(),
		flavour: NewFlavour(),
		bar:     bar,
		models:  make(map[string]tea.Model),
	}

	for _, opt := range opts {
		opt(ret)
	}

	ret.initModels(bar)

	return ret
}
