package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mfulz/chocolate"
)

type testModel struct {
	base  string
	final string
}

func (t testModel) Init() tea.Cmd                            { return nil }
func (t *testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t testModel) View() string                             { return t.final }

var testModelUpdateHandler = func(b *chocolate.ChocolateBar, m tea.Model) func(tea.Msg) tea.Cmd {
	model := m.(*testModel)
	return func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "+":
				p := b.GetParent()
				if p != nil {
					switch p.GetLayout() {
					case chocolate.LIST:
						l, v := b.Y.Get()
						v++
						b.Y.Set(l, v)
						if l != chocolate.DYNAMIC {
							model.final = fmt.Sprintf("%s-%d", model.base, v)
						}
						return nil
					case chocolate.LINEAR:
						l, v := b.X.Get()
						v++
						b.X.Set(l, v)
						if l != chocolate.DYNAMIC {
							model.final = fmt.Sprintf("%s-%d", model.base, v)
						}
						return nil
					}
				}
			case "-":
				p := b.GetParent()
				if p != nil {
					switch p.GetLayout() {
					case chocolate.LIST:
						l, v := b.Y.Get()
						v--
						b.Y.Set(l, v)
						if l != chocolate.DYNAMIC {
							model.final = fmt.Sprintf("%s-%d", model.base, v)
						}
						return nil
					case chocolate.LINEAR:
						l, v := b.X.Get()
						v--
						b.X.Set(l, v)
						if l != chocolate.DYNAMIC {
							model.final = fmt.Sprintf("%s-%d", model.base, v)
						}
						return nil
					}
				}
			}
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

	fifthBar := chocolate.NewChocolateBar(nil,
		chocolate.WithModel(&chocolate.BarModel{
			Model:            fifthModel,
			UpdateHandlerFct: testModelUpdateHandler,
		}),
		chocolate.WithSelectable(),
		chocolate.WithInputOnSelect(),
		chocolate.WithYScaler(chocolate.NewParentScaler(3)),
	)
	fourthBar := chocolate.NewChocolateBar(nil,
		chocolate.WithModel(&chocolate.BarModel{
			Model:            fourthModel,
			UpdateHandlerFct: testModelUpdateHandler,
		}),
		chocolate.WithSelectable(),
		chocolate.WithInputOnSelect(),
	)
	thirdBar := chocolate.NewChocolateBar(nil,
		chocolate.WithModel(&chocolate.BarModel{
			Model:            thirdModel,
			UpdateHandlerFct: testModelUpdateHandler,
		}),
		chocolate.WithSelectable(),
		chocolate.WithInputOnSelect(),
	)
	secondBar := chocolate.NewChocolateBar(nil,
		chocolate.WithModel(&chocolate.BarModel{
			Model:            secondModel,
			UpdateHandlerFct: testModelUpdateHandler,
		}),
		chocolate.WithSelectable(),
		chocolate.WithInputOnSelect(),
		chocolate.WithXScaler(chocolate.NewDynamicScaler()),
	)
	firstBar := chocolate.NewChocolateBar(nil,
		chocolate.WithModel(&chocolate.BarModel{
			Model:            firstModel,
			UpdateHandlerFct: testModelUpdateHandler,
		}),
		chocolate.WithSelectable(),
		chocolate.WithInputOnSelect(),
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
