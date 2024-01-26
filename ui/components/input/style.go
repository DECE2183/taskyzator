package input

import (
	"taskyzator/ui/style"

	"github.com/charmbracelet/lipgloss"
)

var (
	inputTextStyle  = lipgloss.NewStyle().MarginBottom(1)
	inputTitleStyle = inputTextStyle.Copy().Foreground(style.NormalColor)
	inputBoxStyle   = style.RoundedBorder.Copy().BorderForeground(style.ArchivedTaskTextColor).Width(90)
)
