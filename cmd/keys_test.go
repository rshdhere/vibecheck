package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rshdhere/vibecheck/internal/keys"
)

func newTestKeysModel(items ...list.Item) keysModel {
	return keysModel{
		list:      list.New(items, keyItemDelegate{}, 80, 15),
		items:     items,
		textInput: textinput.New(),
		state:     "list",
	}
}

func update(t *testing.T, m keysModel, msg tea.Msg) (keysModel, tea.Cmd) {
	t.Helper()
	next, cmd := m.Update(msg)
	return next.(keysModel), cmd
}

func TestKeyItem(t *testing.T) {
	set := KeyItem{provider: "openai", displayName: "OpenAI", envVar: "OPENAI_API_KEY", hasKey: true, maskedKey: "sk-...1234"}
	unset := KeyItem{provider: "gemini", displayName: "Google Gemini", envVar: "GEMINI_API_KEY"}

	if got := set.Title(); got != "✓ OpenAI" {
		t.Errorf("Title() = %q, want %q", got, "✓ OpenAI")
	}
	if got := unset.Title(); got != "✗ Google Gemini" {
		t.Errorf("Title() = %q, want %q", got, "✗ Google Gemini")
	}
	if got := set.Description(); got != "sk-...1234 • OPENAI_API_KEY" {
		t.Errorf("Description() = %q", got)
	}
	if got := unset.Description(); got != "Not set • GEMINI_API_KEY" {
		t.Errorf("Description() = %q", got)
	}
	if got := set.FilterValue(); got != "openai" {
		t.Errorf("FilterValue() = %q, want openai", got)
	}
}

func TestLoadKeyItems(t *testing.T) {
	isolateHome(t)
	if err := keys.SetAPIKey("anthropic", "sk-ant-0123456789"); err != nil {
		t.Fatal(err)
	}

	items := loadKeyItems()

	if len(items) == 0 {
		t.Fatal("loadKeyItems() returned no items")
	}
	first := items[0].(KeyItem)
	if first.provider != "anthropic" || !first.hasKey || first.displayName != "Anthropic Claude" {
		t.Errorf("first item = %+v, want anthropic with a key listed first", first)
	}
	for _, item := range items {
		ki := item.(KeyItem)
		if ki.provider == "ollama" {
			t.Error("ollama should not be listed: it needs no API key")
		}
		if ki.provider != "anthropic" && ki.hasKey {
			t.Errorf("%s reported a key, want only anthropic", ki.provider)
		}
	}
}

func TestKeysModelSetKey(t *testing.T) {
	isolateHome(t)
	m := newTestKeysModel(KeyItem{provider: "openai", displayName: "OpenAI", envVar: "OPENAI_API_KEY"})

	m, _ = update(t, m, keyEnter)
	if m.state != "input" || m.selectedItem.provider != "openai" {
		t.Fatalf("after enter: state = %q, selected = %q; want input for openai", m.state, m.selectedItem.provider)
	}
	if view := m.View(); !strings.Contains(view, "Set API Key for") || !strings.Contains(view, "OPENAI_API_KEY") {
		t.Errorf("input view missing prompt: %q", view)
	}

	m, _ = update(t, m, keyRunes("sk-test-key"))
	m, cmd := update(t, m, keyEnter)
	if m.state != "list" || cmd == nil {
		t.Fatalf("after save: state = %q, cmd = %v; want list and a reload", m.state, cmd)
	}
	if got, _ := keys.GetAPIKey("openai"); got != "sk-test-key" {
		t.Errorf("stored key = %q, want sk-test-key", got)
	}

	m, _ = update(t, m, cmd())
	for _, item := range m.items {
		if ki := item.(KeyItem); ki.provider == "openai" && !ki.hasKey {
			t.Error("reloaded openai item should report a key")
		}
	}
}

func TestKeysModelDeleteKey(t *testing.T) {
	t.Run("d on a stored key", func(t *testing.T) {
		isolateHome(t)
		if err := keys.SetAPIKey("openai", "sk-old"); err != nil {
			t.Fatal(err)
		}
		m := newTestKeysModel(KeyItem{provider: "openai", hasKey: true})

		_, cmd := update(t, m, keyRunes("d"))
		if cmd == nil {
			t.Error("delete should trigger a reload")
		}
		if _, ok := keys.GetAPIKey("openai"); ok {
			t.Error("key still stored after delete")
		}
	})

	t.Run("empty input", func(t *testing.T) {
		isolateHome(t)
		if err := keys.SetAPIKey("openai", "sk-old"); err != nil {
			t.Fatal(err)
		}
		m := newTestKeysModel(KeyItem{provider: "openai", hasKey: true})

		m, _ = update(t, m, keyEnter)
		m, cmd := update(t, m, keyEnter)
		if m.state != "list" || cmd == nil {
			t.Errorf("state = %q, cmd = %v; want list and a reload", m.state, cmd)
		}
		if _, ok := keys.GetAPIKey("openai"); ok {
			t.Error("key still stored after submitting empty input")
		}
	})
}

