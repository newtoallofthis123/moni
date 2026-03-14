package format

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
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
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, strings.Join(headers, "\t"))
	// separator
	seps := make([]string, len(headers))
	for i, h := range headers {
		seps[i] = strings.Repeat("-", len(h))
	}
	fmt.Fprintln(tw, strings.Join(seps, "\t"))

	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}

	return tw.Flush()
}

// Message prints a simple message (for write commands that don't need formatted output).
func Message(format string, msg string, data any) error {
	if format == "json" {
		return outputJSON(os.Stdout, data)
	}
	fmt.Println(msg)
	return nil
}
