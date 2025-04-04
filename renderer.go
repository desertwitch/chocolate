package chocolate

import (
	"github.com/charmbracelet/lipgloss"
)

type chocolateBarRenderer struct {
	width   int
	height  int
	content *string
}

func (cbr *chocolateBarRenderer) setSize(width, height int) { cbr.width = width; cbr.height = height }
func (cbr *chocolateBarRenderer) render() string {
	if cbr.content != nil {
		return *cbr.content
	}
	return ""
}

type noneRenderer struct {
	chocolateBarRenderer
}

func (hbr *noneRenderer) setSize(_, _ int) {}

func newNoneRenderer() *noneRenderer {
	return &noneRenderer{
		chocolateBarRenderer: chocolateBarRenderer{},
	}
}

func newStaticRenderer(content *string) *chocolateBarRenderer {
	return &chocolateBarRenderer{
		content: content,
	}
}

type styleRenderer struct {
	chocolateBarRenderer
	style        *lipgloss.Style
	defaultStyle lipgloss.Style
	cwidth       int
	cheight      int
}

func (sr *styleRenderer) getStyle() *lipgloss.Style {
	if sr.style != nil {
		return sr.style
	}

	return &sr.defaultStyle
}

func (sr *styleRenderer) setSize(width, height int) {
	sr.chocolateBarRenderer.setSize(width, height)

	sr.cwidth = sr.width - sr.getStyle().GetHorizontalFrameSize()
	sr.cheight = sr.height - sr.getStyle().GetVerticalFrameSize()
}

func (sr *styleRenderer) render() string {
	sr.setSize(sr.width, sr.height)
	return sr.getStyle().
		Width(sr.cwidth).
		Height(sr.cheight).
		Render(sr.chocolateBarRenderer.render())
}

func newStyleRenderer(content *string, style *lipgloss.Style) *styleRenderer {
	return &styleRenderer{
		chocolateBarRenderer: *newStaticRenderer(content),
		style:                style,
		defaultStyle:         lipgloss.NewStyle(),
	}
}

type viewRenderer struct {
	bar barContainer
	styleRenderer
	viewer barViewer
}

func (vr *viewRenderer) render() string {
	var w int
	var h int
	var nw int
	var nh int

	if vr.content != nil {
		w = lipgloss.Width(*vr.content)
		h = lipgloss.Height(*vr.content)
	}
	if vr.viewer != nil {
		*vr.content = vr.viewer.View()
		nw = lipgloss.Width(*vr.content)
		nh = lipgloss.Height(*vr.content)
	}
	if w != nw || h != nh {
		if vr.bar != nil {
			vr.bar.setDirty()
		}
	}

	return vr.styleRenderer.render()
}

func newViewRenderer(viewer barViewer, style *lipgloss.Style) *viewRenderer {
	return &viewRenderer{
		styleRenderer: *newStyleRenderer(new(string), style),
		viewer:        viewer,
	}
}

type modelRenderer struct {
	viewRenderer
	model barModel
}

func (mr *modelRenderer) setSize(width, height int) {
	w := mr.cwidth
	h := mr.cheight

	mr.viewRenderer.setSize(width, height)
	if w != mr.cwidth || h != mr.cheight {
		mr.model.Resize(mr.cwidth, mr.cheight)
	}
}

func newModelRenderer(model barModel, style *lipgloss.Style) *modelRenderer {
	return &modelRenderer{
		viewRenderer: *newViewRenderer(model, style),
		model:        model,
	}
}
