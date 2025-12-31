// Package cmd is provided by cobra-cli to ship command-line tools faster
/*
Copyright © 2025 raashed
*/
package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rshdhere/vibecheck/internal/config"
	"github.com/rshdhere/vibecheck/internal/git"
	"github.com/rshdhere/vibecheck/internal/llm"
	_ "github.com/rshdhere/vibecheck/internal/llm/anthropic"
	_ "github.com/rshdhere/vibecheck/internal/llm/deepseek"
	_ "github.com/rshdhere/vibecheck/internal/llm/gemini"
	_ "github.com/rshdhere/vibecheck/internal/llm/grok"
	_ "github.com/rshdhere/vibecheck/internal/llm/groq"
	_ "github.com/rshdhere/vibecheck/internal/llm/kimi"
	_ "github.com/rshdhere/vibecheck/internal/llm/ollama"
	_ "github.com/rshdhere/vibecheck/internal/llm/openai"
	_ "github.com/rshdhere/vibecheck/internal/llm/perplexity"
	_ "github.com/rshdhere/vibecheck/internal/llm/qwen"
	"github.com/rshdhere/vibecheck/internal/stats"
	"github.com/rshdhere/vibecheck/internal/ui/notify"
	"github.com/spf13/cobra"
)

const (
	promptFlagName   = "prompt"
	providerFlagName = "provider"
)

type ProviderFunc func(context.Context, string, string) (string, error)

var commitCmd = &cobra.Command{
	Use:     "commit",
	Short:   "A command-line tool for easing git commit messages for me(or may be you guys too lol), adding multiple models to it sounds cool right?!",
	Long:    `A complete solution for vibecoders to vibecheck their code and save it locally even before it messess-up your production, vibecheck is a check point were they can automate their commit message to models like gpt-oss:20b, GPT4o-mini, Gemini-2.5-Flash, Claude-3.5-Haiku, Llama-3.3-70b (via Groq), Grok-beta, Kimi K2, Qwen-Turbo, DeepSeek-Chat, and Perplexity Sonar`,
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		diff, err := git.StagedDiff(cmd.Context())
		if err != nil {
			return fmt.Errorf("staged changes: %w", err)
		}
		if strings.TrimSpace(diff) == "" {
			notify.ShowStageReminder()
			return nil
		}

		additionalPrompt, err := cmd.Flags().GetString(promptFlagName)
		if err != nil {
			return fmt.Errorf("get string prompt flag: %w", err)
		}

		providerName, err := cmd.Flags().GetString(providerFlagName)
		if err != nil {
			return fmt.Errorf("get string provider flag: %w", err)
		}

		provider, err := llm.GetProvider(providerName)
		if err != nil {
			return err
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithColor("cyan"))

		s.Suffix = " Generating commit message..."
		s.Start()
		defer s.Stop()

		ctx, cancel := context.WithTimeout(cmd.Context(), time.Second*60)
		defer cancel()

		// Track latency
		startTime := time.Now()
		message, err := provider.GenerateCommitMessage(ctx, diff, additionalPrompt)
		latency := time.Since(startTime).Seconds()
		if err != nil {
			s.Stop()
			if envVar := detectMissingEnvVar(err); envVar != "" {
				notify.ShowMissingAPIKey(providerName, envVar)
				return nil
			}
			if model := detectMissingModel(err); model != "" {
				notify.ShowMissingModel(providerName, model)
				return nil
			}
			return fmt.Errorf("generated commit message: %w", err)
		}
		s.Stop()

		if err := git.CommitWMessage(cmd.Context(), message); err != nil {
			return fmt.Errorf("commit with message: %w", err)
		}

		// Record stats after successful commit
		// Extract first line of commit message for display
		commitMsg := message
		if idx := strings.Index(message, "\n"); idx > 0 {
			commitMsg = message[:idx]
		}
		if err := stats.RecordCommit(providerName, latency, commitMsg); err != nil {
			// Don't fail the commit if stats recording fails
			// Just log it silently
			_ = err
		}

		return nil
	},
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
	commitCmd.Flags().String(promptFlagName, "", "used to provide additional context to llm")
	commitCmd.Flags().String(providerFlagName, config.GetDefaultProvider(), fmt.Sprintf("used to select a particular ai-provider: %v (use 'vibecheck models' to change default)", strings.Join(llm.GetRegisteredNames(), ",")))
}

func detectMissingEnvVar(err error) string {
	const suffix = " environment variable not set"
	msg := err.Error()
	if !strings.HasSuffix(msg, suffix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimSuffix(msg, suffix))
}

func detectMissingModel(err error) string {
	msg := err.Error()
	if !strings.Contains(msg, "model '") || !strings.Contains(msg, "' not found") {
		return ""
	}

	start := strings.Index(msg, "model '")
	if start == -1 {
		return ""
	}
	start += len("model '")
	end := strings.Index(msg[start:], "'")
	if end == -1 {
		return ""
	}
	return msg[start : start+end]
}
