package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mfulz/chocolate"
)

// styles used for flavour
var (
	viewStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("246")).Background(lipgloss.Color("232")).
			BorderForeground(lipgloss.Color("246")).BorderBackground(lipgloss.Color("232"))
	selectedStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("15")).Background(lipgloss.Color("237")).
			BorderForeground(lipgloss.Color("15")).BorderBackground(lipgloss.Color("237"))
	focusedStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("196")).Background(lipgloss.Color("232")).
			BorderForeground(lipgloss.Color("15")).BorderBackground(lipgloss.Color("232"))
)

type menuModel struct {
	entries  []string
	selected int
	choice   string
	width    int
}

func (m *menuModel) Init() tea.Cmd { return nil }

// Resize is used by chocolate.BarModel and will be called
// after layout calculation when the size changed
func (m *menuModel) Resize(width, height int) {
	m.width = width - 4
}

func (m *menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			m.selected++
			if m.selected >= len(m.entries) {
				m.selected = 0
			}
		case "k":
			m.selected--
			if m.selected < 0 {
				m.selected = len(m.entries) - 1
			}
		}
		m.choice = m.entries[m.selected]
	}

	return m, nil
}

func (m *menuModel) View() string {
	entries := []string{}
	vs := viewStyle.Width(m.width)
	ss := selectedStyle.Width(m.width)
	for i, e := range m.entries {
		if i == m.selected {
			entries = append(entries, ss.Render(e))
		} else {
			entries = append(entries, vs.Render(e))
		}
	}

	return lipgloss.JoinVertical(0, entries...)
}

func NewMenuModel(items []string) *menuModel {
	return &menuModel{
		entries:  items,
		selected: 0,
		width:    10,
	}
}

// tea.Model and program start
// it holds a reference to the main chocolate
// which acts as the central control interface
// to the layout and styling functionality
// provided by chocolate
type model struct {
	choc *chocolate.Chocolate
	menu *menuModel
}

func (m *model) Init() tea.Cmd { return nil }
func (m *model) View() string  { return m.choc.View() }
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// chocolate.Chocolate implements the chocolate.BarModel interface
		// The resizing is required initially and for changing the window
		// size to recalculate the layout
		m.choc.Resize(msg.Width, msg.Height)

		// Select the model by name (m.menu.choice) for the bar (contentbar)
		// This is used to have multiple models placed in a container (contentbar)
		// and change them dynamically.
		// This is acting like a state to exchange ui elements depending on
		// other circumstances, like in this example to display another model
		// in the right main part of the application depending on the selected
		// menu entry
		m.choc.SelectModel(m.menu.choice, "contentbar")
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
		_, cmd = m.menu.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.choc.SelectModel(m.menu.choice, "contentbar")

	return m, tea.Batch(cmds...)
}

func main() {
	// the flavour with the three styles from above
	// which is used by chocolate to provide a central styling
	// for everything inside it
	theme := chocolate.NewChocolateFlavour(
		chocolate.WithDefaults(
			&viewStyle,
			&selectedStyle,
			&focusedStyle,
		),
	)

	// the chocolate root layout with the above flavour / theme
	// We're loading the layout definition from a file "layout.cnf"
	// which is using json.
	// Have a look into that file and check the "_comment" entries
	// to get an idea on what the constraints are used for.
	choc := chocolate.NewChocolate(chocolate.WithFlavour(theme))
	if err := choc.FromFile("./layout.cnf"); err != nil {
		panic(err)
	}

	// MakeText is creating a model with a name (first, second, menuheader)
	// and place it under the specified bar (contentbar, menuheader)
	// When flavoured is set to true it will pass the flavour from the
	// parent to the model, which can be used by setStyle (will be shown later)
	// If an optional list of FlavourStyleSelector is provided it will just
	// pass the styles of the flavour to the model, that are specified
	f := choc.MakeText("first", "contentbar", true, chocolate.TS_FOCUSED)
	s := choc.MakeText("second", "contentbar", true, chocolate.TS_FOCUSED)
	mh := choc.MakeText("menuheader", "menuheader", true)

	// Set the content for the TextModel
	f.SetText("First")
	s.SetText("Second")
	mh.SetText("Main Menu")

	// create the menu model and place it as ModelBarModel (BarModel interface, including Resize)
	// with the name "menu" inside the bar "menubar" providing the flavour
	menuModel := NewMenuModel([]string{"first", "second"})
	choc.AddModelBarModel(menuModel, "menu", "menubar", true)

	// Here we can define overrides for specific styles (TS_DEFAULT) of a single model (menuheader, menu, first, second)
	// of a bar (menuheader, menubar, contentbar)
	// The available modifiers are basically identical to lipgloss.Style methods that change the style (borders, colors, etc.)
	choc.AddRootThemeModifier(chocolate.TS_DEFAULT, chocolate.Border(lipgloss.RoundedBorder()))
	choc.AddThemeModifier("menuheader", "menuheader", chocolate.TS_DEFAULT, chocolate.Border(lipgloss.DoubleBorder()), chocolate.BorderLeft(false), chocolate.BorderRight(false), chocolate.BorderTop(false))
	choc.AddThemeModifier("menubar", "menu", chocolate.TS_DEFAULT, chocolate.AlignVertical(lipgloss.Top))
	choc.AddThemeModifier("contentbar", "first", chocolate.TS_FOCUSED, chocolate.Border(lipgloss.RoundedBorder()))
	choc.AddThemeModifier("contentbar", "second", chocolate.TS_FOCUSED, chocolate.Border(lipgloss.RoundedBorder()))

	m := &model{
		choc: choc,
		menu: menuModel,
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
