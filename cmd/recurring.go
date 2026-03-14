package cmd

import (
	"fmt"
	"strconv"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/newtoallofthis123/moni/internal/format"
	"github.com/spf13/cobra"
)

var recurringCmd = &cobra.Command{
	Use:   "recurring",
	Short: "Manage recurring items",
}

var recurringAddCmd = &cobra.Command{
	Use:   "add <description> <amount>",
	Short: "Add a recurring item",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		desc := args[0]
		amount, err := strconv.ParseFloat(args[1], 64)
		if err != nil || amount <= 0 {
			return fmt.Errorf("amount must be a positive number, got %q", args[1])
		}

		catName, _ := cmd.Flags().GetString("cat")
		frequency, _ := cmd.Flags().GetString("every")
		dueDay, _ := cmd.Flags().GetInt("due")
		txnType, _ := cmd.Flags().GetString("type")

		// Validate frequency
		switch frequency {
		case "daily", "weekly", "monthly", "yearly":
		default:
			return fmt.Errorf("--every must be daily, weekly, monthly, or yearly, got %q", frequency)
		}

		// Validate type
		if txnType != "expense" && txnType != "income" {
			return fmt.Errorf("--type must be 'expense' or 'income', got %q", txnType)
		}

		// Resolve category
		var categoryID *int64
		var categoryDisplay string
		if catName != "" {
			cat, err := db.CategoryGetByName(conn, catName)
			if err != nil {
				return fmt.Errorf("category %q not found", catName)
			}
			categoryID = &cat.ID
			categoryDisplay = cat.Name
		}

		r, err := db.RecurringInsert(conn, desc, amount, categoryID, frequency, dueDay, txnType)
		if err != nil {
			return err
		}
		r.CategoryName = categoryDisplay

		return format.Message(outputFormat,
			fmt.Sprintf("Recurring %s %q (%.2f, %s, due day %d) added.", r.Type, r.Description, r.Amount, r.Frequency, r.DueDay),
			r,
		)
	},
}

var recurringListCmd = &cobra.Command{
	Use:   "list",
	Short: "List active recurring items",
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := db.RecurringList(conn)
		if err != nil {
			return err
		}

		if len(items) == 0 {
			fmt.Println("No active recurring items.")
			return nil
		}

		headers := []string{"ID", "Description", "Amount", "Category", "Frequency", "Due Day", "Type"}
		rows := make([][]string, len(items))
		for i, r := range items {
			rows[i] = []string{
				fmt.Sprintf("%d", r.ID),
				r.Description,
				fmt.Sprintf("%.2f", r.Amount),
				r.CategoryName,
				r.Frequency,
				fmt.Sprintf("%d", r.DueDay),
				r.Type,
			}
		}

		return format.Output(outputFormat, headers, rows, items)
	},
}

func init() {
	rootCmd.AddCommand(recurringCmd)
	recurringCmd.AddCommand(recurringAddCmd)
	recurringCmd.AddCommand(recurringListCmd)

	recurringAddCmd.Flags().StringP("cat", "c", "", "Category name")
	recurringAddCmd.Flags().StringP("every", "e", "", "Frequency: daily, weekly, monthly, yearly")
	recurringAddCmd.Flags().IntP("due", "d", 1, "Due day (day of week/month/year)")
	recurringAddCmd.Flags().StringP("type", "t", "expense", "Type: expense or income")
	recurringAddCmd.MarkFlagRequired("every")
}
