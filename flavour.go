package chocolate

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type FlavourStyleSelector string

const (
	TS_DEFAULT  FlavourStyleSelector = "default"
	TS_SELECTED FlavourStyleSelector = "selected"
	TS_FOCUSED  FlavourStyleSelector = "focused"
)

var defaultSelectors []FlavourStyleSelector = []FlavourStyleSelector{
	TS_DEFAULT,
	TS_SELECTED,
	TS_FOCUSED,
}

type chocolateFlavour struct {
	styles map[FlavourStyleSelector]*lipgloss.Style
}

func (t *chocolateFlavour) getStyles(styles ...FlavourStyleSelector) (map[FlavourStyleSelector]*lipgloss.Style, *lipgloss.Style, FlavourStyleSelector) {
	current := new(lipgloss.Style)
	selected := TS_DEFAULT
	var selectedStyles map[FlavourStyleSelector]*lipgloss.Style

	if len(styles) <= 0 {
		selectedStyles = t.styles
		*current = *t.styles[TS_DEFAULT]
	} else {
		selectedStyles = make(map[FlavourStyleSelector]*lipgloss.Style)
		first := true
		for _, style := range styles {
			if s, ok := t.styles[style]; ok {
				if first {
					*current = *s
					selected = style
					first = false
				}
				selectedStyles[style] = s
			}
		}
	}

	return selectedStyles, current, selected
}

func (t *chocolateFlavour) setStyle(selector FlavourStyleSelector, style *lipgloss.Style) {
	if t.styles == nil {
		t.styles = make(map[FlavourStyleSelector]*lipgloss.Style)
	}

	t.styles[selector] = style
}

func (t *chocolateFlavour) SetDefault(style *lipgloss.Style) {
	t.setStyle(TS_DEFAULT, style)

	// for k, v := range t.styles {
	// 	if strings.EqualFold(string(k), string(TS_DEFAULT)) {
	// 		continue
	// 	}
	//
	// 	// *t.styles[k] = v.Inherit(*style)
	// }
}

func (t *chocolateFlavour) SetStyle(selector FlavourStyleSelector, style *lipgloss.Style) {
	if strings.EqualFold(string(TS_DEFAULT), string(selector)) {
		t.SetDefault(style)
		return
	}

	t.setStyle(selector, style)
}

type (
	ChocolateThemeOption func(*chocolateFlavour)
	ThemeStyleModifier   func(lipgloss.Style) lipgloss.Style
)

func Align(p ...lipgloss.Position) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.Align(p...)
	}
}

func AlignHorizontal(p lipgloss.Position) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.AlignHorizontal(p)
	}
}

func AlignVertical(p lipgloss.Position) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.AlignVertical(p)
	}
}

func Background(c lipgloss.TerminalColor) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.Background(c)
	}
}

func Blink(v bool) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.Blink(v)
	}
}

func Border(b lipgloss.Border, sides ...bool) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.Border(b, sides...)
	}
}

func BorderBackground(c ...lipgloss.TerminalColor) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderBackground(c...)
	}
}

func BorderForeground(c ...lipgloss.TerminalColor) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderForeground(c...)
	}
}

func BorderBottom(v bool) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderBottom(v)
	}
}

func BorderBottomBackground(c lipgloss.TerminalColor) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderBottomBackground(c)
	}
}

func BorderBottomForeground(c lipgloss.TerminalColor) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderBottomForeground(c)
	}
}

func BorderLeft(v bool) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderLeft(v)
	}
}

func BorderLeftBackground(c lipgloss.TerminalColor) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderLeftBackground(c)
	}
}

func BorderLeftForeground(c lipgloss.TerminalColor) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderLeftForeground(c)
	}
}

func BorderRight(v bool) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderRight(v)
	}
}

func BorderRightBackground(c lipgloss.TerminalColor) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderRightBackground(c)
	}
}

func BorderRightForeground(c lipgloss.TerminalColor) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderRightForeground(c)
	}
}

func BorderTop(v bool) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderTop(v)
	}
}

func BorderTopBackground(c lipgloss.TerminalColor) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderTopBackground(c)
	}
}

func BorderTopForeground(c lipgloss.TerminalColor) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderTopForeground(c)
	}
}

func BorderStyle(b lipgloss.Border) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.BorderStyle(b)
	}
}

func Faint(v bool) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.Faint(v)
	}
}

func Foreground(c lipgloss.TerminalColor) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.Foreground(c)
	}
}

func Italic(v bool) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.Italic(v)
	}
}

func Margin(i ...int) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.Margin(i...)
	}
}

func MarginBackground(c lipgloss.TerminalColor) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.MarginBackground(c)
	}
}

func MarginBottom(i int) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.MarginBottom(i)
	}
}

func MarginLeft(i int) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.MarginLeft(i)
	}
}

func MarginRight(i int) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.MarginRight(i)
	}
}

func MarginTop(i int) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.MarginTop(i)
	}
}

func Padding(i ...int) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.Padding(i...)
	}
}

func PaddingBottom(i int) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.PaddingBottom(i)
	}
}

func PaddingLeft(i int) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.PaddingLeft(i)
	}
}

func PaddingRight(i int) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.PaddingRight(i)
	}
}

func PaddingTor(i int) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.PaddingTop(i)
	}
}

func Reverse(v bool) ThemeStyleModifier {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.Reverse(v)
	}
}

func DefaultStyle(style *lipgloss.Style) ChocolateThemeOption {
	return func(t *chocolateFlavour) {
		t.SetDefault(style)
	}
}

func WithDefaults(s ...*lipgloss.Style) ChocolateThemeOption {
	return func(t *chocolateFlavour) {
		for i, style := range s {
			t.SetStyle(defaultSelectors[i], style)
			if i >= len(defaultSelectors)-1 {
				return
			}
		}
	}
}

func NewChocolateFlavour(opts ...ChocolateThemeOption) *chocolateFlavour {
	s := lipgloss.NewStyle()
	ret := &chocolateFlavour{}
	ret.SetDefault(&s)

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
