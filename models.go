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

	current *lipgloss.Style

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
	}
}

func newChocolateBarModel[T any](
	model T,
	constrainer barConstrainer,
	renderer barRenderer,
	styles map[FlavourStyleSelector]*lipgloss.Style,
	current *lipgloss.Style,
) *chocolateBarModel[T] {
	return &chocolateBarModel[T]{
		barConstrainer: constrainer,
		barRenderer:    renderer,
		srcModel:       model,
		styles:         styles,
		current:        current,
	}
}

func newHiddenModel() *chocolateBarModel[string] {
	return newChocolateBarModel(
		"",
		newHiddenConstrainer(),
		newNoneRenderer(),
		nil, nil,
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

type teaModel[T any] struct {
	bar *chocolateBarModel[tea.Model]
}

func (tm *teaModel[T]) Resize(width, height int) {
	if tm == nil || tm.bar == nil {
		return
	}
	tm.bar.srcModel, _ = tm.bar.model().Update(tea.WindowSizeMsg{Width: width, Height: height})
}

func newTextBarModel(text string) *chocolateBarModel[*TextModel] {
	tm := &TextModel{
		text: text,
	}

	ret := newChocolateBarModel(
		tm,
		newStyledConstrainer(nil, &tm.text),
		newStaticRenderer(&tm.text),
		nil, nil,
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
		nil, nil,
	)
	tm.bar = ret

	return ret
}

func newFlavouredTextBarModel(text string, flavour *chocolateFlavour, styles ...FlavourStyleSelector) *chocolateBarModel[*TextModel] {
	s, c := flavour.getStyles(styles...)
	tm := &TextModel{
		text: text,
	}

	ret := newChocolateBarModel(
		tm,
		newStyledConstrainer(c, &tm.text),
		newStyleRenderer(&tm.text, c),
		s, c,
	)
	tm.bar = ret

	return ret
}

func newViewBarModel[T barViewer](model T) *chocolateBarModel[T] {
	vr := newViewRenderer(model, nil)
	ret := newChocolateBarModel(
		model,
		newStyledConstrainer(nil, vr.content),
		vr,
		nil, nil,
	)
	vr.bar = ret

	return ret
}

func newStyledViewBarModel[T barViewer](model T, style *lipgloss.Style) *chocolateBarModel[T] {
	vr := newViewRenderer(model, style)
	ret := newChocolateBarModel(
		model,
		newStyledConstrainer(style, vr.content),
		vr,
		nil, nil,
	)
	vr.bar = ret

	return ret
}

func newFlavouredViewBarModel[T barViewer](model T, flavour *chocolateFlavour, styles ...FlavourStyleSelector) *chocolateBarModel[T] {
	s, c := flavour.getStyles(styles...)
	vr := newViewRenderer(model, c)
	ret := newChocolateBarModel(
		model,
		newStyledConstrainer(c, vr.content),
		vr,
		s, c,
	)
	vr.bar = ret

	return ret
}

func newModelBarModel[T barModel](model T) *chocolateBarModel[T] {
	return newChocolateBarModel(
		model,
		newStyledConstrainer(nil),
		newModelRenderer(model, nil),
		nil, nil,
	)
}

func newStyledModelBarModel[T barModel](model T, style *lipgloss.Style) *chocolateBarModel[T] {
	return newChocolateBarModel(
		model,
		newStyledConstrainer(style),
		newModelRenderer(model, style),
		nil, nil,
	)
}

func newFlavouredModelBarModel[T barModel](model T, flavour *chocolateFlavour, styles ...FlavourStyleSelector) *chocolateBarModel[T] {
	s, c := flavour.getStyles(styles...)
	return newChocolateBarModel(
		model,
		newStyledConstrainer(c),
		newModelRenderer(model, c),
		s, c,
	)
}
