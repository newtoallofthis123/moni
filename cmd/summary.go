package cmd

import (
	"fmt"
	"time"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/newtoallofthis123/moni/internal/format"
	"github.com/spf13/cobra"
)

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Spending/income summary for a month",
	RunE: func(cmd *cobra.Command, args []string) error {
		monthStr, _ := cmd.Flags().GetString("month")

		var month time.Time
		if monthStr == "" || monthStr == "current" {
			now := time.Now()
			month = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		} else {
			parsed, err := time.Parse("2006-01", monthStr)
			if err != nil {
				return fmt.Errorf("--month must be YYYY-MM or 'current', got %q", monthStr)
			}
			month = parsed
		}

		s, err := db.SummaryForMonth(conn, month)
		if err != nil {
			return err
		}

		// Main summary table
		headers := []string{"Metric", "Amount"}
		rows := [][]string{
			{"Income", fmt.Sprintf("%.2f", s.TotalIncome)},
			{"Expenses", fmt.Sprintf("%.2f", s.TotalExpenses)},
			{"Net", fmt.Sprintf("%.2f", s.Net)},
		}

		if outputFormat == "json" {
			return format.Output(outputFormat, nil, nil, s)
		}

		fmt.Printf("Summary for %s\n\n", s.Month)
		if err := format.Output(outputFormat, headers, rows, s); err != nil {
			return err
		}

		if len(s.TopCategories) > 0 {
			fmt.Println("\nTop Expense Categories:")
			catHeaders := []string{"Category", "Amount"}
			catRows := make([][]string, len(s.TopCategories))
			for i, c := range s.TopCategories {
				catRows[i] = []string{c.Category, fmt.Sprintf("%.2f", c.Amount)}
			}
			return format.Output(outputFormat, catHeaders, catRows, s.TopCategories)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(summaryCmd)

	summaryCmd.Flags().StringP("month", "m", "", "Month in YYYY-MM format (default: current)")
}
