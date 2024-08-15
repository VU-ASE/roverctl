package components

import "github.com/charmbracelet/bubbles/spinner"

func InitializeSpinner() spinner.Model {
	spin := spinner.New()
	spin.Spinner = spinner.Hamburger
	return spin
}
