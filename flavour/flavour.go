package flavour

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type ColorName string

// default color names
// that every color definition has
// to provide
const (
	COLOR_PRIMARY      ColorName = "primary"
	COLOR_PRIMARY_BG   ColorName = "primaryBack"
	COLOR_SECONDARY    ColorName = "secondary"
	COLOR_SECONDARY_BG ColorName = "secondaryBack"
)

var defaultColors = []ColorName{
	COLOR_PRIMARY,
	COLOR_PRIMARY_BG,
	COLOR_SECONDARY,
	COLOR_SECONDARY_BG,
}

type ColorDefinitons map[ColorName]lipgloss.Color

type StylePreset string

const (
	PRESET_PRIMARY                    StylePreset = "primary"
	PRESET_PRIMARY_NOBORDER           StylePreset = "primaryNoborder"
	PRESET_PRIMARY_NOALIGN            StylePreset = "primaryNoalign"
	PRESET_PRIMARY_NOBORDER_NOALIGN   StylePreset = "primaryNoborderNoalign"
	PRESET_SECONDARY                  StylePreset = "secondary"
	PRESET_SECONDARY_NOBORDER         StylePreset = "secondaryNoborder"
	PRESET_SECONDARY_NOALIGN          StylePreset = "secondaryNoalign"
	PRESET_SECONDARY_NOBORDER_NOALIGN StylePreset = "secondaryNoborderNoalign"
)

var defaultPresets = []StylePreset{
	PRESET_PRIMARY,
	PRESET_PRIMARY_NOBORDER,
	PRESET_PRIMARY_NOALIGN,
	PRESET_PRIMARY_NOBORDER_NOALIGN,
	PRESET_SECONDARY,
	PRESET_SECONDARY_NOBORDER,
	PRESET_SECONDARY_NOALIGN,
	PRESET_SECONDARY_NOBORDER_NOALIGN,
}

type PresetDefinitions map[StylePreset]lipgloss.Style

type Flavour struct {
	colors  ColorDefinitons
	presets PresetDefinitions
	border  lipgloss.Border
	xAlign  lipgloss.Position
	yAlign  lipgloss.Position
}

func (f *Flavour) SetPreset(v StylePreset, c lipgloss.Style) error {
	for _, p := range defaultPresets {
		if v == p {
			return fmt.Errorf("default presets can't be changed '%s'", p)
		}
	}
	f.presets[v] = c
	return nil
}

func (f Flavour) GetPreset(v StylePreset) (lipgloss.Style, error) {
	if p, ok := f.presets[v]; ok {
		return p, nil
	}
	return lipgloss.NewStyle(), fmt.Errorf("preset not exiting '%s'", v)
}

func (f Flavour) GetPresetNoErr(v StylePreset) lipgloss.Style {
	if p, ok := f.presets[v]; ok {
		return p
	}
	return f.presets[PRESET_PRIMARY]
}

func (f *Flavour) SetColor(v ColorName, c lipgloss.Color) {
	f.colors[v] = c
}

func (f Flavour) GetColor(v ColorName) (lipgloss.Color, error) {
	if c, ok := f.colors[v]; ok {
		return c, nil
	}
	return lipgloss.Color(""), fmt.Errorf("color not existing '%s'", v)
}

func (f Flavour) GetFGColor(v ColorName) lipgloss.Color {
	if c, ok := f.colors[v]; ok {
		return c
	}
	return f.colors[COLOR_PRIMARY]
}

func (f Flavour) GetBGColor(v ColorName) lipgloss.Color {
	if c, ok := f.colors[v]; ok {
		return c
	}
	return f.colors[COLOR_PRIMARY_BG]
}

func GetColoredStyle(fg, bg lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(fg).
		BorderForeground(fg).
		Background(bg).
		BorderBackground(bg).
		MarginBackground(bg)
}

