package chocolate

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mfulz/chocolate/flavour"
)

type Bar interface {
	BarSelector
	Renderer
}

type Renderer interface {
	Resize(width, height int)
	PreRender()
	Render()
	GetView() string
}

type childSizes struct {
	width     int
	height    int
	maxWidth  int
	maxHeight int
}

type parentSizes struct {
	childSizes

	contentWidth  int
	contentHeight int
}

type rootRenderer struct {
	BarSelector
	scaler
	parentSizes
	layouter
	view string
}

type parentBar interface {
	AddView(string)
}

func (r *rootRenderer) Resize(width, height int) {
	r.width = width - r.GetStyle().GetHorizontalFrameSize()
	r.height = height - r.GetStyle().GetVerticalFrameSize()

	r.maxWidth = r.width
	r.maxHeight = r.height
}

func (r *rootRenderer) GetStyle() lipgloss.Style {
	return flavour.GetPresetNoErr(flavour.PRESET_PRIMARY)
}

func (r *rootRenderer) PreRender() { r.calcPartSize() }
func (r *rootRenderer) Render() {
	r.view = r.GetStyle().
		Width(r.width).
		Height(r.height).
		Render(r.render(r.GetStyle().Height(r.height).Width(r.width)))
}

func (r *rootRenderer) GetView() string { return r.view }

func newRootBar() *rootRenderer {
	ret := &rootRenderer{
		BarSelector: NewDefaultSelector(),
		scaler: *newScaler(nil,
			withXparent(1, nil),
			withYparent(1, nil)),
	}
	l := &linearLayout{parentSizes: &ret.parentSizes}
	ret.layouter = l

	return ret
}

func (r *linearLayout) GetMaxX() int { return r.maxWidth }
func (r *linearLayout) GetMaxY() int { return r.maxHeight }

type modelRenderer struct {
	BarSelector
	scaler
	childSizes
	view string

	// models to select from
	models map[string]*BarModel
	// running ActModel
	ActModel *BarModel
}

func (b modelRenderer) hasModel() bool {
	if b.ActModel != nil {
		return b.ActModel.Model != nil
	}
	return false
}

func (r *modelRenderer) GetStyle() lipgloss.Style {
	ret := flavour.GetPresetNoErr(flavour.PRESET_PRIMARY)

	// if b.IsSelected(b) && !b.IsRoot(b) {
	// 	ret = ret.BorderForeground(flavour.GetColorNoErr(flavour.COLOR_SECONDARY))
	// }
	// if b.IsFocused(b) {
	// 	ret = flavour.GetPresetNoErr(flavour.PRESET_SECONDARY).
	// 		BorderBackground(flavour.GetColorNoErr(flavour.COLOR_PRIMARY_BG))
	// }

	// if r.hasModel() && r.actModel.FlavourCustomizeHandler != nil {
	// 	ret = r.actModel.FlavourCustomizeHandler(r, r.actModel.Model, ret)()
	// }

	return ret
}

func (r *modelRenderer) Resize(width, height int) {
	w := r.scaler.GetMaxX() - r.GetStyle().GetHorizontalFrameSize()
	h := r.scaler.GetMaxY() - r.GetStyle().GetVerticalFrameSize()

	if w <= 0 {
		w = width - r.GetStyle().GetHorizontalFrameSize()
	}
	if h <= 0 {
		h = width - r.GetStyle().GetVerticalFrameSize()
	}
	r.maxWidth = w
	r.maxHeight = h

	if r.models != nil {
		for _, m := range r.models {
			m.Model, _ = m.Model.Update(tea.WindowSizeMsg{Width: r.maxWidth, Height: r.maxHeight})
		}
	} else if r.ActModel != nil {
		r.ActModel.Model, _ = r.ActModel.Model.Update(tea.WindowSizeMsg{Width: r.maxWidth, Height: r.maxHeight})
	}
}

