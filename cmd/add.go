package cmd

import (
	"fmt"
	"strconv"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/newtoallofthis123/moni/internal/format"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a transaction (expense or income)",
}

var addExpenseCmd = &cobra.Command{
	Use:   "expense <amount>",
	Short: "Log an expense",
	Args:  cobra.ExactArgs(1),
	RunE:  runAddTransaction("expense"),
}

var addIncomeCmd = &cobra.Command{
	Use:   "income <amount>",
	Short: "Log an income",
	Args:  cobra.ExactArgs(1),
	RunE:  runAddTransaction("income"),
}

func runAddTransaction(txnType string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		amount, err := strconv.ParseFloat(args[0], 64)
		if err != nil || amount <= 0 {
			return fmt.Errorf("amount must be a positive number, got %q", args[0])
		}

		catName, _ := cmd.Flags().GetString("cat")
		note, _ := cmd.Flags().GetString("note")
		acctName, _ := cmd.Flags().GetString("account")

		// Resolve account
		var acct, acctErr = resolveAccount(acctName)
		if acctErr != nil {
			return acctErr
		}

		// Resolve category
		var catID *int64
		if catName != "" {
			cat, err := db.CategoryGetByName(conn, catName)
			if err != nil {
				return fmt.Errorf("category %q not found: %w", catName, err)
			}
			catID = &cat.ID
		}

		txn, err := db.TransactionInsert(conn, acct.ID, catID, txnType, amount, note)
		if err != nil {
			return err
		}

		txn.AccountName = acct.Name
		catDisplay := ""
		if catName != "" {
			catDisplay = fmt.Sprintf(" [%s]", catName)
			txn.CategoryName = catName
		}

		return format.Message(outputFormat,
			fmt.Sprintf("%s of %.2f%s recorded in %s (txn #%d)", txnType, amount, catDisplay, acct.Name, txn.ID),
			txn,
		)
	}
}

func resolveAccount(name string) (acct struct {
	ID   int64
	Name string
}, err error) {
	if name != "" {
		a, e := db.AccountGetByName(conn, name)
		if e != nil {
			return acct, fmt.Errorf("account %q not found: %w", name, e)
		}
		acct.ID = a.ID
		acct.Name = a.Name
		return acct, nil
	}
	a, e := db.AccountGetFirst(conn)
	if e != nil {
		return acct, fmt.Errorf("no accounts found — add one with: moni account add <name>")
	}
	acct.ID = a.ID
	acct.Name = a.Name
	return acct, nil
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addExpenseCmd)
	addCmd.AddCommand(addIncomeCmd)

	for _, cmd := range []*cobra.Command{addExpenseCmd, addIncomeCmd} {
		cmd.Flags().StringP("cat", "c", "", "Category name")
		cmd.Flags().StringP("note", "n", "", "Transaction note/description")
		cmd.Flags().StringP("account", "a", "", "Account name (defaults to first account)")
	}
}
