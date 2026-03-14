package cmd

import (
	"fmt"
	"strconv"
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
			format.Empty(outputFormat, "No transactions found.")
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

		if interactive {
			return format.OutputInteractive(headers, rows)
		}
		return format.Output(outputFormat, headers, rows, txns)
	},
}

var transactionDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a transaction (reverses balance)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid transaction ID %q", args[0])
		}

		txn, err := db.TransactionDelete(conn, id)
		if err != nil {
			return err
		}

		return format.Message(outputFormat,
			fmt.Sprintf("Transaction #%d deleted (%.2f %s reversed from %s).", txn.ID, txn.Amount, txn.Type, txn.AccountName),
			txn,
		)
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

var transactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "Manage transactions",
}

func init() {
	rootCmd.AddCommand(transactionsCmd)
	rootCmd.AddCommand(transactionCmd)
	transactionCmd.AddCommand(transactionDeleteCmd)
	transactionsCmd.Flags().StringP("cat", "c", "", "Filter by category name")
	transactionsCmd.Flags().StringP("since", "s", "", "Filter by date: today, week, month, year, or YYYY-MM-DD")
}
