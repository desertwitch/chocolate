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
			BorderForeground(lipgloss.Color("246")).BorderBackground(lipgloss.Color("232")).
			Border(lipgloss.RoundedBorder())
	selectedStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("15")).Background(lipgloss.Color("237")).
			BorderForeground(lipgloss.Color("15")).BorderBackground(lipgloss.Color("237")).
			Border(lipgloss.RoundedBorder())
	focusedStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("15")).Background(lipgloss.Color("232")).
			BorderForeground(lipgloss.Color("196")).BorderBackground(lipgloss.Color("232")).
			Border(lipgloss.RoundedBorder())
)

type textModel struct {
	text string
	alt  string

	cur *string
}

func (t *textModel) Init() tea.Cmd { return nil }
func (t *textModel) View() string  { return *t.cur }
func (t *textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "t":
			if t.cur == &t.alt {
				t.cur = &t.text
			} else {
				t.cur = &t.alt
			}
		case "a":
			t.cur = &t.alt
		case "s":
			t.cur = &t.text
		}
	}

	return t, nil
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
	choc := chocolate.NewChocolateTui(chocolate.WithTuiFlavour(theme))
	if err := choc.FromFile("./layout.cnf"); err != nil {
		panic(err)
	}

	tui1 := &textModel{
		text: "Tui1 Text",
		alt:  "Tui1 Alt Text",
	}
	tui1.cur = &tui1.text

	tui2 := &textModel{
		text: "Tui2 Text",
		alt:  "Tui2 Alt Text",
	}
	tui2.cur = &tui2.text

	tui3 := &textModel{
		text: "Tui3 Text",
		alt:  "Tui3 Alt Text",
	}
	tui3.cur = &tui3.text

	// MakeText is creating a model with a name (first, second, menuheader)
	// and place it under the specified bar (contentbar, menuheader)
	// When flavoured is set to true it will pass the flavour from the
	// parent to the model, which can be used by setStyle (will be shown later)
	// If an optional list of FlavourStyleSelector is provided it will just
	// pass the styles of the flavour to the model, that are specified
	choc.AddTuiModel(tui1, "tui1", "tuimodel1",
		chocolate.EnableAutoFocus(),
	)
	choc.AddTuiModel(tui2, "tui2", "tuimodel2",
		chocolate.EnableFocus(),
		chocolate.EnableNav(),
		chocolate.SetPrev("tui1"),
	)
	// choc.AddTuiModel(tui3, "tui3", "tuimodel3",
	// 	chocolate.EnableSelect(),
	// 	chocolate.SetPrev("tui2"),
	// 	chocolate.SetNext("tui1"),
	// )

	subchoc := choc.MakeChocolateTui("tui3", "tuimodel3",
		chocolate.EnableAutoFocus(),
		chocolate.SetPrev("tui2"),
		chocolate.SetNext("tui1"),
	)
	if err := subchoc.FromFile("./layout.cnf"); err != nil {
		panic(err)
	}

	tui21 := &textModel{
		text: "Tui21 Text",
		alt:  "Tui21 Alt Text",
	}
	tui21.cur = &tui21.text

	tui22 := &textModel{
		text: "Tui22 Text",
		alt:  "Tui22 Alt Text",
	}
	tui22.cur = &tui22.text

	tui23 := &textModel{
		text: "Tui23 Text",
		alt:  "Tui23 Alt Text",
	}
	tui23.cur = &tui23.text

	subchoc.AddTuiModel(tui21, "tui1", "tuimodel1",
		chocolate.EnableAutoFocus(),
	)
	subchoc.AddTuiModel(tui22, "tui2", "tuimodel2",
		chocolate.EnableFocus(),
		chocolate.EnableNav(),
		chocolate.SetPrev("tui1"),
	)
	subchoc.AddTuiModel(tui23, "tui3", "tuimodel3",
		chocolate.EnableSelect(),
		chocolate.SetPrev("tui2"),
		chocolate.SetNext("tui1"),
	)

	p := tea.NewProgram(choc)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
