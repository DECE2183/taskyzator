package style

import "github.com/charmbracelet/lipgloss"

const (
	DoneIcon    = "✔"
	ArchiveIcon = "▪"
)

var (
	ErrorColor            = lipgloss.Color("#F33")
	AccentColor           = lipgloss.Color("#0cc47e")
	DimmedColor           = lipgloss.Color("#a3a3a3")
	ActiveTaskTextColor   = lipgloss.Color("#f7fffc")
	DoneTaskTextColor     = lipgloss.Color("#6ca690")
	ArchivedTaskTextColor = lipgloss.Color("#7d7d7d")
	NormalTextColor       = lipgloss.Color("#FFFFFF")
)

var (
	RoundedBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(AccentColor).Padding(0, 1)
)
