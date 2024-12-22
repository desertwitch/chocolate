package chocolate

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/mfulz/chocolate/flavour"
)

type defaultScaler struct {
	x struct {
		t ScalingType
		v int
	}
	y struct {
		t ScalingType
		v int
	}
}

func (s defaultScaler) GetScaler(axis ScalingAxis) (ScalingType, int) {
	switch axis {
	case X:
		return s.x.t, s.x.v
	case Y:
		return s.y.t, s.y.v
	}

	return s.x.t, s.x.v
}

func (s *defaultScaler) SetScaler(axis ScalingAxis, scalingType ScalingType, value int) {
	switch scalingType {
	case PARENT, FIXED:
		if value <= 0 {
			value = 1
		}
	case DYNAMIC:
		value = 1
	}

	switch axis {
	case X:
		s.x.t = scalingType
		s.x.v = value
	case Y:
		s.y.t = scalingType
		s.y.v = value
	}
}

func NewDefaultScaler() *defaultScaler {
	return &defaultScaler{
		x: struct {
			t ScalingType
			v int
		}{PARENT, 1},
		y: struct {
			t ScalingType
			v int
		}{PARENT, 1},
	}
}

type defaultStyler struct{}

func (s defaultStyler) GetStyle() lipgloss.Style {
	return flavour.GetPresetNoErr(flavour.PRESET_PRIMARY)
}

type defaultSelector struct {
	id         string
	hidden     bool
	selectable bool
}

func (s defaultSelector) GetID() string          { return s.id }
func (s *defaultSelector) SetID(id string)       { s.id = id }
func (s defaultSelector) IsHidden() bool         { return s.hidden }
func (s defaultSelector) IsSelectable() bool     { return s.selectable }
func (s *defaultSelector) Hide(value bool)       { s.hidden = value }
func (s *defaultSelector) Selectable(value bool) { s.selectable = value }

func NewDefaultSelector() *defaultSelector {
	return &defaultSelector{
		id:         uuid.NewString(),
		hidden:     false,
		selectable: false,
	}
}

type DefaultRendererStyleCustomizeHanleFct func(lipgloss.Style) func() lipgloss.Style

type defaultRenderer struct {
	BarStyler
	BarScaler
	BarSelector
	BarController
	ChocolateSelector

	maxWidth      int
	maxHeight     int
	width         int
	height        int
	contentWidth  int
	contentHeight int
	preRendered   bool
	rendered      bool
	view          string

	styleCustomizeHandler DefaultRendererStyleCustomizeHanleFct
}

func (r *defaultRenderer) SetSize(width, height int) {
	if width > 0 {
		r.width = width - r.GetStyle().GetHorizontalFrameSize()
	}
	if height > 0 {
		r.height = height - r.GetStyle().GetVerticalFrameSize()
	}
}

func (r *defaultRenderer) finalizeSizing() {
	pbar := r.GetParent(r)
	if pbar == nil || r.IsHidden() {
		return
	}

	pw, ph := pbar.GetSize()
	pmw, pmh := pbar.GetMaxSize()
	if pw <= 0 {
		pw = pmw
	}
	if ph <= 0 {
		ph = pmh
	}
	if r.width <= 0 {
		SetWidth(r, pw)
	}
	if r.height <= 0 {
		SetHeight(r, ph)
	}
}

func (r *defaultRenderer) resetRender() {
	if !r.IsRoot(r) {
		r.width = 0
		r.height = 0
	}
	r.preRendered = false
	r.rendered = false
	r.contentWidth = 0
	r.contentHeight = 0
}

func (r *defaultRenderer) GetStyle() lipgloss.Style {
	ret := flavour.GetPresetNoErr(flavour.PRESET_PRIMARY_NOBORDER)

	if r.BarStyler != nil {
		ret = r.BarStyler.GetStyle()
	} else {
		// root
		if r.IsRoot(r) {
			ret = flavour.GetPresetNoErr(flavour.PRESET_PRIMARY)
		}
		// selected and not root
		if r.IsSelected(r) && !r.IsRoot(r) {
			ret = ret.BorderForeground(flavour.GetColorNoErr(flavour.COLOR_SECONDARY))
		}

		// focused and not root
		if r.IsFocused(r) && !r.IsRoot(r) {
			ret = flavour.GetPresetNoErr(flavour.PRESET_SECONDARY).
				BorderBackground(flavour.GetColorNoErr(flavour.COLOR_PRIMARY_BG))
		}
	}
	if r.styleCustomizeHandler != nil {
		ret = r.styleCustomizeHandler(ret)()
	}
	return ret
}

func (r *defaultRenderer) Resize(width, height int) {
	// if there is a frame set for the bar
	// this has to be removed from the available
	// content size
	width = width - r.GetStyle().GetHorizontalFrameSize()
	height = height - r.GetStyle().GetVerticalFrameSize()

	// if this is a fixed scaling than we don't have
	// to calculate anything
	if IsXFixed(r) {
		width = GetXValue(r)
		r.width = width
	}
	if IsYFixed(r) {
		height = GetYValue(r)
		r.height = height
	}

	r.maxWidth = width
	r.maxHeight = height

	if r.IsRoot(r) {
		r.width = width
		r.height = height
	}
}

func (r *defaultRenderer) PreRender() bool {
	if r.IsHidden() ||
		r.preRendered {
		return true
	}
	return false
}

func (r *defaultRenderer) Render() {
	if r.rendered || r.IsHidden() {
		return
	}
	r.finalizeSizing()
	r.resetRender()
}

func (r defaultRenderer) GetView() (view string)              { return r.view }
func (r defaultRenderer) GetSize() (width, height int)        { return r.width, r.height }
func (r defaultRenderer) GetMaxSize() (width, height int)     { return r.maxWidth, r.maxHeight }
func (r defaultRenderer) GetContentSize() (width, height int) { return r.contentWidth, r.contentHeight }

func (r *defaultRenderer) SetBarChocolate(chocolateSelector ChocolateSelector) {
	r.ChocolateSelector = chocolateSelector
}

type defaultRendererOption func(*defaultRenderer)

func WithBarStyler(v BarStyler) defaultRendererOption {
	return func(r *defaultRenderer) {
		r.BarStyler = v
	}
}

func WithBarScaler(v BarScaler) defaultRendererOption {
	return func(r *defaultRenderer) {
		r.BarScaler = v
	}
}

func WithBarSelector(v BarSelector) defaultRendererOption {
	return func(r *defaultRenderer) {
		r.BarSelector = v
	}
}

func WithBarController(v BarController) defaultRendererOption {
	return func(r *defaultRenderer) {
		r.BarController = v
	}
}

func WithStyleCustomizeHandler(v DefaultRendererStyleCustomizeHanleFct) defaultRendererOption {
	return func(r *defaultRenderer) {
		r.styleCustomizeHandler = v
	}
}

func NewDefaultRenderer(opts ...defaultRendererOption) *defaultRenderer {
	scaler := NewDefaultScaler()
	controller := NewDefaultSelector()

	ret := &defaultRenderer{
		BarScaler:     scaler,
		BarSelector:   controller,
		BarController: controller,
		maxWidth:      0,
		maxHeight:     0,
		width:         0,
		height:        0,
		contentWidth:  0,
		contentHeight: 0,
		preRendered:   false,
		rendered:      false,
		view:          "",
	}

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
