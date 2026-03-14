package cmd

import (
	"fmt"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/newtoallofthis123/moni/internal/format"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts",
}

var accountAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		acctType, _ := cmd.Flags().GetString("type")

		acct, err := db.AccountInsert(conn, name, acctType)
		if err != nil {
			return err
		}

		return format.Message(outputFormat,
			fmt.Sprintf("Account %q (%s) created.", acct.Name, acct.Type),
			acct,
		)
	},
}

var accountListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		accounts, err := db.AccountList(conn)
		if err != nil {
			return err
		}

		if len(accounts) == 0 {
			fmt.Println("No accounts yet. Add one with: moni account add <name> --type <type>")
			return nil
		}

		headers := []string{"ID", "Name", "Type", "Balance"}
		rows := make([][]string, len(accounts))
		for i, a := range accounts {
			rows[i] = []string{
				fmt.Sprintf("%d", a.ID),
				a.Name,
				a.Type,
				fmt.Sprintf("%.2f", a.Balance),
			}
		}

		return format.Output(outputFormat, headers, rows, accounts)
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(accountAddCmd)
	accountCmd.AddCommand(accountListCmd)

	accountAddCmd.Flags().StringP("type", "t", "bank", "Account type: bank, cash, credit, wallet, other")
}
