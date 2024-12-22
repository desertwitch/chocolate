package chocolate

//
// import (
// 	"log"
//
// 	"github.com/charmbracelet/lipgloss"
// 	"github.com/google/uuid"
// 	"github.com/mfulz/chocolate/flavour"
// )
//
// type BarStyleCustomizeHanleFct func(ChocolateBar, lipgloss.Style) func() lipgloss.Style
//
// // chocolateBar is the main workhorse that will
// // provide most of the functionality and is doing
// // all the calculations and handling of the layout
// // it further holds the tea.Models and wrap the
// // calls around, so that it acts at the end just
// // a view container
// type baseBar struct {
// 	Scaling
// 	id string
//
// 	// backref to the Chocolate the bar
// 	// belongs to
// 	// this is used to centralize some
// 	// parts for the unified theming
// 	// and controls like selector
// 	chocolate *NChocolate
//
// 	// possible maximum content size
// 	// This is used to have a maximum for the content
// 	// the real size will be calculated during the
// 	// view rendering as it is the only possible
// 	// place to handle the scaling which depends on
// 	// possible dynmic content
// 	maxWidth  int
// 	maxHeight int
//
// 	// real size after calculation without margins
// 	width  int
// 	height int
//
// 	// actual content size of children
// 	contentWidth  int
// 	contentHeight int
//
// 	// pre rendered view with maximum content sizes
// 	// this is used to get the correct sizes of the
// 	// models view to be used for dynamic scaling
// 	// and the following calculations
// 	preRendered bool
//
// 	// rendered
// 	view     string
// 	rendered bool
//
// 	// flavourPrefs generation function
// 	// this can be used to override the default
// 	// flavour preferences
// 	StyleCustomizeHandler BarStyleCustomizeHanleFct
//
// 	// custom update function
// 	// this can be used to override the default
// 	// behavior which will only let the bar
// 	// take input focus when a model is attached
// 	// and just pass the tea messages through
// 	// UpdateHandlerFct func(*ChocolateBar) func(tea.Msg) tea.Cmd
//
// 	// if the bar is hidden
// 	// hidden bars are removed from the layout
// 	// rendering and the space is used for the
// 	// other bars
// 	hidden bool
//
// 	// if this bar can be selected
// 	selectable bool
// 	// if this bar should receive input when
// 	// selected
// 	inputOnSelect bool
// }
//
// func (b *baseBar) Has(v interface{}) bool {
// 	checks, ok := v.(*checkAttributes)
// 	if !ok {
// 		return false
// 	}
//
// 	for k, c := range checks.checks {
// 		switch k {
// 		case CA_X_SCALING:
// 			delete(checks.checks, CA_X_SCALING)
// 			if !b.X.Is(c.(ScalingType)) {
// 				return false
// 			}
// 		case CA_Y_SCALING:
// 			delete(checks.checks, CA_Y_SCALING)
// 			if !b.Y.Is(c.(ScalingType)) {
// 				return false
// 			}
// 		case CA_CAN_SELECT:
// 			delete(checks.checks, CA_CAN_SELECT)
// 			if !b.selectable {
// 				return false
// 			}
// 		case CA_INPUT_ON_SELECT:
// 			delete(checks.checks, CA_INPUT_ON_SELECT)
// 			if !b.inputOnSelect {
// 				return false
// 			}
// 		default:
// 			return false
// 		}
// 	}
// 	return true
// }
//
// func (b baseBar) GetID() string {
// 	return b.id
// }
//
// func (b *baseBar) GetStyle() lipgloss.Style {
// 	ret := flavour.GetPresetNoErr(flavour.PRESET_PRIMARY_NOBORDER)
//
// 	// root
// 	if b.getChocolate().IsRoot(b) {
// 		ret = flavour.GetPresetNoErr(flavour.PRESET_PRIMARY)
// 	}
//
// 	// selected and not root
// 	if b.getChocolate().IsSelected(b) && !b.getChocolate().IsRoot(b) {
// 		ret = ret.BorderForeground(flavour.GetColorNoErr(flavour.COLOR_SECONDARY))
// 	}
//
// 	// focused and not root
// 	if b.getChocolate().HasFocus(b) && !b.getChocolate().IsRoot(b) {
// 		ret = flavour.GetPresetNoErr(flavour.PRESET_SECONDARY).
// 			BorderBackground(flavour.GetColorNoErr(flavour.COLOR_PRIMARY_BG))
// 	}
//
// 	if b.StyleCustomizeHandler != nil {
// 		ret = b.StyleCustomizeHandler(b, ret)()
// 	}
//
// 	return ret
// }
//
// func (b *baseBar) Hide(v bool) {
// 	b.hidden = v
// }
//
// func (b *baseBar) Resize(w, h int) {
// 	// if there is a frame set for the bar
// 	// this has to be removed from the available
// 	// content size
// 	width := w - b.GetStyle().GetHorizontalFrameSize()
// 	height := h - b.GetStyle().GetVerticalFrameSize()
//
// 	// if this is a fixed scaling than we don't have
// 	// to calculate anything
// 	if b.X.IsFixed() {
// 		width = b.X.GetValue()
// 		b.width = width
// 	}
// 	if b.Y.IsFixed() {
// 		height = b.Y.GetValue()
// 		b.height = height
// 	}
//
// 	b.maxWidth = width
// 	b.maxHeight = height
//
// 	// the root bar doesn't have to rescale itself
// 	if b.getChocolate().IsRoot(b) {
// 		b.width = width
// 		b.height = height
// 	}
// }
//
// func (b baseBar) GetScaling() (scaling Scaling) {
// 	return b.Scaling
// }
//
// func (b baseBar) getChocolate() *NChocolate {
// 	return b.chocolate
// }
//
// func (b *baseBar) setChocolate(v *NChocolate) {
// 	b.chocolate = v
// }
//
// func (b *baseBar) getParent() (parent LayoutBar) {
// 	return b.getChocolate().GetParent(b)
// }
//
// func (b *baseBar) setID(v string) {
// 	b.id = v
// }
//
// func (b *baseBar) setXScaler(v Scaler) {
// 	b.X = v
// }
//
// func (b *baseBar) setYScaler(v Scaler) {
// 	b.Y = v
// }
//
// func (b *baseBar) setSelectable() {
// 	b.selectable = true
// }
//
// func (b *baseBar) setInputOnSelect() {
// 	b.inputOnSelect = true
// }
//
// func (b baseBar) getContentSize() (int, int) {
// 	return b.contentWidth, b.contentHeight
// }
//
// func (b *baseBar) setWidth(width int) {
// 	b.width = width - b.GetStyle().GetHorizontalFrameSize()
// 	log.Printf("width: %d bwidth: %d\n", width, b.width)
// }
//
// func (b *baseBar) setHeight(height int) {
// 	b.height = height - b.GetStyle().GetVerticalFrameSize()
// }
//
// func (b baseBar) getSize() (int, int) {
// 	return b.width, b.height
// }
//
// func (b baseBar) getMaxSize() (width, height int) {
// 	return b.maxWidth, b.maxHeight
// }
//
// func (b *baseBar) getChildren() (children []ChocolateBar) {
// 	return b.getChocolate().GetChildren(b)
// }
//
// func (b *baseBar) preRender() bool {
// 	if b.hidden ||
// 		b.preRendered {
// 		return true
// 	}
// 	return false
// }
//
// func (b *baseBar) finalizeSizing() {
// 	pbar := b.getParent()
// 	if pbar == nil || b.isHidden() {
// 		return
// 	}
//
// 	pw, ph := pbar.getSize()
// 	pmw, pmh := pbar.getMaxSize()
// 	if pw <= 0 {
// 		pw = pmw
// 	}
// 	if ph <= 0 {
// 		ph = pmh
// 	}
// 	if b.width <= 0 {
// 		b.setWidth(pw)
// 	}
// 	if b.height <= 0 {
// 		b.setHeight(ph)
// 	}
// }
//
// func (b *baseBar) calcParentSizes() {}
// func (b *baseBar) render() {
// 	if b.rendered || b.isHidden() {
// 		return
// 	}
// 	b.finalizeSizing()
// 	b.resetRender()
// }
//
// func (b *baseBar) getView() (view string) {
// 	return b.view
// }
//
// func (b baseBar) isHidden() (hidden bool) {
// 	return b.hidden
// }
//
// func (b *baseBar) resetRender() {
// 	if !b.getChocolate().IsRoot(b) {
// 		b.width = 0
// 		b.height = 0
// 	}
// 	b.preRendered = false
// 	b.rendered = false
// 	b.contentWidth = 0
// 	b.contentHeight = 0
// }
//
// type ChocolateBarOption func(*baseBar)
//
// func SetID(v string) ChocolateBarOption {
// 	return func(c *baseBar) {
// 		log.Println(v)
// 		c.setID(v)
// 	}
// }
//
// func SetXScaler(v Scaler) ChocolateBarOption {
// 	return func(c *baseBar) {
// 		c.setXScaler(v)
// 	}
// }
//
// func SetYScaler(v Scaler) ChocolateBarOption {
// 	return func(c *baseBar) {
// 		c.setYScaler(v)
// 	}
// }
//
// func SetHidden() ChocolateBarOption {
// 	return func(c *baseBar) {
// 		c.Hide(true)
// 	}
// }
//
// func SetSelectable() ChocolateBarOption {
// 	return func(c *baseBar) {
// 		c.setSelectable()
// 	}
// }
//
// func SetInputOnSelect() ChocolateBarOption {
// 	return func(c *baseBar) {
// 		c.setInputOnSelect()
// 	}
// }
//
// func newBaseBar(opts ...ChocolateBarOption) *baseBar {
// 	ret := &baseBar{
// 		id:            uuid.NewString(),
// 		preRendered:   false,
// 		rendered:      false,
// 		view:          "",
// 		width:         0,
// 		height:        0,
// 		contentWidth:  0,
// 		contentHeight: 0,
// 		hidden:        false,
// 		selectable:    false,
// 		inputOnSelect: false,
// 	}
// 	ret.setXScaler(NewParentScaler(1))
// 	ret.setYScaler(NewParentScaler(1))
//
// 	for _, opt := range opts {
// 		opt(ret)
// 	}
//
// 	return ret
// }
