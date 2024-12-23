package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mfulz/chocolate"
	"github.com/mfulz/chocolate/sweets/notify"
)

type testModel struct {
	base  string
	final string
}

func (t testModel) Init() tea.Cmd                            { return nil }
func (t *testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t testModel) View() string                             { return t.final }

var testModelUpdateHandler = func(b chocolate.ChocolateBar, m tea.Model) func(tea.Msg) tea.Cmd {
	return func(msg tea.Msg) tea.Cmd {
		var cmds []tea.Cmd
		model := m.(*testModel)
		var v int
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "+":
				v = chocolate.IncLayoutSize(b)
				cmds = append(cmds, notify.NewNotifyMsg(notify.LEVEL_INFO, fmt.Sprintf("Increased size to: %d", v), time.Second*2))
			case "-":
				v = chocolate.DecLayoutSize(b)
				cmds = append(cmds, notify.NewNotifyMsg(notify.LEVEL_WARN, fmt.Sprintf("Decreased size to: %d", v), time.Second*2))
			}
		}
		if v > 0 {
			model.final = fmt.Sprintf("%s-%d", model.base, v)
			return tea.Batch(cmds...)
		}
		return nil
	}
}

func main() {
	firstModel := &testModel{base: "linear-fixed-first", final: "linear-fixed-first-60"}
	secondModel := &testModel{base: "linear-dynamic-second", final: "linear-dynamic-second"}
	thirdModel := &testModel{base: "linear-parent-list-third", final: "linear-parent-list-third-1"}
	fourthModel := &testModel{base: "linear-parent-list-fourth", final: "linear-parent-list-fourth-1"}
	fifthModel := &testModel{base: "linear-parent-list-fifth", final: "linear-parent-list-fifth-3"}

	firstBar := chocolate.NewModelBar(
		&chocolate.BarModel{
			Model:            firstModel,
			UpdateHandlerFct: testModelUpdateHandler,
		},
		chocolate.WithBarXScaler(chocolate.FIXED, 60),
		chocolate.WithBarSelectable(),
	)

	secondBar := chocolate.NewModelBar(
		&chocolate.BarModel{
			Model:            secondModel,
			UpdateHandlerFct: testModelUpdateHandler,
		},
		chocolate.WithBarXScaler(chocolate.DYNAMIC, 0),
		chocolate.WithBarSelectable(),
	)

	thirdBar := chocolate.NewModelBar(
		&chocolate.BarModel{
			Model:            thirdModel,
			UpdateHandlerFct: testModelUpdateHandler,
		},
		chocolate.WithBarSelectable(),
	)

	fourthBar := chocolate.NewModelBar(
		&chocolate.BarModel{
			Model:            fourthModel,
			UpdateHandlerFct: testModelUpdateHandler,
		},
		chocolate.WithBarSelectable(),
	)

	fifthBar := chocolate.NewModelBar(
		&chocolate.BarModel{
			Model:            fifthModel,
			UpdateHandlerFct: testModelUpdateHandler,
		},
		chocolate.WithBarXScaler(chocolate.PARENT, 3),
		chocolate.WithBarSelectable(),
	)

	containerBar := chocolate.NewLayoutBar(
		chocolate.LIST,
		chocolate.WithBarID("container"),
	)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	if m, err := chocolate.NewNChocolate(
		chocolate.SetLayout(chocolate.LINEAR),
	); // chocolate.WithAutofocus(bar),
	err != nil {
		panic(err)
	} else {
		m.AddBar("root", firstBar)
		m.AddBar("root", secondBar)
		m.AddBar("root", containerBar)
		m.AddBar("container", thirdBar)
		m.AddBar("container", fourthBar)
		m.AddBar("container", fifthBar)
		notify.NewNotificationBar(m, false)

		if _, err := tea.NewProgram(m,
			tea.WithAltScreen()).Run(); err != nil {
			fmt.Println(err)
		}
	}
}
