package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type panes struct {
	Right key.Binding
	Left  key.Binding
}

var Panes = panes{
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "next pane"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "previous pane"),
	),
}
