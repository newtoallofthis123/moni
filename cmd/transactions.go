package cmd

import (
	"fmt"
	"time"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/newtoallofthis123/moni/internal/format"
	"github.com/spf13/cobra"
)

var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "List transactions",
	RunE: func(cmd *cobra.Command, args []string) error {
		catName, _ := cmd.Flags().GetString("cat")
		sinceStr, _ := cmd.Flags().GetString("since")

		since, err := parseSince(sinceStr)
		if err != nil {
			return err
		}

		txns, err := db.TransactionList(conn, catName, since)
		if err != nil {
			return err
		}

		if len(txns) == 0 {
			fmt.Println("No transactions found.")
			return nil
		}

		headers := []string{"ID", "Date", "Type", "Amount", "Category", "Note", "Account"}
		rows := make([][]string, len(txns))
		for i, t := range txns {
			rows[i] = []string{
				fmt.Sprintf("%d", t.ID),
				t.Date.Format("2006-01-02"),
				t.Type,
				fmt.Sprintf("%.2f", t.Amount),
				t.CategoryName,
				t.Note,
				t.AccountName,
			}
		}

		return format.Output(outputFormat, headers, rows, txns)
	},
}

func parseSince(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}

	now := time.Now()
	switch s {
	case "today":
		y, m, d := now.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, now.Location()), nil
	case "week":
		y, m, d := now.Date()
		today := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		weekday := int(today.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		return today.AddDate(0, 0, -(weekday - 1)), nil
	case "month":
		y, m, _ := now.Date()
		return time.Date(y, m, 1, 0, 0, 0, 0, now.Location()), nil
	case "year":
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()), nil
	default:
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid --since value %q: use today, week, month, year, or YYYY-MM-DD", s)
		}
		return t, nil
	}
}

func init() {
	rootCmd.AddCommand(transactionsCmd)
	transactionsCmd.Flags().StringP("cat", "c", "", "Filter by category name")
	transactionsCmd.Flags().StringP("since", "s", "", "Filter by date: today, week, month, year, or YYYY-MM-DD")
}
