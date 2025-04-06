package main

import (
	"strings"

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

type buttonsModel struct {
	buttons     []string
	selected    int
	choice      string
	width       int
	buttonWidth int
}

func (m *buttonsModel) Init() tea.Cmd { return nil }

// Resize is used by chocolate.BarModel and will be called
// after layout calculation when the size changed
func (m *buttonsModel) Resize(width, height int) {
	m.width = width
}

func (m *buttonsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.choice = ""
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "l":
			m.selected++
			if m.selected >= len(m.buttons) {
				m.selected = 0
			}
		case "h":
			m.selected--
			if m.selected < 0 {
				m.selected = len(m.buttons) - 1
			}
		case "enter":
			m.choice = strings.ToLower(m.buttons[m.selected])
		}
	}

	return m, nil
}

func (m *buttonsModel) View() string {
	entries := []string{}
	vs := viewStyle.Width(m.buttonWidth).PaddingRight(5)
	ss := selectedStyle.Width(m.buttonWidth).PaddingRight(5)
	for i, e := range m.buttons {
		if i == m.selected {
			entries = append(entries, ss.Render(e))
		} else {
			entries = append(entries, vs.Render(e))
		}
	}

	return vs.Width(m.width).
		AlignHorizontal(lipgloss.Right).
		Render(lipgloss.JoinHorizontal(0, entries...))
}

func NewButtonsModel(items []string) *buttonsModel {
	return &buttonsModel{
		buttons:     items,
		selected:    0,
		width:       50,
		buttonWidth: 10,
	}
}

// tea.Model and program start
// it holds a reference to the main chocolate
// which acts as the central control interface
// to the layout and styling functionality
// provided by chocolate
type model struct {
	choc          *chocolate.Chocolate
	exitModel     *buttonsModel
	dialogOverlay *chocolate.Overlay
	dialogActive  bool
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
		return m, nil
	case tea.KeyMsg:
		if m.dialogActive {
			_, cmd = m.exitModel.Update(msg)
			cmds = append(cmds, cmd)
			switch msg.String() {
			case "enter":
				if m.exitModel.choice == "yes" {
					return m, tea.Quit
				}
				m.dialogOverlay.Disable()
				m.dialogActive = false
			case "escape":
				m.dialogOverlay.Disable()
				m.dialogActive = false
			}
		} else {
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				m.dialogOverlay.Enable()
				m.dialogActive = true
			}
		}
	}

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
	content := choc.MakeText("content", "mainbar", true, chocolate.TS_FOCUSED)

	// Set the content for the TextModel
	content.SetText("This is just a stupid placeholder text")

	// create an overlay container named dialog which is placed
	// in the mid with an smaller size (parent width - 20, height - 40) for width and height
	// using the flavour of the parent
	overlay := choc.MakeOverlay("overlay", 1, -20, -40, true, chocolate.CENTER, chocolate.CENTER)

	// We're loading the layout definition from a file "overlay.cnf"
	// which is using json.
	// Have a look into that file and check the "_comment" entries
	// to get an idea on what the constraints are used for.
	if err := overlay.FromFile("./dialog.cnf"); err != nil {
		panic(err)
	}

	// create a TextModel which is providing the dialog question
	question := overlay.MakeText("question", "contentbar", true)
	question.SetText("Do you really want to quit?")

	// create the dialog buttons
	buttons := NewButtonsModel([]string{"Yes", "No"})

	// add the buttons to the overlay
	overlay.AddModelBarModel(buttons, "buttons", "buttonbar", true)

	// Here we can define overrides for specific styles (TS_DEFAULT) of a single model (menuheader, menu, first, second)
	// of a bar (menuheader, menubar, contentbar)
	// The available modifiers are basically identical to lipgloss.Style methods that change the style (borders, colors, etc.)
	choc.AddRootThemeModifier(chocolate.TS_DEFAULT, chocolate.Border(lipgloss.RoundedBorder()))
	overlay.AddRootThemeModifier(chocolate.TS_DEFAULT, chocolate.Border(lipgloss.RoundedBorder()))
	overlay.AddThemeModifier("contentbar", "question", chocolate.TS_DEFAULT, chocolate.Border(lipgloss.RoundedBorder()))

	m := &model{
		choc:          choc,
		exitModel:     buttons,
		dialogOverlay: overlay,
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
