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

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	c, err := chocolate.NewNChocolate(chocolate.SetLayout(chocolate.LINEAR))
	if err != nil {
		panic(err)
	}

	firstBar := chocolate.NewModelBar(
		&chocolate.ModelBarModel{Model: firstModel},
		chocolate.ModelBarXScaler(chocolate.FIXED, 60),
		chocolate.ModelBarSelectable(),
	)

	secondBar := chocolate.NewModelBar(
		&chocolate.ModelBarModel{Model: secondModel},
		chocolate.ModelBarXScaler(chocolate.DYNAMIC, 0),
		chocolate.ModelBarSelectable(),
	)

	thirdBar := chocolate.NewModelBar(
		&chocolate.ModelBarModel{Model: thirdModel},
		chocolate.ModelBarSelectable(),
	)

	fourthBar := chocolate.NewModelBar(
		&chocolate.ModelBarModel{Model: fourthModel},
	)

	fifthBar := chocolate.NewModelBar(
		&chocolate.ModelBarModel{Model: fifthModel},
		chocolate.ModelBarYScaler(chocolate.PARENT, 3),
		chocolate.ModelBarID("fifthBar"),
	)

	containerBar := chocolate.NewLayoutBar(
		chocolate.LIST,
		chocolate.LayoutBarID("container"),
	)

	c.AddBar("root", firstBar)
	c.AddBar("root", secondBar)
	c.AddBar("root", containerBar)
	c.AddBar("container", thirdBar)
	c.AddBar("container", fourthBar)
	c.AddBar("container", fifthBar)

	if _, err := tea.NewProgram(c,
		tea.WithAltScreen()).Run(); err != nil {
		fmt.Println(err)
	}
}
