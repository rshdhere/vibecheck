// Package cmd is provided by cobra-cli to ship command-line tools faster
/*
Copyright © 2025 raashed
*/
package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rshdhere/vibecheck/internal/keys"
	"github.com/rshdhere/vibecheck/internal/llm"
	"github.com/spf13/cobra"
)

// KeyItem represents a provider key entry in the list
type KeyItem struct {
	provider    string
	displayName string
	envVar      string
	hasKey      bool
	maskedKey   string
}

func (k KeyItem) Title() string {
	status := "✗"
	if k.hasKey {
		status = "✓"
	}
	return fmt.Sprintf("%s %s", status, k.displayName)
}

func (k KeyItem) Description() string {
	if k.hasKey {
		return fmt.Sprintf("%s • %s", k.maskedKey, k.envVar)
	}
	return fmt.Sprintf("Not set • %s", k.envVar)
}

func (k KeyItem) FilterValue() string {
	return k.provider
}

type keysModel struct {
	list         list.Model
	items        []list.Item
	selectedItem KeyItem
	textInput    textinput.Model
	state        string // "list", "input", "saving"
	quitting     bool
	errorMsg     string
}

func (m keysModel) Init() tea.Cmd {
	return nil
}

func (m keysModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case itemsReloadedMsg:
		// Reload items after save/delete
		items := loadKeyItems()
		m.items = items
		m.list.SetItems(items)
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case "list":
			switch keypress := msg.String(); keypress {
			case "ctrl+c", "q", "esc":
				m.quitting = true
				return m, tea.Quit

			case "enter":
				item, ok := m.list.SelectedItem().(KeyItem)
				if ok {
					m.selectedItem = item
					m.state = "input"
					m.textInput.SetValue("")
					m.textInput.Focus()
					m.errorMsg = ""
					return m, textinput.Blink
				}

			case "d":
				// Delete key
				item, ok := m.list.SelectedItem().(KeyItem)
				if ok && item.hasKey {
					if err := keys.SetAPIKey(item.provider, ""); err != nil {
						m.errorMsg = fmt.Sprintf("Error deleting key: %v", err)
					} else {
						// Reload items
						return m, m.reloadItems()
					}
				}
			}

		case "input":
			switch keypress := msg.String(); keypress {
			case "ctrl+c", "esc":
				m.state = "list"
				m.textInput.Blur()
				return m, nil

			case "enter":
				key := strings.TrimSpace(m.textInput.Value())
				if key == "" {
					// Delete key if empty
					if err := keys.SetAPIKey(m.selectedItem.provider, ""); err != nil {
						m.errorMsg = fmt.Sprintf("Error deleting key: %v", err)
					} else {
						m.state = "list"
						m.textInput.Blur()
						return m, m.reloadItems()
					}
				} else {
					// Save key
					if err := keys.SetAPIKey(m.selectedItem.provider, key); err != nil {
						m.errorMsg = fmt.Sprintf("Error saving key: %v", err)
					} else {
						m.state = "list"
						m.textInput.Blur()
						m.errorMsg = ""
						return m, m.reloadItems()
					}
				}
			}
		}
	}

	// Update focused component
	var cmd tea.Cmd
	if m.state == "input" {
		m.textInput, cmd = m.textInput.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m keysModel) View() string {
	if m.quitting {
		return ""
	}

	// Define color palette
	var (
		primaryColor   = lipgloss.Color("205")
		secondaryColor = lipgloss.Color("140")
		mutedColor     = lipgloss.Color("240")
		borderColor    = lipgloss.Color("238")
		errorColor     = lipgloss.Color("196")
	)

	// Title section
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		MarginBottom(1).
		MarginTop(1)

	title := titleStyle.Render("VIBECHECK API KEYS")

	// Error message
	errorStyle := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		MarginBottom(1)

	errorView := ""
	if m.errorMsg != "" {
		errorView = errorStyle.Render("ERROR: "+m.errorMsg) + "\n"
	}

	// Main content based on state
	var content string
	if m.state == "input" {
		// Input view
		inputBoxStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

		labelStyle := lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			MarginBottom(1)

		hintStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

		providerStyle := lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

		envVarStyle := lipgloss.NewStyle().
			Foreground(mutedColor)

		inputContent := fmt.Sprintf(
			"%s\n%s\n\n%s\n\n%s",
			labelStyle.Render(fmt.Sprintf("Set API Key for %s", providerStyle.Render(m.selectedItem.displayName))),
			envVarStyle.Render(fmt.Sprintf("Environment variable: %s", m.selectedItem.envVar)),
			m.textInput.View(),
			hintStyle.Render("Press Enter to save, Esc to cancel (leave empty to delete)"),
		)

		content = inputBoxStyle.Render(inputContent)
	} else {
		// List view
		listBoxStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

		listBox := listBoxStyle.Render(m.list.View())
		content = listBox
	}

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

	var helpContent string
	if m.state == "input" {
		helpContent = fmt.Sprintf("%s %s  %s %s",
			helpKeyStyle.Render("enter"),
			helpTextStyle.Render("save"),
			helpKeyStyle.Render("esc"),
			helpTextStyle.Render("cancel"),
		)
	} else {
		helpContent = fmt.Sprintf("%s %s  %s %s  %s %s  %s %s",
			helpKeyStyle.Render("↑/↓"),
			helpTextStyle.Render("navigate"),
			helpKeyStyle.Render("enter"),
			helpTextStyle.Render("set key"),
			helpKeyStyle.Render("d"),
			helpTextStyle.Render("delete"),
			helpKeyStyle.Render("q"),
			helpTextStyle.Render("quit"),
		)
	}

	help := helpStyle.Render(helpContent)

	return fmt.Sprintf("%s\n%s%s\n%s", title, errorView, content, help)
}

