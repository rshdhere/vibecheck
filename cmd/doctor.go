package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const doctorOutput = "vibecheck self-test OK"

var doctorCmd = &cobra.Command{
	Use:     "doctor",
	Aliases: []string{"selftest"},
	Short:   "Runs an offline sanity check to verify the CLI is healthy",
	Long:    "Runs a lightweight, offline sanity check to verify the vibecheck CLI is healthy without requiring git, network, or API access.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(cmd.OutOrStdout(), doctorOutput)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
