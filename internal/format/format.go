package format

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/mattn/go-isatty"
)

// Output writes data in the requested format to stdout.
// headers: column names for table/text modes.
// rows: each row is a slice of string values matching headers.
// data: the original typed data for JSON output.
func Output(format string, headers []string, rows [][]string, data any) error {
	switch format {
	case "json":
		return outputJSON(os.Stdout, data)
	case "text":
		return outputText(os.Stdout, headers, rows)
	case "table":
		return outputTable(os.Stdout, headers, rows)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputJSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func outputText(w io.Writer, headers []string, rows [][]string) error {
	for _, row := range rows {
		for i, val := range row {
			if i < len(headers) {
				fmt.Fprintf(w, "%s: %s\n", headers[i], val)
			}
		}
		fmt.Fprintln(w)
	}
	return nil
}

func outputTable(w io.Writer, headers []string, rows [][]string) error {
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(BorderStyle).
		BorderRow(true).
		Headers(headers...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			s := lipgloss.NewStyle().PaddingRight(1).PaddingLeft(1)
			if row == table.HeaderRow {
				return HeaderStyle.PaddingRight(1).PaddingLeft(1)
			}
			if row%2 == 0 {
				return s
			}
			return s.Foreground(lipgloss.Color("245"))
		})

	fmt.Fprintln(w, t.Render())
	return nil
}

// OutputInteractive launches a bubbletea TUI with the given data.
// Falls back to outputTable if stdout is not a TTY.
func OutputInteractive(headers []string, rows [][]string) error {
	if !isTerminal() {
		return outputTable(os.Stdout, headers, rows)
	}
	m := newInteractiveModel(headers, rows)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// Message prints a simple message (for write commands that don't need formatted output).
func Message(format string, msg string, data any) error {
	if format == "json" {
		return outputJSON(os.Stdout, data)
	}
	fmt.Println(SuccessStyle.Render(msg))
	return nil
}

// Label prints a styled section header (e.g. "Summary for March").
func Label(format string, msg string) {
	if format == "json" {
		return
	}
	fmt.Println(LabelStyle.Render(msg))
}

// Empty prints a styled empty-state hint (e.g. "No accounts yet...").
func Empty(format string, msg string) {
	if format == "json" {
		return
	}
	fmt.Println(EmptyStyle.Render(msg))
}

func isTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

