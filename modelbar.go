package chocolate

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mfulz/chocolate/flavour"
)

type (
	ModelUpdateHandlerFct           func(ChocolateBar, tea.Model) func(tea.Msg) tea.Cmd
	ModelFlavourCustomizeHandlerFct func(ChocolateBar, tea.Model, lipgloss.Style) func() lipgloss.Style
)

type BarModel struct {
	Model                   tea.Model
	UpdateHandlerFct        ModelUpdateHandlerFct
	FlavourCustomizeHandler ModelFlavourCustomizeHandlerFct
}

type modelBar struct {
	*baseBar

	// models to select from
	models map[string]*BarModel
	// running actModel
	actModel *BarModel
}

func (b *modelBar) GetStyle() lipgloss.Style {
	ret := flavour.GetPresetNoErr(flavour.PRESET_PRIMARY)

	if b.IsSelected(b) && !b.IsRoot(b) {
		ret = ret.BorderForeground(flavour.GetColorNoErr(flavour.COLOR_SECONDARY))
	}
	if b.IsFocused(b) {
		ret = flavour.GetPresetNoErr(flavour.PRESET_SECONDARY).
			BorderBackground(flavour.GetColorNoErr(flavour.COLOR_PRIMARY_BG))
	}

	if b.hasModel() && b.actModel.FlavourCustomizeHandler != nil {
		ret = b.actModel.FlavourCustomizeHandler(b, b.actModel.Model, ret)()
	}

	return ret
}

func (b *modelBar) Resize(width, height int) {
	pbar := b.GetParent(b)
	if pbar != nil {
		width, height = pbar.GetMaxSize()
	}
	b.baseBar.Resize(width, height)

	if b.models != nil {
		for _, m := range b.models {
			m.Model, _ = m.Model.Update(tea.WindowSizeMsg{Width: b.maxWidth, Height: b.maxHeight})
		}
	} else if b.actModel != nil {
		b.actModel.Model, _ = b.actModel.Model.Update(tea.WindowSizeMsg{Width: b.maxWidth, Height: b.maxHeight})
	}
}

func (b *modelBar) PreRender() bool {
	if b.baseBar.PreRender() {
		return true
	}

	pbar := b.GetParent(b)
	if pbar == nil {
		return false
	}

	b.preRendered = true

	preView := b.actModel.Model.View()
	cw, ch := lipgloss.Size(preView)
	xt, xv := b.GetScaler(XAXIS)
	yt, yv := b.GetScaler(YAXIS)

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

	b.baseBar.finalizeSizing()
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

func (b *modelBar) HandleUpdate(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ModelChangeMsg:
		b.SelectModel(msg.Model)
	}

	if b.hasModel() {
		b.actModel.Model, cmd = b.actModel.Model.Update(msg)
		cmds = append(cmds, cmd)
		if b.actModel.UpdateHandlerFct != nil {
			cmds = append(cmds, b.actModel.UpdateHandlerFct(b, b.actModel.Model)(msg))
		}
	}

	return tea.Batch(cmds...)
}

func (b modelBar) hasModel() bool {
	if b.actModel != nil {
		return b.actModel.Model != nil
	}
	return false
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

func NewModelBar(model *BarModel, opts ...baseBarOption) *modelBar {
	ret := &modelBar{
		actModel: model,
	}
	ret.baseBar = NewBaseBar(
		WithBarStyler(ret),
	)

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}

func NewMultiModelBar(act string, models map[string]*BarModel, opts ...baseBarOption) *modelBar {
	ret := &modelBar{
		models:   models,
		actModel: models[act],
	}

	ret.baseBar = NewBaseBar(
		WithBarStyler(ret),
	)

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
