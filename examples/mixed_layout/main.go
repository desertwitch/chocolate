package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mfulz/chocolate"
)

type testModel string

func (t testModel) Init() tea.Cmd                           { return nil }
func (t testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t testModel) View() string                            { return string(t) }

func main() {
	// firstModel := testModel("linear-fixed-60-first")
	// secondModel := testModel("linear-dynamic-second")
	// thirdModel := testModel("linear-parent-1-list-third")
	// fourthModel := testModel("linear-parent-1-list-fourth")
	// fifthModel := testModel("linear-parent-3-list-fifth")
	// overlayModel := testModel("Overlay")

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// c, err := chocolate.NewNChocolate(chocolate.SetLayout(chocolate.LINEAR))
	c, err := chocolate.NewNChocolate(false)
	if err != nil {
		panic(err)
	}

	// firstBar := chocolate.NewModelBar(
	// 	&chocolate.BarModel{Model: firstModel},
	// 	chocolate.WithBarXScaler(chocolate.FIXED, 60),
	// 	chocolate.WithBarSelectable(),
	// 	chocolate.WithBarID("firstBar"),
	// )
	//
	// secondBar := chocolate.NewModelBar(
	// 	&chocolate.BarModel{Model: secondModel},
	// 	// chocolate.WithBarXScaler(chocolate.DYNAMIC, 0),
	// 	// chocolate.WithBarYScaler(chocolate.DYNAMIC, 0),
	// 	chocolate.WithBarSelectable(),
	// 	chocolate.WithBarID("secondBar"),
	// )
	// secondBar.Focusable(true)
	//
	// thirdBar := chocolate.NewModelBar(
	// 	&chocolate.BarModel{Model: thirdModel},
	// 	chocolate.WithBarSelectable(),
	// 	chocolate.WithBarID("thirdBar"),
	// )
	// thirdBar.Focusable(true)
	//
	// fourthBar := chocolate.NewModelBar(
	// 	&chocolate.BarModel{Model: fourthModel},
	// 	chocolate.WithBarID("fourthBar"),
	// )
	//
	// fifthBar := chocolate.NewModelBar(
	// 	&chocolate.BarModel{Model: fifthModel},
	// 	chocolate.WithBarYScaler(chocolate.PARENT, 3),
	// 	chocolate.WithBarID("fifthBar"),
	// )
	//
	// overlayBar := chocolate.NewModelBar(
	// 	&chocolate.BarModel{Model: overlayModel},
	// 	chocolate.WithBarXScaler(chocolate.FIXED, 20),
	// 	chocolate.WithBarYScaler(chocolate.FIXED, 20),
	// 	chocolate.WithBarID("overlay"),
	// )
	// // overlayBar.Hide(true)
	//
	// containerBar := chocolate.NewLayoutBar(
	// 	chocolate.LIST,
	// 	chocolate.WithBarID("container"),
	// )

	// c.AddBar("root", firstBar)
	// c.AddBar("root", secondBar)
	// c.AddBar("root", thirdBar)
	// c.AddBar("root", containerBar)
	// c.AddBar("container", thirdBar)
	// c.AddBar("container", fourthBar)
	// c.AddBar("container", fifthBar)
	// c.AddOverlayRoot(overlayBar)
	// overlayBar.Hide(true)

	if _, err := tea.NewProgram(c,
		tea.WithAltScreen()).Run(); err != nil {
		fmt.Println(err)
		log.Panic(err)
	}
}
