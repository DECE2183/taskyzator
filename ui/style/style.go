package style

import (
	"taskyzator/config"

	"github.com/charmbracelet/lipgloss"
)

const (
	DoneIcon    = "✔"
	ArchiveIcon = "▪"
)

var (
	ErrorColor            = lipgloss.Color(config.Current.Style.ErrorColor)
	AccentColor           = lipgloss.Color(config.Current.Style.AccentColor)
	DimmedColor           = lipgloss.Color(config.Current.Style.DimmedColor)
	NormalColor           = lipgloss.Color(config.Current.Style.NormalColor)
	ActiveTaskTextColor   = lipgloss.Color(config.Current.Style.ActiveTaskTextColor)
	DoneTaskTextColor     = lipgloss.Color(config.Current.Style.DoneTaskTextColor)
	ArchivedTaskTextColor = lipgloss.Color(config.Current.Style.ArchivedTaskTextColor)
)

var (
	NormalText    = lipgloss.NewStyle().Foreground(NormalColor)
	DimmedText    = lipgloss.NewStyle().Foreground(DimmedColor)
	RoundedBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(AccentColor).Padding(0, 1)
)
