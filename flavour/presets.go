package flavour

import "github.com/charmbracelet/lipgloss"

var DefaultColors = ColorDefinitons{
	COLOR_PRIMARY:      lipgloss.Color("246"),
	COLOR_PRIMARY_BG:   lipgloss.Color("232"),
	COLOR_SECONDARY:    lipgloss.Color("15"),
	COLOR_SECONDARY_BG: lipgloss.Color("237"),

	// sweets/notify
	"Info":  lipgloss.Color("10"),
	"Warn":  lipgloss.Color("11"),
	"Error": lipgloss.Color("9"),
}
