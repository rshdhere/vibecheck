package cmd

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rshdhere/vibecheck/internal/stats"
)

func TestDashboardModelUpdate(t *testing.T) {
	var m tea.Model = dashboardModel{}

	if m.(dashboardModel).Init() == nil {
		t.Error("Init() should load stats and start the refresh tick")
	}

	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	if d := m.(dashboardModel); d.width != 100 || d.height != 40 {
		t.Errorf("size = %dx%d, want 100x40", d.width, d.height)
	}

	loaded := statsLoadedMsg{totalCommits: 3, mostUsedModel: "groq", avgLatency: 1.5, lastUsed: time.Now()}
	m, _ = m.Update(loaded)
	if d := m.(dashboardModel); d.totalCommits != 3 || d.mostUsedModel != "groq" || d.avgLatency != 1.5 {
		t.Errorf("stats not applied: %+v", d)
	}

	if _, cmd := m.Update(tickMsg{}); cmd == nil {
		t.Error("tick should reload stats and schedule the next tick")
	}
	if _, cmd := m.Update(keyRunes("r")); cmd == nil {
		t.Error("r should reload stats")
	}
	if _, cmd := m.Update(keyRunes("x")); cmd != nil {
		t.Error("unbound keys should do nothing")
	}

	m, cmd := m.Update(keyRunes("q"))
	if !m.(dashboardModel).quitting || cmd == nil {
		t.Error("q should quit")
	}
	if m.View() != "" {
		t.Error("view should be empty after quitting")
	}
}

func TestDashboardModelView(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		view := dashboardModel{mostUsedModel: "N/A"}.View()
		for _, want := range []string{"Vibecheck Dashboard", "Never", "None", "0.0s", "No commits yet"} {
			if !strings.Contains(view, want) {
				t.Errorf("view missing %q", want)
			}
		}
	})

	lastUsed := []struct {
		ago  time.Duration
		want string
	}{
		{10 * time.Second, "Just now"},
		{5 * time.Minute, "5m ago"},
		{3 * time.Hour, "3h ago"},
		{50 * time.Hour, "2d ago"},
	}
	for _, tt := range lastUsed {
		t.Run(tt.want, func(t *testing.T) {
			m := dashboardModel{
				width:         120,
				height:        40,
				totalCommits:  2,
				mostUsedModel: "openai",
				avgLatency:    2.25,
				lastUsed:      time.Now().Add(-tt.ago),
				recentCommits: []stats.CommitRecord{{Model: "openai", Latency: 2.25, CommitMsg: "feat: add dashboard"}},
			}
			view := m.View()
			for _, want := range []string{tt.want, "openai", "2.2s", "feat: add dashboard"} {
				if !strings.Contains(view, want) {
					t.Errorf("view missing %q", want)
				}
			}
		})
	}
}

func TestLoadStats(t *testing.T) {
	isolateHome(t)
	if err := stats.RecordCommit("groq", 1.0, "feat: one"); err != nil {
		t.Fatal(err)
	}

	msg, ok := loadStats()().(statsLoadedMsg)
	if !ok {
		t.Fatal("loadStats() did not return statsLoadedMsg")
	}
	if msg.totalCommits != 1 || msg.mostUsedModel != "groq" || len(msg.recentCommits) != 1 {
		t.Errorf("loaded %+v, want the recorded commit", msg)
	}
	if tick() == nil {
		t.Error("tick() returned nil")
	}
}

func TestDashboardCmd(t *testing.T) {
	isolateHome(t)
	scriptTUI(t, "q")

	if err := execute(t, "dashboard"); err != nil {
		t.Fatalf("dashboard error = %v", err)
	}
}
