package chocolate

import (
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
	width         int
	height        int
	maxWidth      int
	maxHeight     int
	contentWidth  int
	contentHeight int
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
	view string
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

func (r *rootRenderer) PreRender() {}
func (r *rootRenderer) Render() {
	r.view = r.GetStyle().
		Width(r.width).
		Height(r.height).
		Render("")
}

func (r *rootRenderer) GetView() string { return r.view }

func newRootBar() *rootRenderer {
	ret := &rootRenderer{
		BarSelector: NewDefaultSelector(),
		scaler: *newScaler(nil,
			withXparent(1, nil),
			withYparent(1, nil)),
	}

	return ret
}

type modelRenderer struct {
	parentSizes
}

type layoutRenderer struct {
	parentSizes
}

// horizontal arranged layout
type linearLayout struct {
	layoutRenderer
	totalParts int
	partSize   int
	partLast   int
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
		l.partSize = (l.maxWidth - l.contentWidth) % l.totalParts
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
