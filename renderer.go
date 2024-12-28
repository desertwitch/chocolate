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
	getScaler() *scaler
	parentBar
	setParentBar(parentBar)
}

type Renderer interface {
	Resize(width, height int)
	PreRender()
	Render()
	GetView() string
}

type rootRenderer struct {
	BarSelector
	*scaler
	joinFct    func(lipgloss.Position, ...string) string
	childViews []string
	view       string
}

func (r *rootRenderer) getScaler() *scaler { return r.scaler }

type parentBar interface {
	AddView(string)
}

func (r *rootRenderer) setParentBar(parent parentBar) {}
func (r *rootRenderer) Resize(width, height int) {
	r.setWidth(width - r.GetStyle().GetHorizontalFrameSize())
	r.setHeight(height - r.GetStyle().GetVerticalFrameSize())

	r.setMaxWidth(r.getWidth())
	r.setMaxHeight(r.getHeight())
}

func (r *rootRenderer) GetStyle() lipgloss.Style {
	return flavour.GetPresetNoErr(flavour.PRESET_PRIMARY)
}

func (r *rootRenderer) PreRender() {}
func (r *rootRenderer) Render() {
	var views []string
	s := r.GetStyle().
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		// Width(r.getWidth())
		Height(r.getHeight())

	for _, b := range r.childViews {
		views = append(views, s.Render(b))
	}

	r.view = r.GetStyle().
		Width(r.getWidth()).
		Height(r.getHeight()).
		Render(s.Render(r.joinFct(0, views...)))
	r.childViews = []string{}
}

func (r *rootRenderer) GetView() string  { return r.view }
func (r *rootRenderer) AddView(v string) { r.childViews = append(r.childViews, v) }

func newRootBar(list bool) *rootRenderer {
	ret := &rootRenderer{
		BarSelector: NewDefaultSelector(),
	}

	if list {
		ret.scaler = newListScaler(
			&parentCreator{1},
			&parentCreator{1},
		)
		ret.joinFct = lipgloss.JoinVertical
	} else {
		ret.scaler = newLinearScaler(
			&parentCreator{1},
			&parentCreator{1},
		)
		ret.joinFct = lipgloss.JoinHorizontal
	}

	return ret
}

type modelRenderer struct {
	BarSelector
	*scaler
	parent parentBar
	view   string

	// models to select from
	models map[string]*BarModel
	// running ActModel
	ActModel *BarModel
}

func (r *modelRenderer) getScaler() *scaler            { return r.scaler }
func (r *modelRenderer) setParentBar(parent parentBar) { r.parent = parent }
func (b *modelRenderer) AddView(view string)           {}

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
	w := r.getParentMaxWidth() - r.GetStyle().GetHorizontalFrameSize()
	h := r.getParentMaxHeight() - r.GetStyle().GetVerticalFrameSize()

	if w <= 0 {
		w = width - r.GetStyle().GetHorizontalFrameSize()
	}
	if h <= 0 {
		h = width - r.GetStyle().GetVerticalFrameSize()
	}
	r.setMaxWidth(w)
	r.setMaxHeight(h)

	if r.models != nil {
		for _, m := range r.models {
			m.Model, _ = m.Model.Update(tea.WindowSizeMsg{Width: w, Height: h})
		}
	} else if r.ActModel != nil {
		r.ActModel.Model, _ = r.ActModel.Model.Update(tea.WindowSizeMsg{Width: w, Height: h})
	}
}

func (r *modelRenderer) PreRender() {
	preView := r.ActModel.Model.View()
	cw, ch := lipgloss.Size(preView)
	cw += r.GetStyle().GetHorizontalFrameSize()
	ch += r.GetStyle().GetVerticalFrameSize()
	r.setContentSize(cw, ch)
}

func (r *modelRenderer) Render() {
	if r.IsHidden() {
		return
	}
	w, h := r.finalizeSize(r.GetStyle().GetFrameSize())
	r.ActModel.Model, _ = r.ActModel.Model.Update(tea.WindowSizeMsg{Width: w, Height: h})

	r.view = r.GetStyle().
		Width(w).
		Height(h).
		Render(r.ActModel.Model.View())
	log.Printf("w, h: %d, %d\n", w, h)
	r.parent.AddView(r.view)
}

func (r *modelRenderer) GetView() string { return r.view }
