/*
Copyright © 2025 raashed
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "A command-line tool for easing git commit messages for me(or may be you guys too lol)",
	Long:  `A complete solution for vibecoders to vibecheck their code and save it locally even before it messess-up your production, vibecheck is a check point were they can automate their commit message to models like Ollama, GPT-5, Sonnet-4.5, Qwen-3 etc`,
}

func init() {
	rootCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
