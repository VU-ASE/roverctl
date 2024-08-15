package components

import "github.com/charmbracelet/bubbles/textinput"

type Textfield struct {
	Input textinput.Model
	Err   error
}

func InitializeTextfield() Textfield {
	ti := textinput.New()
	ti.Placeholder = "..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return Textfield{
		Input: ti,
		Err:   nil,
	}
}
