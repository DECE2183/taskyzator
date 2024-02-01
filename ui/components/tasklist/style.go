package tasklist

import (
	"github.com/dece2183/taskyzator/ui/style"

	"github.com/charmbracelet/lipgloss"
)

var (
	selectedItemStyle = style.RoundedBorder.Copy().BorderForeground(style.DimmedColor)
	itemStyle         = lipgloss.NewStyle().Margin(1).MarginLeft(2)
	titleStyle        = lipgloss.NewStyle()

	activeTaskTextStyle   = lipgloss.NewStyle().Foreground(style.ActiveTaskTextColor)
	doneTaskTextStyle     = lipgloss.NewStyle().Foreground(style.DoneTaskTextColor).Strikethrough(true).StrikethroughSpaces(true)
	archivedTaskTextStyle = lipgloss.NewStyle().Foreground(style.ArchivedTaskTextColor)
)