func (m keysModel) reloadItems() tea.Cmd {
	return func() tea.Msg {
		return itemsReloadedMsg{}
	}
}

type itemsReloadedMsg struct{}

func loadKeyItems() []list.Item {
	registeredProviders := llm.GetRegisteredNames()
	allKeys, _ := keys.GetAllKeys()

	var itemsWithKey []list.Item
	var itemsWithoutKey []list.Item

	for _, provider := range registeredProviders {
		// Skip ollama as it doesn't need an API key
		if provider == "ollama" {
			continue
		}

		envVar, ok := keys.ProviderToEnvVar[provider]
		if !ok {
			continue
		}

		// Get display name from models.go
		displayName := provider
		for _, model := range availableModels {
			if model.name == provider {
				displayName = model.displayName
				break
			}
		}

		maskedKey, hasKey := allKeys[provider]

		item := KeyItem{
			provider:    provider,
			displayName: displayName,
			envVar:      envVar,
			hasKey:      hasKey,
			maskedKey:   maskedKey,
		}

		if hasKey {
			itemsWithKey = append(itemsWithKey, item)
		} else {
			itemsWithoutKey = append(itemsWithoutKey, item)
		}
	}

	// Return items with keys first, then items without keys
	return append(itemsWithKey, itemsWithoutKey...)
}

// Custom delegate for better styling
type keyItemDelegate struct{}

func (d keyItemDelegate) Height() int                             { return 2 }
func (d keyItemDelegate) Spacing() int                            { return 0 }
func (d keyItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d keyItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(KeyItem)
	if !ok {
		return
	}

	// Define colors
	var (
		primaryColor = lipgloss.Color("205")
		mutedColor   = lipgloss.Color("245")
		normalColor  = lipgloss.Color("252")
		successColor = lipgloss.Color("76")
		logoColor    = lipgloss.Color("213")
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

	descStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		Width(50)

	// Selected styles
	selectedProviderStyle := providerStyle.Copy().
		Foreground(primaryColor)

	selectedDescStyle := descStyle.Copy().
		Foreground(normalColor)

	// Cursor and selection
	var line1, line2 string

	if index == m.Index() {
		// Selected item with VC logo
		logo := logoStyle.Render("VC")

		statusIcon := lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("✗")
		if item.hasKey {
			statusIcon = lipgloss.NewStyle().
				Foreground(successColor).
				Render("✓")
		}

		line1 = fmt.Sprintf("%s %s %s",
			logo,
			statusIcon,
			selectedProviderStyle.Render(item.displayName),
		)

		desc := item.Description()
		line2 = lipgloss.NewStyle().
			Foreground(normalColor).
			PaddingLeft(5).
			Render(selectedDescStyle.Render(desc))
	} else {
		// Normal item
		statusIcon := lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("✗")
		if item.hasKey {
			statusIcon = lipgloss.NewStyle().
				Foreground(successColor).
				Render("✓")
		}

		line1 = fmt.Sprintf("    %s %s",
			statusIcon,
			providerStyle.Foreground(normalColor).Render(item.displayName),
		)

		line2 = lipgloss.NewStyle().
			Foreground(mutedColor).
			PaddingLeft(5).
			Render(item.Description())
	}

	fmt.Fprintf(w, "%s\n%s\n", line1, line2)
}

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Interactively manage API keys for AI providers",
	Long:  `Manage API keys for all AI providers. Keys are stored globally in ~/.vibecheck_keys.json and will be used when environment variables are not set.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		items := loadKeyItems()

		// Create list
		const defaultWidth = 80
		const listHeight = 15

		l := list.New(items, keyItemDelegate{}, defaultWidth, listHeight)
		l.Title = ""
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.SetShowHelp(false)
		l.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		l.Styles.HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

		// Create text input
		ti := textinput.New()
		ti.Placeholder = "Enter API key..."
		ti.CharLimit = 200
		ti.Width = 60
		ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

		m := keysModel{
			list:      l,
			items:     items,
			textInput: ti,
			state:     "list",
		}

		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running program: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(keysCmd)
}
