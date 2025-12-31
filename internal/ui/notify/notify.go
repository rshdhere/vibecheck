package notify

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ShowStageReminder displays a minimal bubbletea notification letting the user
// know they need to stage files before running `vibecheck commit`.
func ShowStageReminder() {
	m := messageModel{
		title:       "NO STAGED CHANGES DETECTED !!",
		description: "Use `git add <files>` to stage your changes, then rerun `vibecheck commit`.",
		hint:        "Press any key (or wait a second) to continue.",
	}

	runProgram(m, "NO STAGED CHANGES DETECTED !!. Please stage files and rerun `vibecheck commit`.")
}

func ShowMissingAPIKey(providerName, envVar string) {
	title := fmt.Sprintf("%s API KEY REQUIRED !!", strings.ToUpper(providerName))
	description := fmt.Sprintf("Set %s in your environment or use `vibecheck keys` to store it globally.", envVar)
	hint := fmt.Sprintf("Run: vibecheck keys  OR  export %s=your_key_here", envVar)

	m := messageModel{
		title:       title,
		description: description,
		hint:        hint,
	}

	fallback := fmt.Sprintf("%s Please set %s (via `vibecheck keys` or export) and rerun `vibecheck commit`.", title, envVar)
	runProgram(m, fallback)
}

func ShowMissingModel(providerName, model string) {
	title := fmt.Sprintf("%s MODEL NOT AVAILABLE !!", strings.ToUpper(providerName))
	description := fmt.Sprintf("Model %q is missing locally. Pull or enable it before running `vibecheck commit`.", model)
	hint := fmt.Sprintf("Try: `ollama pull %s` or switch providers via `vibecheck models`.", model)

	m := messageModel{
		title:       title,
		description: description,
		hint:        hint,
	}

	fallback := fmt.Sprintf("%s Please install/pull %s or choose another provider.", title, model)
	runProgram(m, fallback)
}

func runProgram(m messageModel, fallback string) {
	p := tea.NewProgram(m, tea.WithoutSignalHandler())
	if _, err := p.Run(); err != nil {
		fmt.Println(fallback)
	}
}

type messageModel struct {
	width       int
	height      int
	title       string
	description string
	hint        string
}

type autoCloseMsg struct{}

func (m messageModel) Init() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg { return autoCloseMsg{} })
}

func (m messageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		return m, tea.Quit
	case autoCloseMsg:
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m messageModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Background(lipgloss.Color("235")).
		Bold(true).
		Padding(0, 1)

	bodyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244"))

	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		PaddingTop(0).
		PaddingBottom(1).
		PaddingLeft(2).
		PaddingRight(2)

	title := titleStyle.Render(m.title)
	logo := logoStyle.Render("VC")
	description := bodyStyle.Render(m.description)
	hint := hintStyle.Render(m.hint)

	contentWidth := lipgloss.Width(description)
	if w := lipgloss.Width(hint); w > contentWidth {
		contentWidth = w
	}
	headerMinWidth := lipgloss.Width(title)
	if lipgloss.Width(logo) > headerMinWidth {
		headerMinWidth = lipgloss.Width(logo)
	}
	if headerMinWidth > contentWidth {
		contentWidth = headerMinWidth
	}

	logoLine := lipgloss.NewStyle().
		Width(contentWidth).
		Render(strings.Repeat(" ", contentWidth-lipgloss.Width(logo)) + logo)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Width(contentWidth).Render(""),
		logoLine,
		lipgloss.NewStyle().Width(contentWidth).Render(""),
		lipgloss.NewStyle().Width(contentWidth).Render(title),
		lipgloss.NewStyle().Width(contentWidth).Render(description),
		lipgloss.NewStyle().Width(contentWidth).Render(hint),
	)

	box := boxStyle.Render(content)

	return "\n" + box + "\n"
}
