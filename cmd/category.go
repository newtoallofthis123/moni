package cmd

import (
	"fmt"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/newtoallofthis123/moni/internal/format"
	"github.com/spf13/cobra"
)

var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Manage categories",
}

var categoryAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new category",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		catType, _ := cmd.Flags().GetString("type")

		cat, err := db.CategoryInsert(conn, name, catType)
		if err != nil {
			return err
		}

		return format.Message(outputFormat,
			fmt.Sprintf("Category %q (%s) created.", cat.Name, cat.Type),
			cat,
		)
	},
}

var categoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all categories",
	RunE: func(cmd *cobra.Command, args []string) error {
		categories, err := db.CategoryList(conn)
		if err != nil {
			return err
		}

		if len(categories) == 0 {
			format.Empty(outputFormat, "No categories yet. Add one with: moni category add <name>")
			return nil
		}

		headers := []string{"ID", "Name", "Type"}
		rows := make([][]string, len(categories))
		for i, c := range categories {
			rows[i] = []string{
				fmt.Sprintf("%d", c.ID),
				c.Name,
				c.Type,
			}
		}

		if interactive {
			return format.OutputInteractive(headers, rows)
		}
		return format.Output(outputFormat, headers, rows, categories)
	},
}

func init() {
	rootCmd.AddCommand(categoryCmd)
	categoryCmd.AddCommand(categoryAddCmd)
	categoryCmd.AddCommand(categoryListCmd)

	categoryAddCmd.Flags().StringP("type", "t", "expense", "Category type: expense, income, both")
}
