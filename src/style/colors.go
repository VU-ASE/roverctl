package style

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Lipgloss colors that can be easily reused
var AsePrimary = lipgloss.Color("#0077B3")
var WarningPrimary = lipgloss.Color("#FFA500")
var SuccessPrimary = lipgloss.Color("#008000")
var ErrorPrimary = lipgloss.Color("#FF0000")
var GrayPrimary = lipgloss.Color("#808080")

// Lipgloss classes that can be easily reused
var Title = lipgloss.NewStyle().Foreground(AsePrimary)
var Subtitle = lipgloss.NewStyle().Foreground(GrayPrimary)
var Error = lipgloss.NewStyle().Foreground(ErrorPrimary)
var Warning = lipgloss.NewStyle().Foreground(WarningPrimary)
var Success = lipgloss.NewStyle().Foreground(SuccessPrimary)
var Gray = lipgloss.NewStyle().Foreground(GrayPrimary)

// Form theme
var FormTheme = huh.ThemeBase()

func RenderColor(view string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(color).Render(view)
}
