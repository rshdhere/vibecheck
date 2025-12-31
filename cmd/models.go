// Package cmd is provided by cobra-cli to ship command-line tools faster
/*
Copyright © 2025 raashed
*/
package cmd

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rshdhere/vibecheck/internal/config"
	"github.com/rshdhere/vibecheck/internal/llm"
	"github.com/spf13/cobra"
)

// Model represents a provider with its details
type Model struct {
	name        string
	displayName string
	model       string
	badge       string
	description string
}

func (m Model) Title() string       { return m.displayName }
func (m Model) Description() string { return m.description }
func (m Model) FilterValue() string { return m.name }

// Define all available models with their details
var availableModels = []Model{
	{
		name:        "openai",
		displayName: "OpenAI",
		model:       "gpt-4o-mini",
		badge:       "",
		description: "GPT-4o-mini • Fast • Reliable",
	},
	{
		name:        "gemini",
		displayName: "Google Gemini",
		model:       "gemini-2.5-flash",
		badge:       "",
		description: "gemini-2.5-flash • Ultra-Fast • 1M context",
	},
	{
		name:        "anthropic",
		displayName: "Anthropic Claude",
		model:       "claude-3.5-haiku",
		badge:       "",
		description: "claude-3.5-haiku • Fast • Best reasoning",
	},
	{
		name:        "groq",
		displayName: "Groq (Llama)",
		model:       "llama-3.3-70b-versatile",
		badge:       "",
		description: "llama-3.3-70b • Ultra • Free tier available",
	},
	{
		name:        "grok",
		displayName: "xAI Grok",
		model:       "grok-beta",
		badge:       "",
		description: "grok-beta • Fast • X's training data",
	},
	{
		name:        "kimi",
		displayName: "Moonshot Kimi",
		model:       "moonshot-v1-auto",
		badge:       "",
		description: "moonshot-v1-auto • Ultra-Fast • 128K context",
	},
	{
		name:        "qwen",
		displayName: "Alibaba Qwen",
		model:       "qwen-turbo",
		badge:       "",
		description: "qwen-turbo • Ultra-Fast • Multilingual",
	},
	{
		name:        "deepseek",
		displayName: "DeepSeek",
		model:       "deepseek-chat",
		badge:       "",
		description: "deepseek-chat • Ultra-Fast • Best value",
	},
	{
		name:        "perplexity",
		displayName: "Perplexity Sonar",
		model:       "sonar",
		badge:       "",
		description: "sonar • Fast • Search grounded",
	},
	{
		name:        "ollama",
		displayName: "Ollama (Local)",
		model:       "gpt-oss:20b",
		badge:       "Free",
		description: "gpt-oss:20b • Local • Private • No API key",
	},
}

type modelSelection struct {
	list         list.Model
	choice       string
	quitting     bool
	currentModel string
}

func (m modelSelection) Init() tea.Cmd {
	return nil
}

