package chocolate

import (
	"os"
	"sort"

	"github.com/charmbracelet/lipgloss"
)

type Chocolate struct {
	chocolateFlavour

	bars     map[string]*chocolateBar
	overlays map[string]*Overlay

	root      *constraintLayout
	rootModel *chocolateBar
}

func (c *Chocolate) Resize(width, height int) {
	c.rootModel.Resize(width, height)
	for _, o := range c.overlays {
		o.Resize(width, height)
	}
}

func (c *Chocolate) setDirty()                          { c.root.setDirty() }
func (c *Chocolate) addBar(name string, child barChild) { c.root.addBar(name, child) }

func (c *Chocolate) View() string {
	overlaysSorted := []*Overlay{}
	ret := c.rootModel.View()
	w := lipgloss.Width(ret)
	h := lipgloss.Height(ret)
	for _, o := range c.overlays {
		if o.enabled {
			overlaysSorted = append(overlaysSorted, o)
		}
	}
	sort.Slice(overlaysSorted, func(i, j int) bool {
		return overlaysSorted[i].zindex < overlaysSorted[j].zindex
	})

	for _, o := range overlaysSorted {
		oview := o.View()
		ow := lipgloss.Width(oview)
		oh := lipgloss.Height(oview)
		ox, oy := o.calcPosition(w, h, ow, oh)
		ret = placeOverlay(ox, oy, oview, ret)
	}
	return ret
}

func (c *Chocolate) AddConstraint(con Constraint) {
	c.root.AddConstraint(con)
	c.MakeBar(con.Target, false)
}

func (c *Chocolate) FromFile(file string) error {
	if layout, err := os.ReadFile(file); err == nil {
		return c.FromJson(layout)
	} else {
		return err
	}
}

func (c *Chocolate) FromJson(layout []byte) error {
	if err := c.root.fromJson(layout); err != nil {
		return err
	}

	for _, con := range c.root.constraints {
		c.MakeBar(con.Target, false)
	}

	return nil
}

func (c *Chocolate) AddThemeModifier(name string, model string, style FlavourStyleSelector, modifiers ...ThemeStyleModifier) {
	if b, ok := c.bars[name]; ok {
		b.addThemeModifier(model, style, modifiers...)
	}
}

func (c *Chocolate) MakeBar(name string, canhide bool) {
	if _, ok := c.bars[name]; ok {
		return
	}
	bar := newChocolateBar(
		"",
		nil,
		canhide,
	)
	c.addBar(name, bar)

	if c.bars == nil {
		c.bars = make(map[string]*chocolateBar)
	}
	c.bars[name] = bar
}

func (c *Chocolate) MakeChocolate(name string, bar string, flavoured bool, styles ...FlavourStyleSelector) *Chocolate {
	b, ok := c.bars[bar]
	if !ok {
		return nil
	}
	var model *chocolateBarModel[*Chocolate]
	if flavoured {
		model = newFlavouredModelBarModel(NewChocolate(WithFlavour(&c.chocolateFlavour)), &c.chocolateFlavour, styles...)
	} else {
		model = newModelBarModel(NewChocolate(WithFlavour(&c.chocolateFlavour)))
	}
	b.addModel(name, model)
	// b.SelectModel(name)

	return model.model()
}

func (c *Chocolate) MakeOverlay(
	name string,
	zindex int,
	width float64, height float64,
	flavoured bool,
	pos ...OverlayPosition,
) *Overlay {
	if o, ok := c.overlays[name]; ok {
		return o
	}
	if c.overlays == nil {
		c.overlays = make(map[string]*Overlay)
	}
	var choc *Chocolate
	if flavoured {
		choc = NewChocolate(WithFlavour(&c.chocolateFlavour))
	} else {
		choc = NewChocolate()
	}
	o := newOverlay(choc, zindex, width, height, pos...)
	c.overlays[name] = o

	return o
}

func (c *Chocolate) MakeText(name string, bar string, flavoured bool, styles ...FlavourStyleSelector) *TextModel {
	b, ok := c.bars[bar]
	if !ok {
		return nil
	}
	var model *chocolateBarModel[*TextModel]
	if flavoured {
		model = newFlavouredTextBarModel("", &c.chocolateFlavour, styles...)
	} else {
		model = newTextBarModel("")
	}
	b.addModel(name, model)
	// b.SelectModel(name)

	return model.model()
}

func (c *Chocolate) MakeStyledText(name string, bar string, style *lipgloss.Style) *TextModel {
	b, ok := c.bars[bar]
	if !ok {
		return nil
	}
	model := newStyledTextBarModel("", style)
	b.addModel(name, model)
	// b.SelectModel(name)

	return model.model()
}

func (c *Chocolate) AddViewBarModel(model barViewer, name string, bar string, flavoured bool, styles ...FlavourStyleSelector) {
	b, ok := c.bars[bar]
	if !ok {
		return
	}
	var _model *chocolateBarModel[barViewer]
	if flavoured {
		_model = newFlavouredViewBarModel(model, &c.chocolateFlavour, styles...)
	} else {
		_model = newViewBarModel(model)
	}
	b.addModel(name, _model)
	// b.SelectModel(name)
}

func (c *Chocolate) AddStyledViewBarModel(model barViewer, name string, bar string, style *lipgloss.Style) {
	b, ok := c.bars[bar]
	if !ok {
		return
	}
	_model := newStyledViewBarModel(model, style)
	b.addModel(name, _model)
	// b.SelectModel(name)
}

func (c *Chocolate) SelectModel(name string, bar string) {
	if b, ok := c.bars[bar]; ok {
		b.SelectModel(name)
	}
}

func (c *Chocolate) SelectStyle(name FlavourStyleSelector, bar string) {
	if b, ok := c.bars[bar]; ok {
		b.selectStyle(name)
	}
}

func (c *Chocolate) SelectRootStyle(name FlavourStyleSelector) {
	c.rootModel.selectStyle(name)
}

func (c *Chocolate) SetCanHide(bar string, v bool) {
	if b, ok := c.bars[bar]; ok {
		b.setCanHide(v)
	}
}

func (c *Chocolate) Hide(bar string) {
	if b, ok := c.bars[bar]; ok {
		b.hide()
	}
}

func (c *Chocolate) Unhide(bar string) {
	if b, ok := c.bars[bar]; ok {
		b.unhide()
	}
}

type ChocolateOption func(*Chocolate)

func WithFlavour(flavour *chocolateFlavour) ChocolateOption {
	return func(c *Chocolate) {
		c.chocolateFlavour = *flavour
	}
}

func NewChocolate(opts ...ChocolateOption) *Chocolate {
	ret := &Chocolate{
		chocolateFlavour: *NewChocolateFlavour(),
		root:             newConstraintLayout(),
	}

	for _, opt := range opts {
		opt(ret)
	}

	ret.rootModel = newChocolateBar("default",
		newFlavouredModelBarModel(ret.root, &ret.chocolateFlavour),
		false,
	)
	ret.rootModel.setParent(ret)

	return ret
}
