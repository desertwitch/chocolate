package chocolate

import (
	tea "github.com/charmbracelet/bubbletea"
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
	case XAXIS:
		return s.x.t, s.x.v
	case YAXIS:
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
	case XAXIS:
		s.x.t = scalingType
		s.x.v = value
	case YAXIS:
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

type defaultPlacer struct {
	x struct {
		t PlacementType
		v int
	}
	y struct {
		t PlacementType
		v int
	}
}

func (s defaultPlacer) GetPlacement(axis ScalingAxis) (PlacementType, int) {
	switch axis {
	case XAXIS:
		return s.x.t, s.x.v
	case YAXIS:
		return s.y.t, s.y.v
	}

	return s.x.t, s.x.v
}

func (s *defaultPlacer) SetPlacement(axis ScalingAxis, placementType PlacementType, value int) {
	switch placementType {
	case START, CENTER, END:
		value = 0
	case POSITION:
		if value < 0 {
			value = 0
		}
	}

	switch axis {
	case XAXIS:
		s.x.t = placementType
		s.x.v = value
	case YAXIS:
		s.y.t = placementType
		s.y.v = value
	}
}

func NewDefaultPlacer() *defaultPlacer {
	return &defaultPlacer{
		x: struct {
			t PlacementType
			v int
		}{CENTER, 0},
		y: struct {
			t PlacementType
			v int
		}{CENTER, 0},
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
	focusable  bool
	overlay    bool
}

func (s defaultSelector) GetID() string          { return s.id }
func (s *defaultSelector) SetID(id string)       { s.id = id }
func (s defaultSelector) IsHidden() bool         { return s.hidden }
func (s defaultSelector) IsSelectable() bool     { return s.selectable }
func (s defaultSelector) IsFocusable() bool      { return s.focusable }
func (s *defaultSelector) isOverlay() bool       { return s.overlay }
func (s *defaultSelector) Hide(value bool)       { s.hidden = value }
func (s *defaultSelector) Selectable(value bool) { s.selectable = value }
func (s *defaultSelector) Focusable(value bool)  { s.focusable = value }
func (s *defaultSelector) setOverlay()           { s.overlay = true }

func NewDefaultSelector() *defaultSelector {
	return &defaultSelector{
		id:         uuid.NewString(),
		hidden:     false,
		selectable: false,
		focusable:  false,
		overlay:    false,
	}
}

type BaseBarStyleCustomizeHanleFct func(ChocolateBar, lipgloss.Style) func() lipgloss.Style

type baseBar struct {
	BarStyler
	BarScaler
	BarPlacer
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
	zindex        int

	styleCustomizeHandler BaseBarStyleCustomizeHanleFct
}

func (r *baseBar) SetSize(width, height int) {
	if width > 0 {
		r.width = width - r.GetStyle().GetHorizontalFrameSize()
	}
	if height > 0 {
		r.height = height - r.GetStyle().GetVerticalFrameSize()
	}
}

func (r *baseBar) finalizeSizing() {
	pbar := r.GetParent(r)
	if pbar == nil || r.IsHidden() || r.isOverlay() {
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

func (r *baseBar) resetRender() {
	if !r.IsRoot(r) && !r.isOverlay() {
		r.width = 0
		r.height = 0
	}
	r.preRendered = false
	r.rendered = false
	r.contentWidth = 0
	r.contentHeight = 0
}

func (r *baseBar) GetStyle() lipgloss.Style {
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
		ret = r.styleCustomizeHandler(r, ret)()
	}
	return ret
}

func (r *baseBar) Resize(width, height int) {
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

	if r.isOverlay() {
		if IsXParent(r) {
			xv := GetXValue(r)
			width = int(float64(width) * (float64(xv) / 100.0))
		}
		if IsYParent(r) {
			yv := GetYValue(r)
			height = int(float64(height) * (float64(yv) / 100.0))
		}
		r.width = width
		r.height = height
	}

	r.maxWidth = width
	r.maxHeight = height

	if r.IsRoot(r) {
		r.width = width
		r.height = height
	}
}

func (r *baseBar) PreRender() bool {
	if r.IsHidden() ||
		r.preRendered {
		return true
	}
	return false
}

func (r *baseBar) Render() {
	if r.rendered || r.IsHidden() {
		return
	}
	r.finalizeSizing()
	r.resetRender()
}

func (r baseBar) GetView() (view string)              { return r.view }
func (r baseBar) GetSize() (width, height int)        { return r.width, r.height }
func (r baseBar) GetMaxSize() (width, height int)     { return r.maxWidth, r.maxHeight }
func (r baseBar) GetContentSize() (width, height int) { return r.contentWidth, r.contentHeight }
func (r baseBar) GetLayout() (layout LayoutType)      { return NONE }
func (r baseBar) SetLayout(layout LayoutType)         {}
func (r baseBar) GetModel() tea.Model                 { return nil }
func (r baseBar) SelectModel(string)                  {}
func (r baseBar) GetZindex() int                      { return r.zindex }
func (r *baseBar) SetZindex(index int)                { r.zindex = index }

func (r *baseBar) setBarStyler(v BarStyler)         { r.BarStyler = v }
func (r *baseBar) setBarScaler(v BarScaler)         { r.BarScaler = v }
func (r *baseBar) setBarPlacer(v BarPlacer)         { r.BarPlacer = v }
func (r *baseBar) setBarSelector(v BarSelector)     { r.BarSelector = v }
func (r *baseBar) setBarController(v BarController) { r.BarController = v }
func (r *baseBar) setStyleCustomizeHandler(v BaseBarStyleCustomizeHanleFct) {
	r.styleCustomizeHandler = v
}

func (r *baseBar) setBarChocolate(chocolateSelector ChocolateSelector) {
	r.ChocolateSelector = chocolateSelector
}
func (r *baseBar) HandleUpdate(msg tea.Msg) tea.Cmd { return nil }

type BaseBarOption func(ChocolateBar)

func WithBarID(v string) BaseBarOption {
	return func(b ChocolateBar) {
		b.SetID(v)
	}
}

func WithBarSelectable() BaseBarOption {
	return func(b ChocolateBar) {
		b.Selectable(true)
	}
}

func WithBarFocusable() BaseBarOption {
	return func(b ChocolateBar) {
		b.Focusable(true)
	}
}

func WithBarStyler(v BarStyler) BaseBarOption {
	return func(b ChocolateBar) {
		b.setBarStyler(v)
	}
}

func WithBarScaler(v BarScaler) BaseBarOption {
	return func(b ChocolateBar) {
		b.setBarScaler(v)
	}
}

func WithBarPlacer(v BarPlacer) BaseBarOption {
	return func(b ChocolateBar) {
		b.setBarPlacer(v)
	}
}

func WithBarSelector(v BarSelector) BaseBarOption {
	return func(b ChocolateBar) {
		b.setBarSelector(v)
	}
}

func WithBarController(v BarController) BaseBarOption {
	return func(b ChocolateBar) {
		b.setBarController(v)
	}
}

func WithStyleCustomizeHandler(v BaseBarStyleCustomizeHanleFct) BaseBarOption {
	return func(b ChocolateBar) {
		b.setStyleCustomizeHandler(v)
	}
}

func WithBarXScaler(scalingType ScalingType, value int) BaseBarOption {
	return func(b ChocolateBar) {
		b.SetScaler(XAXIS, scalingType, value)
	}
}

func WithBarYScaler(scalingType ScalingType, value int) BaseBarOption {
	return func(b ChocolateBar) {
		b.SetScaler(YAXIS, scalingType, value)
	}
}

func WithBarXPlacer(placementType PlacementType, value int) BaseBarOption {
	return func(b ChocolateBar) {
		b.SetPlacement(XAXIS, placementType, value)
	}
}

func WithBarYPlacer(placementType PlacementType, value int) BaseBarOption {
	return func(b ChocolateBar) {
		b.SetPlacement(YAXIS, placementType, value)
	}
}

func NewBaseBar(opts ...BaseBarOption) *baseBar {
	scaler := NewDefaultScaler()
	placer := NewDefaultPlacer()
	controller := NewDefaultSelector()

	ret := &baseBar{
		BarScaler:     scaler,
		BarPlacer:     placer,
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