func (f *Flavour) initPresets() {
	pfg := f.GetFGColor(COLOR_PRIMARY)
	pbg := f.GetFGColor(COLOR_PRIMARY_BG)
	sfg := f.GetFGColor(COLOR_SECONDARY)
	sbg := f.GetFGColor(COLOR_SECONDARY_BG)

	primary := GetColoredStyle(pfg, pbg)
	secondary := GetColoredStyle(sfg, sbg)

	f.presets[PRESET_PRIMARY] = primary.
		Border(f.border).
		AlignHorizontal(f.xAlign).
		AlignVertical(f.yAlign)

	f.presets[PRESET_PRIMARY_NOBORDER] = primary.
		AlignHorizontal(f.xAlign).
		AlignVertical(f.yAlign)

	f.presets[PRESET_PRIMARY_NOALIGN] = primary.
		Border(f.border)

	f.presets[PRESET_PRIMARY_NOBORDER_NOALIGN] = primary

	f.presets[PRESET_SECONDARY] = secondary.
		Border(f.border).
		AlignHorizontal(f.xAlign).
		AlignVertical(f.yAlign)

	f.presets[PRESET_SECONDARY_NOALIGN] = secondary.
		Border(f.border)

	f.presets[PRESET_SECONDARY_NOBORDER] = secondary.
		AlignHorizontal(f.xAlign).
		AlignVertical(f.yAlign)

	f.presets[PRESET_SECONDARY_NOBORDER_NOALIGN] = secondary
}

var defaultFlavour = DefaultFlavour()

func SetPreset(v StylePreset, c lipgloss.Style) error {
	return defaultFlavour.SetPreset(v, c)
}

func GetPreset(v StylePreset) (lipgloss.Style, error) {
	return defaultFlavour.GetPreset(v)
}

func GetPresetNoErr(v StylePreset) lipgloss.Style {
	return defaultFlavour.GetPresetNoErr(v)
}

func SetColor(v ColorName, c lipgloss.Color) {
	defaultFlavour.SetColor(v, c)
}

func GetColor(v ColorName) (lipgloss.Color, error) {
	return defaultFlavour.GetColor(v)
}

func GetColorNoErr(v ColorName) lipgloss.Color {
	return defaultFlavour.GetFGColor(v)
}

type newFlavourOptions func(*Flavour) error

func WithColors(v ColorDefinitons) newFlavourOptions {
	return func(f *Flavour) error {
		for _, color := range defaultColors {
			if _, ok := v[color]; !ok {
				return fmt.Errorf("missing color definition '%s'", color)
			}
		}
		return nil
	}
}

func WithBorder(v lipgloss.Border) newFlavourOptions {
	return func(f *Flavour) error {
		f.border = v
		return nil
	}
}

func WithAlign(v lipgloss.Position) newFlavourOptions {
	return func(f *Flavour) error {
		f.xAlign = v
		f.yAlign = v
		return nil
	}
}

func WithXalign(v lipgloss.Position) newFlavourOptions {
	return func(f *Flavour) error {
		f.xAlign = v
		return nil
	}
}

func WithYalign(v lipgloss.Position) newFlavourOptions {
	return func(f *Flavour) error {
		f.yAlign = v
		return nil
	}
}

func DefaultFlavour() *Flavour {
	ret := &Flavour{
		colors:  DefaultColors,
		presets: make(PresetDefinitions),
		border:  lipgloss.RoundedBorder(),
		xAlign:  lipgloss.Center,
		yAlign:  lipgloss.Center,
	}

	ret.initPresets()
	return ret
}

func NewFlavour(opts ...newFlavourOptions) (*Flavour, error) {
	ret := &Flavour{
		colors:  DefaultColors,
		presets: make(PresetDefinitions),
		border:  lipgloss.RoundedBorder(),
		xAlign:  lipgloss.Center,
		yAlign:  lipgloss.Center,
	}

	for _, opt := range opts {
		if err := opt(ret); err != nil {
			return nil, err
		}
	}

	ret.initPresets()
	return ret, nil
}
