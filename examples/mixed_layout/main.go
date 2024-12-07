package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mfulz/chocolate"
)

type testModel string

func (t testModel) Init() tea.Cmd                           { return nil }
func (t testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t testModel) View() string                            { return string(t) }

func main() {
	firstModel := testModel("linear-fixed-60-first")
	secondModel := testModel("linear-dynamic-second")
	thirdModel := testModel("linear-parent-1-list-third")
	fourthModel := testModel("linear-parent-1-list-fourth")
	fifthModel := testModel("linear-parent-3-list-fifth")

	fifthBar := chocolate.NewChocolateBar(nil,
		chocolate.WithModel(&chocolate.BarModel{
			Model: fifthModel,
		}),
		chocolate.WithYScaler(chocolate.NewParentScaler(3)),
	)
	fourthBar := chocolate.NewChocolateBar(nil,
		chocolate.WithModel(&chocolate.BarModel{
			Model: fourthModel,
		}),
	)
	thirdBar := chocolate.NewChocolateBar(nil,
		chocolate.WithModel(&chocolate.BarModel{
			Model: thirdModel,
		}),
	)
	secondBar := chocolate.NewChocolateBar(nil,
		chocolate.WithModel(&chocolate.BarModel{
			Model: secondModel,
		}),
		chocolate.WithXScaler(chocolate.NewDynamicScaler()),
	)
	firstBar := chocolate.NewChocolateBar(nil,
		chocolate.WithModel(&chocolate.BarModel{
			Model: firstModel,
		}),
		chocolate.WithXScaler(chocolate.NewFixedScaler(60)),
	)
	containerBar := chocolate.NewChocolateBar([]*chocolate.ChocolateBar{
		thirdBar,
		fourthBar,
		fifthBar,
	},
		chocolate.WithLayout(chocolate.LIST),
	)
	bar := chocolate.NewChocolateBar([]*chocolate.ChocolateBar{
		firstBar,
		secondBar,
		containerBar,
	},
		chocolate.WithLayout(chocolate.LINEAR),
	)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	if m, err := chocolate.NewChocolate(bar); // chocolate.WithAutofocus(bar),
	err != nil {
		panic(err)
	} else {
		if _, err := tea.NewProgram(m,
			tea.WithAltScreen()).Run(); err != nil {
			fmt.Println(err)
		}
	}
}
