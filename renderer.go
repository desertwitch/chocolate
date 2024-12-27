package chocolate

type Renderer interface {
	Resize(width, height int)
	PreRender()
	Render()
	GetView() string
}

type rootRenderer struct {
	width  int
	height int
}

func (r *rootRenderer) Resize(width, height int) {
	r.width = width
	r.height = height
}

type modelRenderer struct {
	xscaler   Scaler
	yscaler   Scaler
	width     int
	height    int
	maxWidth  int
	maxHeight int
}

type layoutRenderer struct {
	width     int
	height    int
	maxWidth  int
	maxHeight int

	contentWidth  int
	contentHeight int
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
