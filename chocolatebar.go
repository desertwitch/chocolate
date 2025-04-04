package chocolate

import (
	"fmt"
	"strings"

	"github.com/lithdew/casso"
)

type chocolateModel interface {
	barConstrainer
	barRenderer
	selectStyle(FlavourStyleSelector)
	addThemeModifier(FlavourStyleSelector, ...ThemeStyleModifier)
	setBar(*chocolateBar)
}

type chocolateBar struct {
	_width  int
	_height int
	_xpos   int
	_ypos   int
	_xend   int
	_yend   int

	canhide bool

	cElem constraintElement

	parent barContainer

	models   map[string]chocolateModel
	selected chocolateModel
	hidden   chocolateModel
	current  chocolateModel
}

func (cb *chocolateBar) update(solver *casso.Solver) {
	width := int(solver.Val(cb.cElem.width))
	height := int(solver.Val(cb.cElem.height))
	xPos := int(solver.Val(cb.cElem.xpos))
	yPos := int(solver.Val(cb.cElem.ypos))

	if cb._xpos != xPos ||
		cb._ypos != yPos {
		cb.setDirty()
	}

	cb._xpos = xPos
	cb._ypos = yPos
	cb.Resize(width, height)
}

func (cb *chocolateBar) Resize(width, height int) {
	if cb.current == nil {
		return
	}
	if cb._width != width ||
		cb._height != height {
		cb.setDirty()
	}

	cb._width = width
	cb._height = height
	cb._xend = cb._xpos + cb._width
	cb._yend = cb._ypos + cb._height
	cb.current.setSize(cb._width, cb._height)
}

func (cb *chocolateBar) width() int                    { return cb._width }
func (cb *chocolateBar) height() int                   { return cb._height }
func (cb *chocolateBar) xpos() int                     { return cb._xpos }
func (cb *chocolateBar) ypos() int                     { return cb._ypos }
func (cb *chocolateBar) xend() int                     { return cb._xend }
func (cb *chocolateBar) yend() int                     { return cb._yend }
func (cb *chocolateBar) getCelem() constraintElement   { return cb.cElem }
func (cb *chocolateBar) anyZero() bool                 { return cb._width < 1 || cb._height < 1 }
func (cb *chocolateBar) setParent(parent barContainer) { cb.parent = parent }
func (cb *chocolateBar) setDirty() {
	if cb.parent != nil {
		cb.parent.setDirty()
	}
}

func (cb *chocolateBar) getInitConstraints() []casso.Constraint {
	if cb.current == nil {
		return nil
	}
	ret := []casso.Constraint{}

	wcon, hcon := cb.current.sizeConstraints()
	for _, i := range wcon {
		c := casso.NewConstraint(casso.Op(i.Relation), i.Value, cb.cElem.width.T(1))
		ret = append(ret, c)
	}
	for _, i := range hcon {
		c := casso.NewConstraint(casso.Op(i.Relation), i.Value, cb.cElem.height.T(1))
		ret = append(ret, c)
	}

	return ret
}

func (cb *chocolateBar) SelectModel(name string) error {
	if model, ok := cb.models[strings.ToLower(name)]; !ok {
		return fmt.Errorf("invalid model '%s'", name)
	} else {
		cb.selectModel(model)
	}

	return nil
}

func (cb *chocolateBar) addModel(name string, model chocolateModel) {
	if cb.models == nil {
		cb.models = make(map[string]chocolateModel)
		defer cb.SelectModel(name)
	}

	cb.models[strings.ToLower(name)] = model
	model.setBar(cb)
}

func (cb *chocolateBar) hide() {
	if !cb.canhide {
		return
	}
	if cb.hidden == nil {
		cb.hidden = newHiddenModel()
	}

	if cb.current != cb.hidden {
		cb.setDirty()
		cb.current = cb.hidden
	}
}

func (cb *chocolateBar) unhide() {
	if cb.current != cb.selected {
		cb.setDirty()
		cb.current = cb.selected
	}
}

func (cb *chocolateBar) selectStyle(v FlavourStyleSelector) {
	if cb.current == nil {
		return
	}
	cb.current.selectStyle(v)
}

func (cb *chocolateBar) addThemeModifier(name string, style FlavourStyleSelector, modifiers ...ThemeStyleModifier) {
	if model, ok := cb.models[strings.ToLower(name)]; ok {
		model.addThemeModifier(style, modifiers...)
	}
}

func (cb *chocolateBar) View() string {
	if cb.current == nil {
		return ""
	}
	return cb.current.render()
}

func (cb *chocolateBar) setCanHide(v bool) { cb.canhide = v }

func (cb *chocolateBar) selectModel(model chocolateModel) {
	if cb.current != model {
		cb.setDirty()
		cb.current = model
		cb.selected = model
	}
}

func (cb *chocolateBar) sizeConstraints() (width, height []barSizeConstraint) {
	if cb.current == nil {
		return nil, nil
	}
	return cb.current.sizeConstraints()
}

func (cb *chocolateBar) constraintTarget(a ConstraintAttribute) bool {
	if cb.current == nil {
		return false
	}
	return cb.current.constraintTarget(a)
}

func (cb *chocolateBar) canBias() bool {
	if cb.current == nil {
		return false
	}
	return cb.current.canBias()
}

func newChocolateBar(current string, model chocolateModel, canhide bool) *chocolateBar {
	ret := &chocolateBar{
		cElem: constraintElement{
			width:  casso.New(),
			height: casso.New(),
			xpos:   casso.New(),
			ypos:   casso.New(),
		},
		canhide: canhide,
	}

	if model != nil {
		ret.addModel(current, model)
		ret.SelectModel(current)
	}

	return ret
}