func (r *modelRenderer) GetMaxX() int { return r.maxWidth }
func (r *modelRenderer) GetMaxY() int { return r.maxHeight }

func (r *modelRenderer) PreRender() {
	preView := r.ActModel.Model.View()
	cw, ch := lipgloss.Size(preView)
	r.SetContentSize(cw, ch)
}

func (r *modelRenderer) Render() {
	if r.IsHidden() {
		return
	}
	r.width, r.height = r.FinalSize(r.width, r.height)
	r.width -= r.GetStyle().GetHorizontalFrameSize()
	r.height -= r.GetStyle().GetVerticalFrameSize()

	r.ActModel.Model, _ = r.ActModel.Model.Update(tea.WindowSizeMsg{Width: r.width, Height: r.height})

	r.view = r.GetStyle().
		Width(r.width).
		Height(r.height).
		Render(r.ActModel.Model.View())
	log.Printf("w, h: %d, %d\n", r.width, r.height)
	r.p.AddView(r.view)
}

func (r *modelRenderer) GetView() string { return r.view }

type layoutRenderer struct {
	parentSizes
}

type layouter interface {
	ContentSizer
	ParentSizer
	calcPartSize()
	render(lipgloss.Style) string
	parentBar
}

// horizontal arranged layout
type linearLayout struct {
	*parentSizes
	totalParts int
	partSize   int
	partLast   int
	childViews []string
}

func (r *linearLayout) AddView(v string) {
	r.childViews = append(r.childViews, v)
}

func (l *linearLayout) render(s lipgloss.Style) string {
	var views []string
	sb := s.
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		Height(l.height)

	for _, b := range l.childViews {
		views = append(views, sb.Render(b))
	}
	ret := sb.Render(lipgloss.JoinHorizontal(0, views...))
	l.childViews = []string{}
	return ret
}

func (l *linearLayout) AddPartsX(parts int) { l.totalParts += parts }
func (l *linearLayout) AddPartsY(parts int) {}

func (r *linearLayout) AddContentX(width int) { r.contentWidth += width }
func (r *linearLayout) AddContentY(height int) {
	if height > r.contentHeight {
		r.contentHeight = height
	}
}

func (l *linearLayout) calcPartSize() {
	if l.totalParts > 0 {
		l.partSize = (l.maxWidth - l.contentWidth) / l.totalParts
		l.partLast = (l.maxWidth - l.contentWidth) % l.totalParts
	}
}

func (l *linearLayout) TakePartsY(parts int) int {
	if l.height > 0 {
		return l.height
	}
	return l.maxHeight
}

func (l *linearLayout) TakePartsX(parts int) int {
	if l.totalParts <= 0 {
		return 0
	}

	l.totalParts -= parts
	width := parts * l.partSize
	if l.totalParts <= 0 {
		width += l.partLast
	}

	return width
}

// vertical arranged layout
type listLayout struct {
	layoutRenderer
	totalParts int
	partSize   int
	partLast   int
}

func (l *listLayout) AddPartsX(parts int) {}
func (l *listLayout) AddPartsY(parts int) { l.totalParts += parts }

func (r *listLayout) AddContentY(height int) { r.contentHeight += height }
func (r *listLayout) AddContentX(width int) {
	if width > r.contentWidth {
		r.contentWidth = width
	}
}

func (l *listLayout) calcPartSize() {
	if l.totalParts > 0 {
		l.partSize = (l.maxHeight - l.contentHeight) / l.totalParts
		l.partSize = (l.maxHeight - l.contentHeight) % l.totalParts
	}
}

func (l *listLayout) TakePartsX(parts int) int {
	if l.width > 0 {
		return l.width
	}
	return l.maxWidth
}

func (l *listLayout) TakePartsY(parts int) int {
	if l.totalParts <= 0 {
		return 0
	}

	l.totalParts -= parts
	height := parts * l.partSize
	if l.totalParts <= 0 {
		height += l.partLast
	}

	return height
}
