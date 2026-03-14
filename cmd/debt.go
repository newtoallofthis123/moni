package cmd

import (
	"fmt"
	"strconv"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/newtoallofthis123/moni/internal/format"
	"github.com/spf13/cobra"
)

var debtCmd = &cobra.Command{
	Use:   "debt",
	Short: "Manage debts",
}

var debtAddCmd = &cobra.Command{
	Use:   "add <person> <amount> <i_owe|they_owe>",
	Short: "Record a debt",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		personName := args[0]
		amount, err := strconv.ParseFloat(args[1], 64)
		if err != nil || amount <= 0 {
			return fmt.Errorf("amount must be a positive number, got %q", args[1])
		}
		direction := args[2]
		if direction != "i_owe" && direction != "they_owe" {
			return fmt.Errorf("direction must be 'i_owe' or 'they_owe', got %q", direction)
		}

		note, _ := cmd.Flags().GetString("note")

		p, err := db.PersonGetByName(conn, personName)
		if err != nil {
			return fmt.Errorf("person %q not found — add them first with: moni person add %s", personName, personName)
		}

		d, err := db.DebtInsert(conn, p.ID, amount, direction, note)
		if err != nil {
			return err
		}
		d.PersonName = p.Name

		return format.Message(outputFormat,
			fmt.Sprintf("Debt of %.2f (%s) with %s recorded.", d.Amount, d.Direction, p.Name),
			d,
		)
	},
}

var debtSettleCmd = &cobra.Command{
	Use:   "settle <person> <amount>",
	Short: "Settle debt with a person",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		personName := args[0]
		amount, err := strconv.ParseFloat(args[1], 64)
		if err != nil || amount <= 0 {
			return fmt.Errorf("amount must be a positive number, got %q", args[1])
		}

		p, err := db.PersonGetByName(conn, personName)
		if err != nil {
			return fmt.Errorf("person %q not found", personName)
		}

		settled, err := db.DebtSettle(conn, p.ID, amount)
		if err != nil {
			return err
		}

		result := struct {
			Person  string  `json:"person"`
			Settled float64 `json:"settled"`
		}{Person: p.Name, Settled: settled}

		return format.Message(outputFormat,
			fmt.Sprintf("Settled %.2f with %s.", settled, p.Name),
			result,
		)
	},
}

var debtListCmd = &cobra.Command{
	Use:   "list",
	Short: "List open debts",
	RunE: func(cmd *cobra.Command, args []string) error {
		debts, err := db.DebtListOpen(conn)
		if err != nil {
			return err
		}

		if len(debts) == 0 {
			format.Empty(outputFormat, "No open debts.")
			return nil
		}

		headers := []string{"ID", "Person", "Amount", "Direction", "Note"}
		rows := make([][]string, len(debts))
		for i, d := range debts {
			rows[i] = []string{
				fmt.Sprintf("%d", d.ID),
				d.PersonName,
				fmt.Sprintf("%.2f", d.Amount),
				d.Direction,
				d.Note,
			}
		}

		if interactive {
			return format.OutputInteractive(headers, rows)
		}
		return format.Output(outputFormat, headers, rows, debts)
	},
}

var debtDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a debt",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid debt ID %q", args[0])
		}

		if err := db.DebtDelete(conn, id); err != nil {
			return err
		}

		return format.Message(outputFormat,
			fmt.Sprintf("Debt #%d deleted.", id),
			struct {
				ID     int64  `json:"id"`
				Status string `json:"status"`
			}{ID: id, Status: "deleted"},
		)
	},
}

func init() {
	rootCmd.AddCommand(debtCmd)
	debtCmd.AddCommand(debtAddCmd)
	debtCmd.AddCommand(debtSettleCmd)
	debtCmd.AddCommand(debtListCmd)
	debtCmd.AddCommand(debtDeleteCmd)

	debtAddCmd.Flags().StringP("note", "n", "", "Note about this debt")
}
