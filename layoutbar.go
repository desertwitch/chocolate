package chocolate

import (
	"log"

	"github.com/charmbracelet/lipgloss"
)

type layoutBar struct {
	*defaultRenderer

	layout     LayoutType
	totalParts int
	partSize   int
	partLast   int
}

func (b *layoutBar) Resize(width, height int) {
	pbar := b.GetParent(b)
	if pbar != nil {
		width, height = pbar.GetMaxSize()
	}
	b.defaultRenderer.Resize(width, height)
}

func (b *layoutBar) PreRender() bool {
	if b.defaultRenderer.PreRender() {
		return true
	}

	b.preRendered = true
	for _, child := range b.GetChildren(b) {
		if child.IsHidden() {
			continue
		}
		cw, ch := child.GetContentSize()
		xt, xv := child.GetScaler(X)
		yt, yv := child.GetScaler(Y)

		switch xt {
		case DYNAMIC, FIXED:
			b.contentWidth += cw
		case PARENT:
			if b.layout == LINEAR {
				b.totalParts += xv
			}
		}
		switch yt {
		case DYNAMIC, FIXED:
			b.contentHeight += ch
		case PARENT:
			if b.layout == LIST {
				b.totalParts += yv
			}
		}
	}
	b.calcParentSizes()

	return true
}

func (b *layoutBar) calcParentsHorizontal() {
	if b.totalParts > 0 {
		partSize := (b.maxWidth - b.contentWidth) / b.totalParts
		partLast := (b.maxWidth - b.contentWidth) % b.totalParts

		for _, child := range b.GetChildren(b) {
			if !IsXParent(child) || child.IsHidden() {
				continue
			}

			childParts := GetXValue(child)
			width := childParts * partSize
			b.totalParts -= childParts
			if b.totalParts <= 0 {
				width += partLast
			}
			log.Printf("width: %d\n", width)
			SetWidth(child, width)
			b.contentWidth += width
		}
	}

	if !b.IsRoot(b) {
		b.width = b.contentWidth
	}
	b.contentWidth += b.GetStyle().GetHorizontalFrameSize()
}

func (b *layoutBar) calcParentsVertical() {
	if b.totalParts > 0 {
		partSize := (b.maxHeight - b.contentHeight) / b.totalParts
		partLast := (b.maxHeight - b.contentHeight) % b.totalParts

		for _, child := range b.GetChildren(b) {
			log.Printf("WTF VERT\n")
			if !IsYParent(child) || child.IsHidden() {
				continue
			}

			childParts := GetYValue(child)
			height := childParts * partSize
			b.totalParts -= childParts
			if b.totalParts <= 0 {
				height += partLast
			}
			log.Printf("height: %d\n", height)
			SetHeight(child, height)
			b.contentHeight += height
		}
	}

	if !b.IsRoot(b) {
		b.height = b.contentHeight
	}
	b.contentHeight += b.GetStyle().GetVerticalFrameSize()
}

func (b *layoutBar) calcParentSizes() {
	if b.IsHidden() {
		return
	}

	switch b.layout {
	case LINEAR:
		b.calcParentsHorizontal()
	case LIST:
		b.calcParentsVertical()
	}
}

func (b *layoutBar) Render() {
	var bars []string

	if b.rendered || b.IsHidden() {
		return
	}
	b.finalizeSizing()

	children := b.GetChildren(b)
	switch b.layout {
	case LINEAR:
		for _, c := range children {
			if c.IsHidden() {
				continue
			}
			s := b.GetStyle().
				BorderTop(false).
				BorderBottom(false).
				BorderLeft(false).
				BorderRight(false).
				Height(b.height)
			bars = append(bars, s.Render(c.GetView()))
		}
		s := b.GetStyle()
		if b.IsRoot(b) {
			s = s.Width(b.width)
		}
		b.view = s.Render(lipgloss.JoinHorizontal(0, bars...))
	case LIST:
		for _, c := range children {
			if c.IsHidden() {
				continue
			}
			s := b.GetStyle().
				BorderTop(false).
				BorderBottom(false).
				BorderLeft(false).
				BorderRight(false).
				Width(b.width)
			bars = append(bars, s.Render(c.GetView()))
		}
		s := b.GetStyle()
		if b.IsRoot(b) {
			s = s.Height(b.height)
		}
		b.view = s.Render(lipgloss.JoinVertical(0, bars...))
	}
	b.rendered = true

	if b.IsRoot(b) {
		w, h := lipgloss.Size(b.view)
		w -= b.GetStyle().GetHorizontalFrameSize()
		h -= b.GetStyle().GetVerticalFrameSize()
		if w > b.width || h > b.height {
			b.view = "Window too small"
		}
	}

	b.resetRender()
}

func (b *layoutBar) resetRender() {
	b.defaultRenderer.resetRender()
	b.totalParts = 0
}

func (b *layoutBar) SetLayout(layout LayoutType) {
	b.layout = layout
}

func (b layoutBar) GetLayout() LayoutType {
	return b.layout
}

type LLayoutBarOption func(*layoutBar)

func LayoutBarID(id string) LLayoutBarOption {
	return func(m *layoutBar) {
		m.SetID(id)
	}
}

func LayoutBarXScaler(scalingType ScalingType, value int) LLayoutBarOption {
	return func(m *layoutBar) {
		m.SetScaler(X, scalingType, value)
	}
}

func LayoutBarYScaler(scalingType ScalingType, value int) LLayoutBarOption {
	return func(m *layoutBar) {
		m.SetScaler(Y, scalingType, value)
	}
}

func NewLayoutBar(layout LayoutType, opts ...LLayoutBarOption) *layoutBar {
	ret := &layoutBar{
		layout: layout,
	}
	ret.defaultRenderer = NewDefaultRenderer()

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
