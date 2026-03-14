package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/newtoallofthis123/moni/internal/format"
	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link <txn_id>",
	Short: "Link a transaction to persons",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		txnID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("transaction ID must be a number, got %q", args[0])
		}

		personsFlag, _ := cmd.Flags().GetStringSlice("persons")
		if len(personsFlag) == 0 {
			return fmt.Errorf("--persons is required")
		}

		note, _ := cmd.Flags().GetString("note")

		// Verify transaction exists
		exists, err := db.TransactionExists(conn, txnID)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("transaction %d not found", txnID)
		}

		// Link each person
		var linked []string
		for _, name := range personsFlag {
			name = strings.TrimSpace(name)
			p, err := db.PersonGetByName(conn, name)
			if err != nil {
				return fmt.Errorf("person %q not found — add them first with: moni person add %s", name, name)
			}

			tp, err := db.TransactionPersonLink(conn, txnID, p.ID, note)
			if err != nil {
				return err
			}
			tp.PersonName = p.Name
			linked = append(linked, p.Name)
		}

		result := struct {
			TransactionID int64    `json:"transaction_id"`
			Persons       []string `json:"persons"`
			Note          string   `json:"note,omitempty"`
		}{TransactionID: txnID, Persons: linked, Note: note}

		return format.Message(outputFormat,
			fmt.Sprintf("Transaction %d linked to %s.", txnID, strings.Join(linked, ", ")),
			result,
		)
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)

	linkCmd.Flags().StringSliceP("persons", "p", nil, "Person names to link (comma-separated)")
	linkCmd.Flags().StringP("note", "n", "", "Note for the link")
	linkCmd.MarkFlagRequired("persons")
}
