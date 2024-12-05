package chocolate

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Chocolate struct {
	// Key mappings
	KeyMap KeyMap

	// root bar
	bar *ChocolateBar

	// theme
	flavour Flavour
}

func (c *Chocolate) handleResize(size tea.WindowSizeMsg) {
	log.Printf("chocolate: w=%d h=%d\n", size.Width, size.Height)
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
	return c.bar.Render()
}

type chocolateOptions func(*Chocolate)

func NewChocolate(bar *ChocolateBar, opts ...chocolateOptions) *Chocolate {
	ret := &Chocolate{
		KeyMap:  DefaultKeyMap(),
		flavour: NewFlavour(),
		bar:     bar,
	}

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
