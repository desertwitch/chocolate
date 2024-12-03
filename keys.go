package chocolate

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines keybindings. It satisfies the help.KeyMap interface, which
// is used to render the menu
type KeyMap struct {
	// Keybinds used by chocolate
	NextBar key.Binding
	PrevBar key.Binding
	Child   key.Binding
	Parent  key.Binding
}

// DefaultKeyMap returns a default set of keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		NextBar: key.NewBinding(
			key.WithKeys("tab", "n"),
			key.WithHelp("tab/n", "Go to next"),
		),
		PrevBar: key.NewBinding(
			key.WithKeys("shift+tab", "p"),
			key.WithHelp("shift+tab/p", "Go to prev"),
		),
		Child: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Go to child"),
		),
		Parent: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Go to parent"),
		),
	}
}
