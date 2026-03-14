package cmd

import (
	"fmt"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/newtoallofthis123/moni/internal/format"
	"github.com/spf13/cobra"
)

var personCmd = &cobra.Command{
	Use:   "person",
	Short: "Manage persons",
}

var personAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a person",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		phone, _ := cmd.Flags().GetString("phone")

		p, err := db.PersonInsert(conn, name, phone)
		if err != nil {
			return err
		}

		return format.Message(outputFormat,
			fmt.Sprintf("Person %q added.", p.Name),
			p,
		)
	},
}

var personHistoryCmd = &cobra.Command{
	Use:   "history <name>",
	Short: "Show a person's transactions and debts",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		p, err := db.PersonGetByName(conn, name)
		if err != nil {
			return err
		}

		history, err := db.PersonHistory(conn, p.ID)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return format.Output(outputFormat, nil, nil, history)
		}

		// Text/table: show transactions then debts
		if len(history.Transactions) > 0 {
			fmt.Println("Transactions:")
			headers := []string{"ID", "Date", "Type", "Amount", "Category", "Account", "Note"}
			rows := make([][]string, len(history.Transactions))
			for i, t := range history.Transactions {
				rows[i] = []string{
					fmt.Sprintf("%d", t.ID),
					t.Date.Format("2006-01-02"),
					t.Type,
					fmt.Sprintf("%.2f", t.Amount),
					t.CategoryName,
					t.AccountName,
					t.Note,
				}
			}
			if err := format.Output(outputFormat, headers, rows, history.Transactions); err != nil {
				return err
			}
		} else {
			fmt.Println("No linked transactions.")
		}

		fmt.Println()

		if len(history.Debts) > 0 {
			fmt.Println("Debts:")
			headers := []string{"ID", "Amount", "Direction", "Settled", "Note"}
			rows := make([][]string, len(history.Debts))
			for i, d := range history.Debts {
				settled := "no"
				if d.Settled {
					settled = "yes"
				}
				rows[i] = []string{
					fmt.Sprintf("%d", d.ID),
					fmt.Sprintf("%.2f", d.Amount),
					d.Direction,
					settled,
					d.Note,
				}
			}
			if err := format.Output(outputFormat, headers, rows, history.Debts); err != nil {
				return err
			}
		} else {
			fmt.Println("No debts.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(personCmd)
	personCmd.AddCommand(personAddCmd)
	personCmd.AddCommand(personHistoryCmd)

	personAddCmd.Flags().StringP("phone", "p", "", "Phone number")
}
