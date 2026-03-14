package format

import "github.com/charmbracelet/lipgloss"

var (
	HeaderStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4"))
	EvenRowStyle = lipgloss.NewStyle()
	OddRowStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	BorderStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	EmptyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true)
	LabelStyle   = lipgloss.NewStyle().Bold(true).Underline(true)
)
