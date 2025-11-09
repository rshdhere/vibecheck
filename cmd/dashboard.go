// Package cmd is provided by cobra-cli to ship command-line tools faster
/*
Copyright © 2025 raashed
*/
package cmd

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rshdhere/vibecheck/internal/stats"
	"github.com/spf13/cobra"
)

type dashboardModel struct {
	totalCommits  int
	mostUsedModel string
	avgLatency    float64
	lastUsed      time.Time
	recentCommits []stats.CommitRecord
	width         int
	height        int
	quitting      bool
}

type tickMsg struct{}

func (m dashboardModel) Init() tea.Cmd {
	return tea.Batch(
		loadStats(),
		tick(),
	)
}

func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "r":
			// Refresh on 'r' key
			return m, loadStats()
		}

	case tickMsg:
		return m, tea.Batch(
			loadStats(),
			tick(),
		)

	case statsLoadedMsg:
		m.totalCommits = msg.totalCommits
		m.mostUsedModel = msg.mostUsedModel
		m.avgLatency = msg.avgLatency
		m.lastUsed = msg.lastUsed
		m.recentCommits = msg.recentCommits
		return m, nil
	}

	return m, nil
}

func (m dashboardModel) View() string {
	if m.quitting {
		return ""
	}

	// Define color palette
	var (
		primaryColor   = lipgloss.Color("205")
		secondaryColor = lipgloss.Color("140")
		mutedColor     = lipgloss.Color("240")
		borderColor    = lipgloss.Color("238")
		successColor   = lipgloss.Color("76")
	)

	// Title section
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		MarginBottom(1).
		MarginTop(1)

	title := titleStyle.Render("Vibecheck Dashboard")

	// Stats section
	statsBoxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(mutedColor)

	valueStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true)

	// Format last used time
	lastUsedStr := "Never"
	if !m.lastUsed.IsZero() {
		now := time.Now()
		diff := now.Sub(m.lastUsed)
		if diff < time.Minute {
			lastUsedStr = "Just now"
		} else if diff < time.Hour {
			minutes := int(diff.Minutes())
			lastUsedStr = fmt.Sprintf("%dm ago", minutes)
		} else if diff < 24*time.Hour {
			hours := int(diff.Hours())
			lastUsedStr = fmt.Sprintf("%dh ago", hours)
		} else {
			days := int(diff.Hours() / 24)
			lastUsedStr = fmt.Sprintf("%dd ago", days)
		}
	}

	// Format average latency
	avgLatencyStr := "0.0s"
	if m.avgLatency > 0 {
		avgLatencyStr = fmt.Sprintf("%.1fs", m.avgLatency)
	}

	// Format most used model
	modelDisplay := m.mostUsedModel
	if modelDisplay == "N/A" {
		modelDisplay = "None"
	}

	statsContent := fmt.Sprintf(
		"%s %s\n%s %s\n%s %s\n%s %s",
		labelStyle.Render("Total commits AI-generated:"),
		valueStyle.Render(fmt.Sprintf("%d", m.totalCommits)),
		labelStyle.Render("Most used model:"),
		valueStyle.Render(modelDisplay),
		labelStyle.Render("Average latency:"),
		valueStyle.Render(avgLatencyStr),
		labelStyle.Render("Last used:"),
		valueStyle.Render(lastUsedStr),
	)

	statsBox := statsBoxStyle.Render(statsContent)

	// Recent commits section
	commitsBoxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		MarginBottom(1)

	commitsTitleStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true).
		MarginBottom(1)

	commitItemStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		MarginLeft(2)

	commitCheckStyle := lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true)

	commitsContent := commitsTitleStyle.Render("Recent commits:")
	if len(m.recentCommits) == 0 {
		commitsContent += "\n" + commitItemStyle.Render("No commits yet")
	} else {
		for _, commit := range m.recentCommits {
			commitsContent += "\n" + fmt.Sprintf(
				"%s %s",
				commitCheckStyle.Render("✓"),
				commitItemStyle.Render(commit.CommitMsg),
			)
		}
	}

	commitsBox := commitsBoxStyle.Render(commitsContent)

	// Help bar
	helpStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderForeground(borderColor).
		PaddingTop(1)

	helpKeyStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true)

	helpTextStyle := lipgloss.NewStyle().
		Foreground(mutedColor)

	helpContent := fmt.Sprintf("%s %s  %s %s  %s",
		helpKeyStyle.Render("r"),
		helpTextStyle.Render("refresh"),
		helpKeyStyle.Render("q"),
		helpTextStyle.Render("quit"),
		helpTextStyle.Render("(auto-refreshes every 5s)"),
	)

	help := helpStyle.Render(helpContent)

	// Combine all sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		statsBox,
		commitsBox,
		help,
	)

	// Center content if window is large enough
	if m.width > 0 {
		content = lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			content,
		)
	}

	return content
}

type statsLoadedMsg struct {
	totalCommits  int
	mostUsedModel string
	avgLatency    float64
	lastUsed      time.Time
	recentCommits []stats.CommitRecord
}

func loadStats() tea.Cmd {
	return func() tea.Msg {
		totalCommits, _ := stats.GetTotalCommits()
		mostUsedModel, _ := stats.GetMostUsedModel()
		avgLatency, _ := stats.GetAverageLatency()
		lastUsed, _ := stats.GetLastUsed()
		recentCommits, _ := stats.GetRecentCommits(10)

		return statsLoadedMsg{
			totalCommits:  totalCommits,
			mostUsedModel: mostUsedModel,
			avgLatency:    avgLatency,
			lastUsed:      lastUsed,
			recentCommits: recentCommits,
		}
	}
}

func tick() tea.Cmd {
	return tea.Tick(5*time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Display statistics dashboard for AI-generated commits",
	Long:  `Display a live dashboard showing statistics about your AI-generated commits, including total commits, most used model, average latency, and recent commits.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		m := dashboardModel{
			width:  80,
			height: 24,
		}

		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running dashboard: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
