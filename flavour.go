package chocolate

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type FlavourType int

const (
	FLAVOUR_PRIMARY FlavourType = iota
	FLAVOUR_SECONDARY
)

type ColorType int

const (
	FOREGROUND_PRIMARY ColorType = iota
	FOREGROUND_SECONDARY
	BACKGROUND_PRIMARY
	BACKGROUND_SECONDARY
	FOREGROUND_HIGHLIGHT_PRIMARY
	FOREGROUND_HIGHLIGHT_SECONDARY
	BACKGROUND_HIGHLIGHT_PRIMARY
	BACKGROUND_HIGHLIGHT_SECONDARY
)

type BorderType int

const (
	ROUND BorderType = iota
	BLOCK
	DOUBLE
	HIDDEN
)

type Flavour interface {
	GetColor(ColorType) lipgloss.Color
	GetBorder() lipgloss.Border

	SetColor(ColorType, uint8)
	SetBorder(BorderType)

	GetFrameSize() (int, int)
}

type flavour struct {
	colors [8]uint8 // ansi color codes
	border BorderType
}

func (f flavour) GetColor(v ColorType) lipgloss.Color {
	return lipgloss.Color(fmt.Sprint(f.colors[v]))
}

func (f flavour) GetBorder() lipgloss.Border {
	switch f.border {
	case ROUND:
		return lipgloss.RoundedBorder()
	case BLOCK:
		return lipgloss.BlockBorder()
	case DOUBLE:
		return lipgloss.DoubleBorder()
	case HIDDEN:
		return lipgloss.HiddenBorder()
	default:
		return lipgloss.Border{}
	}
}

func (f *flavour) SetColor(c ColorType, v uint8) {
	f.colors[c] = v
}

func (f *flavour) SetBorder(v BorderType) {
	f.border = v
}

func (f *flavour) GetFrameSize() (int, int) {
	s := lipgloss.NewStyle().
		Border(f.GetBorder())

	return s.GetFrameSize()
}

type flavourOptions func(*flavour)

func Border(v BorderType) flavourOptions {
	return func(f *flavour) {
		f.border = v
	}
}

func Color(c ColorType, v uint8) flavourOptions {
	return func(f *flavour) {
		f.colors[c] = v
	}
}

func ColorPreset(v [8]uint8) flavourOptions {
	return func(f *flavour) {
		f.colors = v
	}
}

func NewFlavour(opts ...flavourOptions) Flavour {
	ret := &flavour{
		border: ROUND,
		colors: WhiteBlack,
	}

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
