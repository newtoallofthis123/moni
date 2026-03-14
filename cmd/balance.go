package cmd

import (
	"fmt"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/newtoallofthis123/moni/internal/format"
	"github.com/spf13/cobra"
)

// BalanceSummary is the JSON-friendly output for the balance command.
type BalanceSummary struct {
	Accounts []AccountBalance `json:"accounts"`
	Total    float64          `json:"total"`
}

// AccountBalance is a single account's balance for the balance command.
type AccountBalance struct {
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Show all account balances",
	RunE: func(cmd *cobra.Command, args []string) error {
		accounts, err := db.AccountList(conn)
		if err != nil {
			return err
		}

		if len(accounts) == 0 {
			fmt.Println("No accounts yet. Add one with: moni account add <name> --type <type>")
			return nil
		}

		var total float64
		balances := make([]AccountBalance, len(accounts))
		for i, a := range accounts {
			balances[i] = AccountBalance{Name: a.Name, Type: a.Type, Balance: a.Balance}
			total += a.Balance
		}

		summary := BalanceSummary{Accounts: balances, Total: total}

		headers := []string{"Account", "Type", "Balance"}
		rows := make([][]string, len(accounts)+1)
		for i, b := range balances {
			rows[i] = []string{b.Name, b.Type, fmt.Sprintf("%.2f", b.Balance)}
		}
		rows[len(accounts)] = []string{"TOTAL", "", fmt.Sprintf("%.2f", total)}

		return format.Output(outputFormat, headers, rows, summary)
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}