func (m modelSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			item, ok := m.list.SelectedItem().(Model)
			if ok {
				m.choice = item.name
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m modelSelection) View() string {
	if m.quitting {
		if m.choice == "" {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("Selection cancelled\n")
		}
		return ""
	}

	// Define color palette
	var (
		primaryColor   = lipgloss.Color("205")
		secondaryColor = lipgloss.Color("140")
		mutedColor     = lipgloss.Color("240")
		borderColor    = lipgloss.Color("238")
	)

	// Title section
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		MarginBottom(1).
		MarginTop(1)

	title := titleStyle.Render("VIBECHECK MODEL SELECTION")

	// Subtitle with current default
	subtitleStyle := lipgloss.NewStyle().
		Foreground(mutedColor)

	currentLabelStyle := lipgloss.NewStyle().
		Foreground(mutedColor)

	currentValueStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true)

	subtitle := subtitleStyle.Render(
		fmt.Sprintf("%s %s",
			currentLabelStyle.Render("Current:"),
			currentValueStyle.Render(m.currentModel),
		),
	)

	// Border for the list
	listBoxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	listBox := listBoxStyle.Render(m.list.View())

	// Help bar at the bottom
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

	helpContent := fmt.Sprintf("%s %s  %s %s  %s %s",
		helpKeyStyle.Render("↑/↓"),
		helpTextStyle.Render("navigate"),
		helpKeyStyle.Render("enter"),
		helpTextStyle.Render("select"),
		helpKeyStyle.Render("q"),
		helpTextStyle.Render("quit"),
	)

	help := helpStyle.Render(helpContent)

	return fmt.Sprintf("%s\n%s\n%s\n%s", title, subtitle, listBox, help)
}

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Interactively select the default AI provider for vibecheck",
	Long:  `Display all available AI providers and allow you to select a new default. The --provider flag will still work to override the default.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get registered providers from llm package
		registeredProviders := llm.GetRegisteredNames()

		// Filter models to only include registered ones
		var items []list.Item
		for _, model := range availableModels {
			// Check if this model is registered
			found := false
			for _, registered := range registeredProviders {
				if model.name == registered {
					found = true
					break
				}
			}
			if found {
				items = append(items, model)
			}
		}

		// Get current default
		currentDefault := config.GetDefaultProvider()

		// Create list
		const defaultWidth = 80
		const listHeight = 22

		l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
		l.Title = ""
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.SetShowHelp(false)
		l.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		l.Styles.HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

		// Find and set cursor to current default
		for i, item := range items {
			if model, ok := item.(Model); ok && model.name == currentDefault {
				l.Select(i)
				break
			}
		}

		m := modelSelection{
			list:         l,
			currentModel: currentDefault,
		}

		p := tea.NewProgram(m)
		finalModel, err := p.Run()
		if err != nil {
			return fmt.Errorf("error running program: %w", err)
		}

		if m, ok := finalModel.(modelSelection); ok {
			if m.choice != "" && m.choice != currentDefault {
				if err := config.SetDefaultProvider(m.choice); err != nil {
					return fmt.Errorf("failed to save configuration: %w", err)
				}

				successStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("140")).
					Bold(true)

				providerStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("205")).
					Bold(true)

				labelStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("252"))

				fmt.Printf("\n%s %s %s\n\n",
					successStyle.Render("SUCCESS"),
					labelStyle.Render("Default provider set to"),
					providerStyle.Render(m.choice))
			} else if m.choice == currentDefault {
				fmt.Println(lipgloss.NewStyle().
					Foreground(lipgloss.Color("240")).
					Render(fmt.Sprintf("\nNo changes made. Current default: %s\n", currentDefault)))
			}
		}

		return nil
	},
}

// Custom delegate for better styling
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(Model)
	if !ok {
		return
	}

	// Define colors
	var (
		primaryColor   = lipgloss.Color("205")
		secondaryColor = lipgloss.Color("140")
		mutedColor     = lipgloss.Color("245")
		normalColor    = lipgloss.Color("252")
		logoColor      = lipgloss.Color("213")
	)

	// VC Logo style
	logoStyle := lipgloss.NewStyle().
		Foreground(logoColor).
		Background(lipgloss.Color("235")).
		Bold(true).
		Padding(0, 1)

	// Styles
	providerStyle := lipgloss.NewStyle().
		Bold(true).
		Width(22)

	modelStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		Width(30)

	badgeStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true)

	// Selected styles
	selectedProviderStyle := providerStyle.Copy().
		Foreground(primaryColor)

	selectedModelStyle := modelStyle.Copy().
		Foreground(normalColor)

	// Cursor and selection
	var line1, line2 string

	if index == m.Index() {
		// Selected item with VC logo
		logo := logoStyle.Render("VC")

		badge := ""
		if item.badge != "" {
			badge = " " + badgeStyle.Render(item.badge)
		}

		line1 = fmt.Sprintf("%s %s%s",
			logo,
			selectedProviderStyle.Render(item.displayName),
			badge,
		)

		line2 = lipgloss.NewStyle().
			Foreground(normalColor).
			PaddingLeft(5).
			Render(selectedModelStyle.Render(item.model))
	} else {
		// Normal item
		badge := ""
		if item.badge != "" {
			badge = " " + lipgloss.NewStyle().
				Foreground(mutedColor).
				Render(item.badge)
		}

		line1 = fmt.Sprintf("    %s%s",
			providerStyle.Foreground(normalColor).Render(item.displayName),
			badge,
		)

		line2 = lipgloss.NewStyle().
			Foreground(mutedColor).
			PaddingLeft(5).
			Render(item.model)
	}

	fmt.Fprintf(w, "%s\n%s\n", line1, line2)
}

func init() {
	rootCmd.AddCommand(modelsCmd)
}
