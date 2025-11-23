package stats

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetStatsPath(t *testing.T) {
	path, err := getStatsPath()
	if err != nil {
		t.Fatalf("getStatsPath() error = %v", err)
	}
	if path == "" {
		t.Error("getStatsPath() returned empty path")
	}
	if filepath.Base(path) != ".vibecheck_stats.json" {
		t.Errorf("getStatsPath() = %v, want path ending with .vibecheck_stats.json", path)
	}
}

func TestLoad(t *testing.T) {
	t.Run("non-existent stats returns empty", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldHome := os.Getenv("HOME")
		defer os.Setenv("HOME", oldHome)
		os.Setenv("HOME", tmpDir)

		stats, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if stats == nil {
			t.Error("Load() returned nil")
		}
		if len(stats.Commits) != 0 {
			t.Errorf("Load() Commits length = %v, want 0", len(stats.Commits))
		}
	})

	t.Run("existing stats loads correctly", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldHome := os.Getenv("HOME")
		defer os.Setenv("HOME", oldHome)
		os.Setenv("HOME", tmpDir)

		statsPath := filepath.Join(tmpDir, ".vibecheck_stats.json")
		statsData := `{"commits": [{"timestamp": "2025-01-01T00:00:00Z", "model": "openai", "latency": 1.5, "commit_msg": "test commit"}]}`
		if err := os.WriteFile(statsPath, []byte(statsData), 0644); err != nil {
			t.Fatalf("Failed to write stats: %v", err)
		}

		stats, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if len(stats.Commits) != 1 {
			t.Errorf("Load() Commits length = %v, want 1", len(stats.Commits))
		}
	})
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	stats := &Stats{
		Commits: []CommitRecord{
			{
				Timestamp: time.Now(),
				Model:     "openai",
				Latency:   1.5,
				CommitMsg: "test commit",
			},
		},
	}
	if err := Save(stats); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify it was saved
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() after Save() error = %v", err)
	}
	if len(loaded.Commits) != 1 {
		t.Errorf("Load() after Save() Commits length = %v, want 1", len(loaded.Commits))
	}
	if loaded.Commits[0].Model != "openai" {
		t.Errorf("Load() after Save() Model = %v, want openai", loaded.Commits[0].Model)
	}
}

func TestRecordCommit(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	if err := RecordCommit("gemini", 2.3, "feat: add new feature"); err != nil {
		t.Fatalf("RecordCommit() error = %v", err)
	}

	stats, err := Load()
	if err != nil {
		t.Fatalf("Load() after RecordCommit() error = %v", err)
	}
	if len(stats.Commits) != 1 {
		t.Fatalf("Load() after RecordCommit() Commits length = %v, want 1", len(stats.Commits))
	}

	commit := stats.Commits[0]
	if commit.Model != "gemini" {
		t.Errorf("RecordCommit() Model = %v, want gemini", commit.Model)
	}
	if commit.Latency != 2.3 {
		t.Errorf("RecordCommit() Latency = %v, want 2.3", commit.Latency)
	}
	if commit.CommitMsg != "feat: add new feature" {
		t.Errorf("RecordCommit() CommitMsg = %v, want feat: add new feature", commit.CommitMsg)
	}
	if commit.Timestamp.IsZero() {
		t.Error("RecordCommit() Timestamp is zero")
	}
}

func TestGetTotalCommits(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Initially should be 0
	total, err := GetTotalCommits()
	if err != nil {
		t.Fatalf("GetTotalCommits() error = %v", err)
	}
	if total != 0 {
		t.Errorf("GetTotalCommits() = %v, want 0", total)
	}

	// Add some commits
	RecordCommit("openai", 1.0, "commit 1")
	RecordCommit("gemini", 2.0, "commit 2")
	RecordCommit("anthropic", 3.0, "commit 3")

	total, err = GetTotalCommits()
	if err != nil {
		t.Fatalf("GetTotalCommits() error = %v", err)
	}
	if total != 3 {
		t.Errorf("GetTotalCommits() = %v, want 3", total)
	}
}

func TestGetMostUsedModel(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Empty stats
	model, err := GetMostUsedModel()
	if err != nil {
		t.Fatalf("GetMostUsedModel() error = %v", err)
	}
	if model != "N/A" {
		t.Errorf("GetMostUsedModel() = %v, want N/A", model)
	}

	// Add commits
	RecordCommit("openai", 1.0, "commit 1")
	RecordCommit("openai", 1.5, "commit 2")
	RecordCommit("gemini", 2.0, "commit 3")
	RecordCommit("openai", 1.2, "commit 4")

	model, err = GetMostUsedModel()
	if err != nil {
		t.Fatalf("GetMostUsedModel() error = %v", err)
	}
	if model != "openai" {
		t.Errorf("GetMostUsedModel() = %v, want openai", model)
	}
}

