package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	conn         *sql.DB
)

var rootCmd = &cobra.Command{
	Use:   "moni",
	Short: "Personal finance CLI",
	Long:  "moni — a local-first personal finance tracker backed by SQLite.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip DB connection for init command (it handles its own)
		if cmd.Name() == "init" {
			return nil
		}

		var err error
		conn, err = db.Open()
		if err != nil {
			return fmt.Errorf("cannot open database: %w\nRun 'moni init' first.", err)
		}
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if conn != nil {
			conn.Close()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format: text, table, json")
}
