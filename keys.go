package chocolate

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines keybindings. It satisfies the help.KeyMap interface, which
// is used to render the menu
type KeyMap struct {
	// Keybinds used by chocolate
	Quit    key.Binding
	NextBar key.Binding
	PrevBar key.Binding
	Focus   key.Binding
	Release key.Binding
}

// DefaultKeyMap returns a default set of keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "Quit"),
		),
		NextBar: key.NewBinding(
			key.WithKeys("tab", "n"),
			key.WithHelp("tab/n", "Go to next"),
		),
		PrevBar: key.NewBinding(
			key.WithKeys("shift+tab", "p"),
			key.WithHelp("shift+tab/p", "Go to prev"),
		),
		Focus: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Hand input over to selected"),
		),
		Release: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Release input from selected"),
		),
	}
}
