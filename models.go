package chocolate

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type chocolateBarModel[T any] struct {
	bar *chocolateBar
	barRenderer
	barConstrainer

	styles        map[FlavourStyleSelector]*lipgloss.Style
	styleModifier map[FlavourStyleSelector][]ThemeStyleModifier

	current  *lipgloss.Style
	selected FlavourStyleSelector

	srcModel T
}

func (cbm *chocolateBarModel[T]) setBar(v *chocolateBar) { cbm.bar = v }
func (cbm *chocolateBarModel[T]) model() T               { return cbm.srcModel }

func (cbm *chocolateBarModel[T]) setDirty() {
	if cbm.bar != nil {
		cbm.bar.setDirty()
	}
}

func (cbm *chocolateBarModel[T]) selectStyle(style FlavourStyleSelector) {
	s := FlavourStyleSelector(strings.ToLower(string(style)))
	if sel, ok := cbm.styles[s]; ok {
		if cbm.current.GetHorizontalFrameSize() != sel.GetHorizontalFrameSize() ||
			cbm.current.GetVerticalFrameSize() != sel.GetVerticalFrameSize() {
			cbm.setDirty()
		}
		*cbm.current = *sel
		cbm.selected = s

		if selMod, ok := cbm.styleModifier[s]; ok {
			for _, mod := range selMod {
				*cbm.current = mod(*cbm.current)
			}
		}
	}
}

func (cbm *chocolateBarModel[T]) addThemeModifier(style FlavourStyleSelector, modifiers ...ThemeStyleModifier) {
	if len(modifiers) <= 0 {
		return
	}
	s := FlavourStyleSelector(strings.ToLower(string(style)))
	if _, ok := cbm.styles[s]; ok {
		if cbm.styleModifier == nil {
			cbm.styleModifier = make(map[FlavourStyleSelector][]ThemeStyleModifier)
		}
		if cbm.styleModifier[s] == nil {
			cbm.styleModifier[s] = make([]ThemeStyleModifier, 0)
		}
		cbm.styleModifier[s] = append(cbm.styleModifier[s], modifiers...)
		if cbm.selected == style {
			cbm.setDirty()
			for _, mod := range modifiers {
				*cbm.current = mod(*cbm.current)
			}
		}
	}
}

func newChocolateBarModel[T any](
	model T,
	constrainer barConstrainer,
	renderer barRenderer,
	styles map[FlavourStyleSelector]*lipgloss.Style,
	current *lipgloss.Style,
	selected FlavourStyleSelector,
) *chocolateBarModel[T] {
	return &chocolateBarModel[T]{
		barConstrainer: constrainer,
		barRenderer:    renderer,
		srcModel:       model,
		styles:         styles,
		current:        current,
		selected:       selected,
	}
}

func newHiddenModel() *chocolateBarModel[string] {
	return newChocolateBarModel(
		"",
		newHiddenConstrainer(),
		newNoneRenderer(),
		nil, nil, "",
	)
}

type TextModel struct {
	bar  *chocolateBarModel[*TextModel]
	text string
}

func (tm *TextModel) SetText(v string) {
	if tm == nil {
		return
	}
	tm.text = v
	if tm.bar != nil {
		tm.bar.setDirty()
	}
}

func newTextBarModel(text string) *chocolateBarModel[*TextModel] {
	tm := &TextModel{
		text: text,
	}

	ret := newChocolateBarModel(
		tm,
		newStyledConstrainer(nil, &tm.text),
		newStaticRenderer(&tm.text),
		nil, nil, "",
	)
	tm.bar = ret

	return ret
}

func newStyledTextBarModel(text string, style *lipgloss.Style) *chocolateBarModel[*TextModel] {
	tm := &TextModel{
		text: text,
	}

	ret := newChocolateBarModel(
		tm,
		newStyledConstrainer(style, &tm.text),
		newStyleRenderer(&tm.text, style),
		nil, nil, "",
	)
	tm.bar = ret

	return ret
}

func newFlavouredTextBarModel(text string, flavour *chocolateFlavour, styles ...FlavourStyleSelector) *chocolateBarModel[*TextModel] {
	s, c, sel := flavour.getStyles(styles...)
	tm := &TextModel{
		text: text,
	}

	ret := newChocolateBarModel(
		tm,
		newStyledConstrainer(c, &tm.text),
		newStyleRenderer(&tm.text, c),
		s, c, sel,
	)
	tm.bar = ret

	return ret
}

func newViewBarModel[T BarViewer](model T) *chocolateBarModel[T] {
	vr := newViewRenderer(model, nil)
	ret := newChocolateBarModel(
		model,
		newStyledConstrainer(nil, vr.content),
		vr,
		nil, nil, "",
	)
	vr.bar = ret

	return ret
}

func newStyledViewBarModel[T BarViewer](model T, style *lipgloss.Style) *chocolateBarModel[T] {
	vr := newViewRenderer(model, style)
	ret := newChocolateBarModel(
		model,
		newStyledConstrainer(style, vr.content),
		vr,
		nil, nil, "",
	)
	vr.bar = ret

	return ret
}

func newFlavouredViewBarModel[T BarViewer](model T, flavour *chocolateFlavour, styles ...FlavourStyleSelector) *chocolateBarModel[T] {
	s, c, sel := flavour.getStyles(styles...)
	vr := newViewRenderer(model, c)
	ret := newChocolateBarModel(
		model,
		newStyledConstrainer(c, vr.content),
		vr,
		s, c, sel,
	)
	vr.bar = ret

	return ret
}

func newModelBarModel[T BarModel](model T) *chocolateBarModel[T] {
	return newChocolateBarModel(
		model,
		newStyledConstrainer(nil),
		newModelRenderer(model, nil),
		nil, nil, "",
	)
}

func newStyledModelBarModel[T BarModel](model T, style *lipgloss.Style) *chocolateBarModel[T] {
	return newChocolateBarModel(
		model,
		newStyledConstrainer(style),
		newModelRenderer(model, style),
		nil, nil, "",
	)
}

func newFlavouredModelBarModel[T BarModel](model T, flavour *chocolateFlavour, styles ...FlavourStyleSelector) *chocolateBarModel[T] {
	s, c, sel := flavour.getStyles(styles...)
	return newChocolateBarModel(
		model,
		newStyledConstrainer(c),
		newModelRenderer(model, c),
		s, c, sel,
	)
}

type teaModel struct {
	tea.Model
}

func (tmbm *teaModel) Resize(width, height int) {
	tmbm.Update(
		tea.WindowSizeMsg{
			Width:  width,
			Height: height,
		},
	)
}

func newTeaModel(model tea.Model) *teaModel {
	return &teaModel{
		model,
	}
}
