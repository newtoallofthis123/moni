package cmd

import (
	"fmt"
	"strconv"

	"github.com/newtoallofthis123/moni/internal/db"
	"github.com/newtoallofthis123/moni/internal/format"
	"github.com/spf13/cobra"
)

var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Manage savings buckets",
}

var bucketCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a savings bucket",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		target, _ := cmd.Flags().GetFloat64("target")
		if target <= 0 {
			return fmt.Errorf("--target must be a positive number")
		}

		b, err := db.BucketInsert(conn, name, target)
		if err != nil {
			return err
		}

		return format.Message(outputFormat,
			fmt.Sprintf("Bucket %q created with target %.2f.", b.Name, b.Target),
			b,
		)
	},
}

var bucketAddCmd = &cobra.Command{
	Use:   "add <name> <amount>",
	Short: "Add funds to a bucket",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		amount, err := strconv.ParseFloat(args[1], 64)
		if err != nil || amount <= 0 {
			return fmt.Errorf("amount must be a positive number, got %q", args[1])
		}

		bucket, err := db.BucketGetByName(conn, name)
		if err != nil {
			return fmt.Errorf("bucket %q not found — create it first with: moni bucket create %s --target <amount>", name, name)
		}

		b, err := db.BucketAddFunds(conn, bucket.ID, amount)
		if err != nil {
			return err
		}

		pct := (b.Current / b.Target) * 100
		return format.Message(outputFormat,
			fmt.Sprintf("Added %.2f to %q — now at %.2f/%.2f (%.0f%%).", amount, b.Name, b.Current, b.Target, pct),
			b,
		)
	},
}

var bucketStatusCmd = &cobra.Command{
	Use:   "status [name]",
	Short: "Show bucket progress",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			b, err := db.BucketGetByName(conn, args[0])
			if err != nil {
				return err
			}
			pct := (b.Current / b.Target) * 100
			headers := []string{"Name", "Current", "Target", "Progress"}
			rows := [][]string{{
				b.Name,
				fmt.Sprintf("%.2f", b.Current),
				fmt.Sprintf("%.2f", b.Target),
				fmt.Sprintf("%.0f%%", pct),
			}}
			if interactive {
				return format.OutputInteractive(headers, rows)
			}
			return format.Output(outputFormat, headers, rows, b)
		}

		buckets, err := db.BucketList(conn)
		if err != nil {
			return err
		}
		if len(buckets) == 0 {
			format.Empty(outputFormat, "No buckets.")
			return nil
		}

		headers := []string{"Name", "Current", "Target", "Progress"}
		rows := make([][]string, len(buckets))
		for i, b := range buckets {
			pct := (b.Current / b.Target) * 100
			rows[i] = []string{
				b.Name,
				fmt.Sprintf("%.2f", b.Current),
				fmt.Sprintf("%.2f", b.Target),
				fmt.Sprintf("%.0f%%", pct),
			}
		}
		if interactive {
			return format.OutputInteractive(headers, rows)
		}
		return format.Output(outputFormat, headers, rows, buckets)
	},
}

var bucketDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a bucket",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		bucket, err := db.BucketGetByName(conn, name)
		if err != nil {
			return err
		}

		if err := db.BucketDelete(conn, bucket.ID); err != nil {
			return err
		}

		return format.Message(outputFormat,
			fmt.Sprintf("Bucket %q deleted.", name),
			bucket,
		)
	},
}

var bucketEditCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit a bucket's target",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		bucket, err := db.BucketGetByName(conn, name)
		if err != nil {
			return err
		}

		newName := name
		if cmd.Flags().Changed("name") {
			newName, _ = cmd.Flags().GetString("name")
		}
		target := bucket.Target
		if cmd.Flags().Changed("target") {
			target, _ = cmd.Flags().GetFloat64("target")
			if target <= 0 {
				return fmt.Errorf("--target must be a positive number")
			}
		}

		b, err := db.BucketEdit(conn, bucket.ID, newName, target)
		if err != nil {
			return err
		}

		return format.Message(outputFormat,
			fmt.Sprintf("Bucket %q updated (target: %.2f).", b.Name, b.Target),
			b,
		)
	},
}

func init() {
	rootCmd.AddCommand(bucketCmd)
	bucketCmd.AddCommand(bucketCreateCmd)
	bucketCmd.AddCommand(bucketAddCmd)
	bucketCmd.AddCommand(bucketStatusCmd)
	bucketCmd.AddCommand(bucketDeleteCmd)
	bucketCmd.AddCommand(bucketEditCmd)

	bucketCreateCmd.Flags().Float64P("target", "t", 0, "Target amount for the bucket")
	bucketCreateCmd.MarkFlagRequired("target")

	bucketEditCmd.Flags().StringP("name", "n", "", "New bucket name")
	bucketEditCmd.Flags().Float64P("target", "t", 0, "New target amount")
}
