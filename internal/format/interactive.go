package format

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type interactiveModel struct {
	table table.Model
}

func newInteractiveModel(headers []string, rows [][]string) interactiveModel {
	columns := make([]table.Column, len(headers))
	for i, h := range headers {
		width := len(h)
		for _, r := range rows {
			if i < len(r) && len(r[i]) > width {
				width = len(r[i])
			}
		}
		columns[i] = table.Column{Title: h, Width: width}
	}

	tableRows := make([]table.Row, len(rows))
	for i, r := range rows {
		tableRows[i] = table.Row(r)
	}

	s := table.DefaultStyles()
	s.Header = HeaderStyle.Padding(0, 1)
	s.Selected = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("4"))

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(tableRows),
		table.WithFocused(true),
		table.WithHeight(20),
		table.WithStyles(s),
	)

	return interactiveModel{table: t}
}

func (m interactiveModel) Init() tea.Cmd { return nil }

func (m interactiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m interactiveModel) View() string {
	return m.table.View() + "\n  q/esc to quit\n"
}
