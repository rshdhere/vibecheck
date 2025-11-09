// Package stats handles tracking and retrieval of commit statistics
package stats

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// CommitRecord represents a single commit record
type CommitRecord struct {
	Timestamp time.Time `json:"timestamp"`
	Model     string    `json:"model"`
	Latency   float64   `json:"latency"` // in seconds
	CommitMsg string    `json:"commit_msg"`
}

// Stats represents all commit statistics
type Stats struct {
	Commits []CommitRecord `json:"commits"`
}

// getStatsPath returns the path to the stats file
func getStatsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".vibecheck_stats.json"), nil
}

// Load reads statistics from disk
func Load() (*Stats, error) {
	path, err := getStatsPath()
	if err != nil {
		return nil, err
	}

	// If stats doesn't exist, return empty stats
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Stats{Commits: []CommitRecord{}}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var stats Stats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// Save writes statistics to disk
func Save(stats *Stats) error {
	path, err := getStatsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// RecordCommit adds a new commit record to the stats
func RecordCommit(model string, latency float64, commitMsg string) error {
	stats, err := Load()
	if err != nil {
		return err
	}

	record := CommitRecord{
		Timestamp: time.Now(),
		Model:     model,
		Latency:   latency,
		CommitMsg: commitMsg,
	}

	stats.Commits = append(stats.Commits, record)
	return Save(stats)
}

// GetTotalCommits returns the total number of commits
func GetTotalCommits() (int, error) {
	stats, err := Load()
	if err != nil {
		return 0, err
	}
	return len(stats.Commits), nil
}

// GetMostUsedModel returns the most frequently used model
func GetMostUsedModel() (string, error) {
	stats, err := Load()
	if err != nil {
		return "", err
	}

	if len(stats.Commits) == 0 {
		return "N/A", nil
	}

	modelCount := make(map[string]int)
	for _, commit := range stats.Commits {
		modelCount[commit.Model]++
	}

	maxCount := 0
	mostUsed := ""
	for model, count := range modelCount {
		if count > maxCount {
			maxCount = count
			mostUsed = model
		}
	}

	return mostUsed, nil
}

// GetAverageLatency returns the average latency in seconds
func GetAverageLatency() (float64, error) {
	stats, err := Load()
	if err != nil {
		return 0, err
	}

	if len(stats.Commits) == 0 {
		return 0, nil
	}

	var total float64
	for _, commit := range stats.Commits {
		total += commit.Latency
	}

	return total / float64(len(stats.Commits)), nil
}

// GetLastUsed returns the time of the last commit
func GetLastUsed() (time.Time, error) {
	stats, err := Load()
	if err != nil {
		return time.Time{}, err
	}

	if len(stats.Commits) == 0 {
		return time.Time{}, nil
	}

	// Commits are appended, so the last one is the most recent
	return stats.Commits[len(stats.Commits)-1].Timestamp, nil
}

// GetRecentCommits returns the most recent commits (up to limit)
func GetRecentCommits(limit int) ([]CommitRecord, error) {
	stats, err := Load()
	if err != nil {
		return nil, err
	}

	if len(stats.Commits) == 0 {
		return []CommitRecord{}, nil
	}

	// Sort by timestamp descending
	commits := make([]CommitRecord, len(stats.Commits))
	copy(commits, stats.Commits)
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Timestamp.After(commits[j].Timestamp)
	})

	if limit > len(commits) {
		limit = len(commits)
	}

	return commits[:limit], nil
}
