package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mfulz/chocolate"
	flavour "github.com/mfulz/chocolate/flavour"
)

func (i menuModel) FilterValue() string { return "" }

type menuItemDelegate struct{}

func (d menuItemDelegate) Height() int                             { return 1 }
func (d menuItemDelegate) Spacing() int                            { return 0 }
func (d menuItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d menuItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*menuModel)
	if !ok {
		return
	}

	s := flavour.GetPresetNoErr(flavour.PRESET_PRIMARY_NOBORDER).
		Width(i.width)
	fn := s.Render

	if index == m.Index() {
		fn = flavour.GetPresetNoErr(flavour.PRESET_SECONDARY_NOBORDER).
			Width(i.width).
			Render
	}

	fmt.Fprint(w, fn(i.name))
}

type MainChangeMsg string

type menuModel struct {
	items  list.Model
	name   string
	dst    string
	choice string
	width  int
}

func (m menuModel) Init() tea.Cmd { return nil }

func (m *menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.items.SetWidth(msg.Width)
		m.items.SetHeight(msg.Height)
		for _, i := range m.items.Items() {
			ie := i.(*menuModel)
			ie.Update(msg)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.choice = ""
			i, ok := m.items.SelectedItem().(*menuModel)
			if ok {
				var cmds []tea.Cmd
				cmds = append(cmds,
					func() tea.Msg {
						return chocolate.ModelChangeMsg{
							Id:    "main",
							Model: i.dst,
						}
					},
					func() tea.Msg {
						return chocolate.ForceSelectMsg("main")
					},
				)
				return m, tea.Batch(cmds...)
			}
		}
	}

	var cmd tea.Cmd
	m.items, cmd = m.items.Update(msg)
	return m, cmd
}

func (m menuModel) View() string {
	return m.items.View()
}

func NewMenuModel(name string, items []list.Item, dst string) *menuModel {
	const defaultWidth = 50
	const defaultHeight = 50
	l := list.New(items, menuItemDelegate{}, defaultWidth, defaultHeight)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)

	ret := &menuModel{
		items: l,
		name:  name,
		dst:   dst,
	}

	return ret
}

type mainModel string

func (t mainModel) Init() tea.Cmd                           { return nil }
func (t mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t mainModel) View() string                            { return string(t) }

var menuBarFlavourCustomizer = func(
	b chocolate.ChocolateBar,
	m tea.Model, s lipgloss.Style,
) func() lipgloss.Style {
	return func() lipgloss.Style {
		return flavour.GetPresetNoErr(flavour.PRESET_PRIMARY_NOBORDER).
			MarginTop(1).
			MarginLeft(3).
			MarginRight(3)
	}
}

var mainBarUpdateHandler = func(b chocolate.ChocolateBar, m tea.Model) func(tea.Msg) tea.Cmd {
	return func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "esc":
				b.SelectModel("dummy")
				b.ForceSelect(b.GetByID("menu"))
			}
		}
		return nil
	}
}

func main() {
	mainDummy := mainModel("")
	mainFirst := mainModel("first")
	mainSecond := mainModel("second")

	mainModels := make(map[string]*chocolate.BarModel)
	mainModels["dummy"] = &chocolate.BarModel{Model: mainDummy}
	mainModels["first"] = &chocolate.BarModel{Model: mainFirst, UpdateHandlerFct: mainBarUpdateHandler}
	mainModels["second"] = &chocolate.BarModel{Model: mainSecond, UpdateHandlerFct: mainBarUpdateHandler}

	mainContentBar := chocolate.NewMultiModelBar(
		"dummy",
		mainModels,
		chocolate.WithBarID("main"),
	)

	menuModel := NewMenuModel("Main Menu",
		[]list.Item{
			NewMenuModel("First", nil, "first"),
			NewMenuModel("Second", nil, "second"),
		},
		"",
	)

	menuBar := chocolate.NewModelBar(
		&chocolate.BarModel{
			Model:                   menuModel,
			FlavourCustomizeHandler: menuBarFlavourCustomizer,
		},
		chocolate.WithBarID("menu"),
		chocolate.WithBarXScaler(chocolate.FIXED, 20),
	)

	if m, err := chocolate.NewNChocolate(
		chocolate.SetLayout(chocolate.LINEAR),
	); err != nil {
		panic(err)
	} else {
		m.AddBar("root", menuBar)
		m.AddBar("root", mainContentBar)
		m.ForceSelect(m.GetByID("menu"))
		if _, err := tea.NewProgram(m,
			tea.WithAltScreen()).Run(); err != nil {
			fmt.Println(err)
		}
	}
}
