package chocolate

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mfulz/chocolate/flavour"
)

type (
	ModelUpdateHandlerFct           func(NChocolateBar, tea.Model) func(tea.Msg) tea.Cmd
	ModelFlavourCustomizeHandlerFct func(NChocolateBar, tea.Model, lipgloss.Style) func() lipgloss.Style
	BarFlavourCustomizeHandlerFct   func(NChocolateBar, lipgloss.Style) func() lipgloss.Style
)

type ModelBarModel struct {
	Model                   tea.Model
	UpdateHandlerFct        ModelUpdateHandlerFct
	FlavourCustomizeHandler ModelFlavourCustomizeHandlerFct
}

type modelBar struct {
	*defaultRenderer

	// models to select from
	models map[string]*ModelBarModel
	// running actModel
	actModel *ModelBarModel
}

func (b *modelBar) GetStyle() lipgloss.Style {
	ret := flavour.GetPresetNoErr(flavour.PRESET_PRIMARY)

	if b.IsSelected(b) && !b.IsRoot(b) {
		ret = ret.BorderForeground(flavour.GetColorNoErr(flavour.COLOR_SECONDARY))
	}

	return ret
}

func (b *modelBar) Resize(width, height int) {
	pbar := b.GetParent(b)
	if pbar != nil {
		width, height = pbar.GetMaxSize()
	}
	b.defaultRenderer.Resize(width, height)

	if b.models != nil {
		for _, m := range b.models {
			m.Model, _ = m.Model.Update(tea.WindowSizeMsg{Width: b.maxWidth, Height: b.maxHeight})
		}
	} else if b.actModel != nil {
		b.actModel.Model, _ = b.actModel.Model.Update(tea.WindowSizeMsg{Width: b.maxWidth, Height: b.maxHeight})
	}
}

func (b *modelBar) PreRender() bool {
	if b.defaultRenderer.PreRender() {
		return true
	}

	pbar := b.GetParent(b)
	if pbar == nil {
		return false
	}

	b.preRendered = true

	preView := b.actModel.Model.View()
	cw, ch := lipgloss.Size(preView)
	xt, xv := b.GetScaler(X)
	yt, yv := b.GetScaler(Y)

	switch xt {
	case FIXED:
		b.contentWidth = xv + b.GetStyle().GetHorizontalFrameSize()
		b.width = xv
	case DYNAMIC:
		b.contentWidth = cw + b.GetStyle().GetHorizontalFrameSize()
		b.width = cw
	}

	switch yt {
	case FIXED:
		b.contentHeight = yv + b.GetStyle().GetVerticalFrameSize()
		b.height = yv
	case DYNAMIC:
		b.contentHeight = ch + b.GetStyle().GetVerticalFrameSize()
		b.height = ch
	}

	return true
}

func (b *modelBar) finalizeSizing() {
	if b.IsHidden() {
		return
	}

	b.defaultRenderer.finalizeSizing()
	b.actModel.Model, _ = b.actModel.Model.Update(tea.WindowSizeMsg{Width: b.width, Height: b.height})
}

func (b *modelBar) Render() {
	if b.rendered || b.IsHidden() {
		return
	}
	b.finalizeSizing()

	b.view = b.GetStyle().
		Width(b.width).
		Height(b.height).
		Render(b.actModel.Model.View())
	b.rendered = true

	b.resetRender()
}

func (b modelBar) GetModel() tea.Model {
	return b.actModel.Model
}

func (b *modelBar) SelectModel(v string) {
	if b.models == nil {
		return
	}
	if m, ok := b.models[v]; ok {
		b.actModel = m
	}
}

type ModelBarOption func(*modelBar)

func ModelBarID(id string) ModelBarOption {
	return func(m *modelBar) {
		m.SetID(id)
	}
}

func ModelBarXScaler(scalingType ScalingType, value int) ModelBarOption {
	return func(m *modelBar) {
		m.SetScaler(X, scalingType, value)
	}
}

func ModelBarYScaler(scalingType ScalingType, value int) ModelBarOption {
	return func(m *modelBar) {
		m.SetScaler(Y, scalingType, value)
	}
}

func ModelBarSelectable() ModelBarOption {
	return func(m *modelBar) {
		m.Selectable(true)
	}
}

func NewModelBar(model *ModelBarModel, opts ...ModelBarOption) *modelBar {
	ret := &modelBar{
		actModel: model,
	}
	ret.defaultRenderer = NewDefaultRenderer(
		WithBarStyler(ret),
	)

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}

func NewMultiModelBar(act string, models map[string]*ModelBarModel, opts ...ModelBarOption) *modelBar {
	ret := &modelBar{
		models:   models,
		actModel: models[act],
	}

	ret.defaultRenderer = NewDefaultRenderer(
		WithBarStyler(ret),
	)

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
