package cmd

import (
	"fmt"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize moni database",
	Long:  "Creates ~/.moni/ directory and database with all tables and default categories.",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := db.Open()
		if err != nil {
			return fmt.Errorf("cannot create database: %w", err)
		}
		defer conn.Close()

		if err := db.Migrate(conn); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}

		dbPath, _ := db.DBPath()
		fmt.Printf("moni initialized at %s\n", dbPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
