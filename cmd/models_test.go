package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rshdhere/vibecheck/internal/config"
)

func modelItems() []list.Item {
	items := make([]list.Item, len(availableModels))
	for i, m := range availableModels {
		items[i] = m
	}
	return items
}

func TestModelItem(t *testing.T) {
	m := availableModels[0]
	if m.Title() != m.displayName || m.Description() != m.description || m.FilterValue() != m.name {
		t.Errorf("list item accessors do not match %+v", m)
	}
}

func TestModelSelectionUpdate(t *testing.T) {
	m := modelSelection{list: list.New(modelItems(), itemDelegate{}, 80, 22), currentModel: "openai"}

	next, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = next.(modelSelection)
	if m.list.Width() != 100 {
		t.Errorf("list width = %d, want 100", m.list.Width())
	}

	next, _ = m.Update(keyRunes("j"))
	m = next.(modelSelection)
	if view := m.View(); !strings.Contains(view, "VIBECHECK MODEL SELECTION") || !strings.Contains(view, "Current:") {
		t.Errorf("view missing title or current model: %q", view)
	}

	next, cmd := m.Update(keyEnter)
	m = next.(modelSelection)
	if m.choice != availableModels[1].name || cmd == nil {
		t.Errorf("choice = %q, want %q and a quit", m.choice, availableModels[1].name)
	}
}

func TestModelSelectionQuit(t *testing.T) {
	m := modelSelection{list: list.New(modelItems(), itemDelegate{}, 80, 22)}

	next, cmd := m.Update(keyRunes("q"))
	m = next.(modelSelection)
	if !m.quitting || cmd == nil {
		t.Fatal("q should quit")
	}
	if !strings.Contains(m.View(), "Selection cancelled") {
		t.Errorf("view = %q, want cancellation notice", m.View())
	}

	m.choice = "openai"
	if m.View() != "" {
		t.Error("view should be empty after quitting with a choice")
	}
}

func TestItemDelegateRender(t *testing.T) {
	items := modelItems()
	l := list.New(items, itemDelegate{}, 80, 22)
	d := itemDelegate{}

	if d.Height() != 2 || d.Spacing() != 0 || d.Update(nil, &l) != nil {
		t.Error("unexpected delegate layout")
	}
	for i, item := range items {
		for _, selected := range []int{0, i} {
			l.Select(selected)
			var buf bytes.Buffer
			d.Render(&buf, l, i, item)
			m := item.(Model)
			if !strings.Contains(buf.String(), m.displayName) || !strings.Contains(buf.String(), m.model) {
				t.Errorf("render of %s = %q, missing name or model", m.name, buf.String())
			}
			if m.badge != "" && !strings.Contains(buf.String(), m.badge) {
				t.Errorf("render of %s missing badge %q", m.name, m.badge)
			}
		}
	}

	var buf bytes.Buffer
	d.Render(&buf, l, 0, KeyItem{})
	if buf.Len() != 0 {
		t.Error("render should ignore non-Model items")
	}
}

func TestModelsCmd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"select next provider", "j\r", availableModels[1].name},
		{"keep current provider", "\r", "openai"},
		{"cancel", "q", "openai"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolateHome(t)
			scriptTUI(t, tt.input)

			if err := execute(t, "models"); err != nil {
				t.Fatalf("models error = %v", err)
			}
			if got := config.GetDefaultProvider(); got != tt.want {
				t.Errorf("default provider = %q, want %q", got, tt.want)
			}
		})
	}
}
