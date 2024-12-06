package chocolate

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
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
	NONE BorderType = iota
	ROUND
	BLOCK
	DOUBLE
	HIDDEN
)

type Alignment int

const (
	NO Alignment = iota
	START
	END
	CENTER
)

type FlavourPrefs struct {
	foreground          ColorType
	background          ColorType
	foregroundBorder    ColorType
	backgroundBorder    ColorType
	borderType          BorderType
	horizontalAlignment Alignment
	verticalAlignment   Alignment
}

func (p FlavourPrefs) Foreground(v ColorType) FlavourPrefs {
	p.foreground = v
	return p
}

func (p FlavourPrefs) Background(v ColorType) FlavourPrefs {
	p.background = v
	return p
}

func (p FlavourPrefs) ForegroundBorder(v ColorType) FlavourPrefs {
	p.foregroundBorder = v
	return p
}

func (p FlavourPrefs) BackgroundBorder(v ColorType) FlavourPrefs {
	p.backgroundBorder = v
	return p
}

func (p FlavourPrefs) BorderType(v BorderType) FlavourPrefs {
	p.borderType = v
	return p
}

func (p FlavourPrefs) HorizontalAlignment(v Alignment) FlavourPrefs {
	p.horizontalAlignment = v
	return p
}

func (p FlavourPrefs) VerticalAlignment(v Alignment) FlavourPrefs {
	p.verticalAlignment = v
	return p
}

func NewFlavourPrefs() FlavourPrefs {
	ret := FlavourPrefs{
		foreground:          FOREGROUND_PRIMARY,
		background:          BACKGROUND_PRIMARY,
		foregroundBorder:    FOREGROUND_PRIMARY,
		backgroundBorder:    BACKGROUND_PRIMARY,
		borderType:          NONE,
		horizontalAlignment: CENTER,
		verticalAlignment:   CENTER,
	}

	return ret
}

type Flavour interface {
	GetColor(ColorType) lipgloss.Color
	GetBorder() lipgloss.Border
	GetBorderType() BorderType

	SetColor(ColorType, uint8)
	SetBorder(BorderType)

	GetStyle(FlavourPrefs) lipgloss.Style
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

func (f flavour) GetBorderType() BorderType {
	return f.border
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

func (f *flavour) GetHorizontalFrameSize() int {
	s := lipgloss.NewStyle().
		Border(f.GetBorder())

	return s.GetHorizontalFrameSize()
}

func (f *flavour) GetVerticalFrameSize() int {
	s := lipgloss.NewStyle().
		Border(f.GetBorder())

	return s.GetVerticalFrameSize()
}

func (f *flavour) GetStyle(v FlavourPrefs) lipgloss.Style {
	s := lipgloss.NewStyle()

	if v.borderType != NONE {
		s = s.Border(f.GetBorder())
	}

	switch v.horizontalAlignment {
	case START:
		s = s.AlignHorizontal(lipgloss.Top)
	case END:
		s = s.AlignHorizontal(lipgloss.Bottom)
	case CENTER:
		s = s.AlignHorizontal(lipgloss.Center)
	}

	switch v.verticalAlignment {
	case START:
		s = s.AlignVertical(lipgloss.Left)
	case END:
		s = s.AlignVertical(lipgloss.Right)
	case CENTER:
		s = s.AlignVertical(lipgloss.Center)
	}

	s = s.Foreground(f.GetColor(v.foreground))
	s = s.Background(f.GetColor(v.background))
	s = s.BorderForeground(f.GetColor(v.foregroundBorder))
	s = s.BorderBackground(f.GetColor(v.backgroundBorder))

	return s
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
