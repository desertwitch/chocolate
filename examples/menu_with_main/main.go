package main

import (
	"fmt"
	"io"

	"gitea.olznet.de/mfulz/chocolate"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

	fp := chocolate.NewFlavourPrefs()
	s := i.flavour.GetStyle(fp)
	fn := s.Render

	if index == m.Index() {
		fn = i.flavour.GetStyle(fp.
			Foreground(chocolate.FOREGROUND_HIGHLIGHT_PRIMARY).
			Background(chocolate.BACKGROUND_HIGHLIGHT_PRIMARY),
		).Render
	}

	fmt.Fprint(w, fn(i.name))
}

type MainChangeMsg string

type menuModel struct {
	items   list.Model
	name    string
	dst     string
	flavour chocolate.Flavour
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m *menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.items.SetWidth(msg.Width)
		m.items.SetHeight(msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.items.SelectedItem().(*menuModel)
			if ok {
				return m, func() tea.Msg { return MainChangeMsg(i.dst) }
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

func NewMenuModel(name string, items []list.Item, dst string, flavour chocolate.Flavour) *menuModel {
	const defaultWidth = 50
	const defaultHeight = 50
	l := list.New(items, menuItemDelegate{}, defaultWidth, defaultHeight)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)

	ret := &menuModel{
		items:   l,
		name:    name,
		dst:     dst,
		flavour: flavour,
	}

	return ret
}

type mainModel string

func (t mainModel) Init() tea.Cmd                           { return nil }
func (t mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t mainModel) View() string                            { return string(t) }

var flavour = chocolate.NewFlavour()

var menuBarFlavourPrefsHandler = func(b *chocolate.ChocolateBar) func() chocolate.FlavourPrefs {
	return func() chocolate.FlavourPrefs {
		return chocolate.NewFlavourPrefs()
	}
}

var menuBarUpdateHandler = func(b *chocolate.ChocolateBar, m tea.Model) func(tea.Msg) tea.Cmd {
	return func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case MainChangeMsg:
			model := string(msg)
			bar := b.GetChoc().GetBarByID("main")
			bar.SelectModel(model)
			b.GetChoc().ForceSelect(bar)
		}
		return nil
	}
}

var mainBarUpdateHandler = func(b *chocolate.ChocolateBar, m tea.Model) func(tea.Msg) tea.Cmd {
	return func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "esc":
				b.SelectModel("dummy")
				b.GetChoc().ForceSelect(b.GetChoc().GetBarByID("menu"))
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
	mainModels["dummy"] = &chocolate.BarModel{Model: mainDummy, UpdateHandlerFct: mainBarUpdateHandler}
	mainModels["first"] = &chocolate.BarModel{Model: mainFirst, UpdateHandlerFct: mainBarUpdateHandler}
	mainModels["second"] = &chocolate.BarModel{Model: mainSecond, UpdateHandlerFct: mainBarUpdateHandler}

	mainContentBar := chocolate.NewChocolateBar(nil,
		chocolate.WithID("main"),
		chocolate.WithModels(mainModels, "dummy"),
	)

	menuModel := NewMenuModel("Main Menu",
		[]list.Item{
			NewMenuModel("First", nil, "first", flavour),
			NewMenuModel("Second", nil, "second", flavour),
		},
		"",
		flavour,
	)

	menuBar := chocolate.NewChocolateBar(nil,
		chocolate.WithModel(&chocolate.BarModel{Model: menuModel, UpdateHandlerFct: menuBarUpdateHandler}),
		chocolate.WithID("menu"),
		chocolate.WithFlavourPrefsHandle(menuBarFlavourPrefsHandler),
		chocolate.WithXScaler(chocolate.NewFixedScaler(20)),
	)

	bar := chocolate.NewChocolateBar([]*chocolate.ChocolateBar{
		menuBar,
		mainContentBar,
	},
		chocolate.WithLayout(chocolate.LINEAR),
	)

	if m, err := chocolate.NewChocolate(bar,
		chocolate.WithAutofocus(menuBar),
	); err != nil {
		panic(err)
	} else {
		if _, err := tea.NewProgram(m,
			tea.WithAltScreen()).Run(); err != nil {
			fmt.Println(err)
		}
	}
}
