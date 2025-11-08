// Package cmd is provided by cobra-cli to ship command-line tools faster
/*
Copyright © 2025 raashed
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// version will be set at build time via ldflags
var version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "vibecheck",
	Short:   "A command-line tool for easing git commit messages for me(or may be you guys too lol)",
	Long:    `A complete solution for vibecoders to vibecheck their code and save it locally even before it messess-up your production, vibecheck is a check point were they can automate their commit message to models like gpt-oss:20b, GPT4o-mini`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vibecheck.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Make version flag available for all subcommands
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}
