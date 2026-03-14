package cmd

import (
	"fmt"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/newtoallofthis123/moni/internal/format"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage database backups",
}

var backupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a backup of the database",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := db.BackupCreate(conn)
		if err != nil {
			return err
		}

		return format.Message(outputFormat,
			fmt.Sprintf("Backup created: %s", path),
			db.BackupInfo{Path: path},
		)
	},
}

var backupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all backups",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		infos, err := db.BackupList()
		if err != nil {
			return err
		}
		if len(infos) == 0 {
			fmt.Println("No backups found.")
			return nil
		}

		headers := []string{"Date", "Size"}
		rows := make([][]string, len(infos))
		for i, info := range infos {
			rows[i] = []string{info.Date, formatBytes(info.Size)}
		}
		return format.Output(outputFormat, headers, rows, infos)
	},
}

var backupDeleteCmd = &cobra.Command{
	Use:   "delete <date>",
	Short: "Delete a backup by date (e.g. 2026-03-14)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		date := args[0]
		if err := db.BackupDelete(date); err != nil {
			return err
		}

		return format.Message(outputFormat,
			fmt.Sprintf("Backup %q deleted.", date),
			db.BackupInfo{Date: date},
		)
	},
}

var backupPruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove old backups, keeping the newest N",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		keep, _ := cmd.Flags().GetInt("keep")

		removed, err := db.BackupPrune(keep)
		if err != nil {
			return err
		}

		if len(removed) == 0 {
			fmt.Println("Nothing to prune.")
			return nil
		}

		return format.Message(outputFormat,
			fmt.Sprintf("Pruned %d backup(s).", len(removed)),
			removed,
		)
	},
}

var backupRestoreCmd = &cobra.Command{
	Use:   "restore <date>",
	Short: "Restore database from a backup (e.g. 2026-03-14)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		date := args[0]
		if err := db.BackupRestore(conn, date); err != nil {
			return err
		}
		// conn was closed inside BackupRestore; prevent PersistentPostRun from double-closing.
		conn = nil

		return format.Message(outputFormat,
			fmt.Sprintf("Database restored from backup %q.", date),
			db.BackupInfo{Date: date},
		)
	},
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.AddCommand(backupCreateCmd)
	backupCmd.AddCommand(backupListCmd)
	backupCmd.AddCommand(backupDeleteCmd)
	backupCmd.AddCommand(backupPruneCmd)
	backupCmd.AddCommand(backupRestoreCmd)

	backupPruneCmd.Flags().IntP("keep", "k", 5, "Number of newest backups to keep")
}