func TestGetAverageLatency(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Empty stats
	avg, err := GetAverageLatency()
	if err != nil {
		t.Fatalf("GetAverageLatency() error = %v", err)
	}
	if avg != 0 {
		t.Errorf("GetAverageLatency() = %v, want 0", avg)
	}

	// Add commits
	RecordCommit("openai", 1.0, "commit 1")
	RecordCommit("gemini", 2.0, "commit 2")
	RecordCommit("anthropic", 3.0, "commit 3")

	avg, err = GetAverageLatency()
	if err != nil {
		t.Fatalf("GetAverageLatency() error = %v", err)
	}
	expected := 2.0 // (1.0 + 2.0 + 3.0) / 3
	if avg != expected {
		t.Errorf("GetAverageLatency() = %v, want %v", avg, expected)
	}
}

func TestGetLastUsed(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Empty stats
	lastUsed, err := GetLastUsed()
	if err != nil {
		t.Fatalf("GetLastUsed() error = %v", err)
	}
	if !lastUsed.IsZero() {
		t.Errorf("GetLastUsed() = %v, want zero time", lastUsed)
	}

	// Add a commit
	before := time.Now()
	RecordCommit("openai", 1.0, "commit 1")
	after := time.Now()

	lastUsed, err = GetLastUsed()
	if err != nil {
		t.Fatalf("GetLastUsed() error = %v", err)
	}
	if lastUsed.IsZero() {
		t.Error("GetLastUsed() is zero after RecordCommit")
	}
	if lastUsed.Before(before) || lastUsed.After(after) {
		t.Errorf("GetLastUsed() = %v, want between %v and %v", lastUsed, before, after)
	}
}

func TestGetRecentCommits(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Empty stats
	commits, err := GetRecentCommits(10)
	if err != nil {
		t.Fatalf("GetRecentCommits() error = %v", err)
	}
	if len(commits) != 0 {
		t.Errorf("GetRecentCommits() length = %v, want 0", len(commits))
	}

	// Add multiple commits with small delays to ensure different timestamps
	RecordCommit("openai", 1.0, "commit 1")
	time.Sleep(10 * time.Millisecond)
	RecordCommit("gemini", 2.0, "commit 2")
	time.Sleep(10 * time.Millisecond)
	RecordCommit("anthropic", 3.0, "commit 3")
	time.Sleep(10 * time.Millisecond)
	RecordCommit("groq", 4.0, "commit 4")
	time.Sleep(10 * time.Millisecond)
	RecordCommit("grok", 5.0, "commit 5")

	// Get recent commits
	commits, err = GetRecentCommits(3)
	if err != nil {
		t.Fatalf("GetRecentCommits() error = %v", err)
	}
	if len(commits) != 3 {
		t.Errorf("GetRecentCommits(3) length = %v, want 3", len(commits))
	}

	// Should be sorted by most recent first
	if commits[0].Model != "grok" {
		t.Errorf("GetRecentCommits() first commit Model = %v, want grok", commits[0].Model)
	}
	if commits[1].Model != "groq" {
		t.Errorf("GetRecentCommits() second commit Model = %v, want groq", commits[1].Model)
	}
	if commits[2].Model != "anthropic" {
		t.Errorf("GetRecentCommits() third commit Model = %v, want anthropic", commits[2].Model)
	}

	// Test with limit larger than total
	commits, err = GetRecentCommits(100)
	if err != nil {
		t.Fatalf("GetRecentCommits() error = %v", err)
	}
	if len(commits) != 5 {
		t.Errorf("GetRecentCommits(100) length = %v, want 5", len(commits))
	}
}

func TestLoadWithInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create invalid JSON stats file
	statsPath := filepath.Join(tmpDir, ".vibecheck_stats.json")
	invalidJSON := `{"commits": [{"timestamp": "invalid"}]}`
	if err := os.WriteFile(statsPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("Failed to write invalid stats: %v", err)
	}

	_, err := Load()
	if err == nil {
		t.Error("Load() with invalid JSON should return error")
	}
}

func TestGetMostUsedModelWithTie(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create a tie situation
	RecordCommit("openai", 1.0, "commit 1")
	RecordCommit("gemini", 2.0, "commit 2")
	RecordCommit("openai", 1.5, "commit 3")
	RecordCommit("gemini", 2.5, "commit 4")

	model, err := GetMostUsedModel()
	if err != nil {
		t.Fatalf("GetMostUsedModel() error = %v", err)
	}
	// Should return one of the tied models
	if model != "openai" && model != "gemini" {
		t.Errorf("GetMostUsedModel() = %v, want openai or gemini", model)
	}
}
