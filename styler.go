package chocolate

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mfulz/chocolate/flavour"
)

type Styler interface {
	Get() lipgloss.Style
}

type BaseStyler struct{}

func (f *BaseStyler) Get() lipgloss.Style {
	return flavour.GetPresetNoErr(flavour.PRESET_PRIMARY_NOBORDER)
}

type RootStyler struct{}

func (r *RootStyler) Get() lipgloss.Style {
	return flavour.GetPresetNoErr(flavour.PRESET_PRIMARY)
}