func TestKeysModelErrors(t *testing.T) {
	isolateHome(t)
	bogus := KeyItem{provider: "bogus", displayName: "Bogus", hasKey: true}

	t.Run("delete", func(t *testing.T) {
		m, _ := update(t, newTestKeysModel(bogus), keyRunes("d"))
		if !strings.Contains(m.errorMsg, "Error deleting key") {
			t.Errorf("errorMsg = %q, want a delete error", m.errorMsg)
		}
		if !strings.Contains(m.View(), "ERROR: Error deleting key") {
			t.Error("view should show the error")
		}
	})

	t.Run("save", func(t *testing.T) {
		m, _ := update(t, newTestKeysModel(bogus), keyEnter)
		m, _ = update(t, m, keyRunes("x"))
		m, _ = update(t, m, keyEnter)
		if !strings.Contains(m.errorMsg, "Error saving key") {
			t.Errorf("errorMsg = %q, want a save error", m.errorMsg)
		}
	})

	t.Run("clear", func(t *testing.T) {
		m, _ := update(t, newTestKeysModel(bogus), keyEnter)
		m, _ = update(t, m, keyEnter)
		if !strings.Contains(m.errorMsg, "Error deleting key") {
			t.Errorf("errorMsg = %q, want a delete error", m.errorMsg)
		}
	})
}

func TestKeysModelNavigation(t *testing.T) {
	m := newTestKeysModel(KeyItem{provider: "openai"})

	m, _ = update(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	if m.list.Width() != 120 {
		t.Errorf("list width = %d, want 120", m.list.Width())
	}

	m, _ = update(t, m, keyEnter)
	m, _ = update(t, m, keyEsc)
	if m.state != "list" {
		t.Errorf("esc from input: state = %q, want list", m.state)
	}
	if view := m.View(); !strings.Contains(view, "VIBECHECK API KEYS") || !strings.Contains(view, "navigate") {
		t.Errorf("list view missing title or help: %q", view)
	}

	m, _ = update(t, m, keyRunes("j"))
	m, cmd := update(t, m, keyRunes("q"))
	if !m.quitting || cmd == nil {
		t.Error("q should quit")
	}
	if m.View() != "" {
		t.Error("view should be empty after quitting")
	}
}

func TestKeyItemDelegateRender(t *testing.T) {
	items := []list.Item{
		KeyItem{provider: "openai", displayName: "OpenAI", envVar: "OPENAI_API_KEY", hasKey: true, maskedKey: "sk-...1234"},
		KeyItem{provider: "gemini", displayName: "Google Gemini", envVar: "GEMINI_API_KEY"},
	}
	l := list.New(items, keyItemDelegate{}, 80, 15)
	d := keyItemDelegate{}

	if d.Height() != 2 || d.Spacing() != 0 || d.Update(nil, &l) != nil {
		t.Error("unexpected delegate layout")
	}
	for i := range items {
		var buf bytes.Buffer
		for _, selected := range []int{0, 1} {
			l.Select(selected)
			buf.Reset()
			d.Render(&buf, l, i, items[i])
			if !strings.Contains(buf.String(), items[i].(KeyItem).displayName) {
				t.Errorf("render of item %d (selected %d) = %q, missing display name", i, selected, buf.String())
			}
		}
	}

	var buf bytes.Buffer
	d.Render(&buf, l, 0, Model{name: "not-a-key-item"})
	if buf.Len() != 0 {
		t.Error("render should ignore non-KeyItem items")
	}
}

func TestKeysCmdSavesKey(t *testing.T) {
	isolateHome(t)
	scriptTUI(t, "\rsk-from-tui\rq")

	if err := execute(t, "keys"); err != nil {
		t.Fatalf("keys error = %v", err)
	}

	stored, err := keys.GetAllKeys()
	if err != nil || len(stored) != 1 {
		t.Fatalf("stored keys = %v, %v; want exactly one", stored, err)
	}
}
