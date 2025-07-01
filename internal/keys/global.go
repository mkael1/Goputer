package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type global struct {
	Quit key.Binding
	Help key.Binding
}

var Global = global{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c|q", "exit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}
