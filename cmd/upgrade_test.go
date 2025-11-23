package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFindAssetForPlatform(t *testing.T) {
	tests := []struct {
		name     string
		release  *GitHubRelease
		osName   string
		archName string
		wantName string
		wantURL  string
	}{
		{
			name: "linux amd64",
			release: &GitHubRelease{
				Assets: []struct {
					Name               string `json:"name"`
					BrowserDownloadURL string `json:"browser_download_url"`
				}{
					{Name: "vibecheck_Linux_x86_64.tar.gz", BrowserDownloadURL: "https://example.com/vibecheck_Linux_x86_64.tar.gz"},
				},
			},
			osName:   "linux",
			archName: "amd64",
			wantName: "vibecheck_Linux_x86_64.tar.gz",
			wantURL:  "https://example.com/vibecheck_Linux_x86_64.tar.gz",
		},
		{
			name: "darwin arm64",
			release: &GitHubRelease{
				Assets: []struct {
					Name               string `json:"name"`
					BrowserDownloadURL string `json:"browser_download_url"`
				}{
					{Name: "vibecheck_Darwin_arm64.tar.gz", BrowserDownloadURL: "https://example.com/vibecheck_Darwin_arm64.tar.gz"},
				},
			},
			osName:   "darwin",
			archName: "arm64",
			wantName: "vibecheck_Darwin_arm64.tar.gz",
			wantURL:  "https://example.com/vibecheck_Darwin_arm64.tar.gz",
		},
		{
			name: "windows amd64",
			release: &GitHubRelease{
				Assets: []struct {
					Name               string `json:"name"`
					BrowserDownloadURL string `json:"browser_download_url"`
				}{
					{Name: "vibecheck_Windows_x86_64.zip", BrowserDownloadURL: "https://example.com/vibecheck_Windows_x86_64.zip"},
				},
			},
			osName:   "windows",
			archName: "amd64",
			wantName: "vibecheck_Windows_x86_64.zip",
			wantURL:  "https://example.com/vibecheck_Windows_x86_64.zip",
		},
		{
			name: "not found",
			release: &GitHubRelease{
				Assets: []struct {
					Name               string `json:"name"`
					BrowserDownloadURL string `json:"browser_download_url"`
				}{
					{Name: "vibecheck_Linux_x86_64.tar.gz", BrowserDownloadURL: "https://example.com/vibecheck_Linux_x86_64.tar.gz"},
				},
			},
			osName:   "windows",
			archName: "amd64",
			wantName: "",
			wantURL:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: We can't actually override runtime.GOOS/GOARCH, so we'll test
			// with the current platform or skip platform-specific tests
			// For now, we'll test the logic with a mock approach
			if tt.osName != runtime.GOOS || tt.archName != runtime.GOARCH {
				// Skip tests for other platforms
				t.Skipf("Skipping test for %s/%s (current: %s/%s)", tt.osName, tt.archName, runtime.GOOS, runtime.GOARCH)
			}

			gotName, gotURL := findAssetForPlatform(tt.release)
			if gotName != tt.wantName {
				t.Errorf("findAssetForPlatform() name = %v, want %v", gotName, tt.wantName)
			}
			if gotURL != tt.wantURL {
				t.Errorf("findAssetForPlatform() url = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func TestIsWritable(t *testing.T) {
	t.Run("writable directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		if !isWritable(tmpDir) {
			t.Error("isWritable() = false for writable directory")
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		if isWritable("/nonexistent/directory/path") {
			t.Error("isWritable() = true for non-existent directory")
		}
	})
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "source.txt")
	dst := filepath.Join(tmpDir, "dest.txt")

	// Create source file
	content := "test content"
	if err := os.WriteFile(src, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Copy file
	if err := copyFile(src, dst); err != nil {
		t.Fatalf("copyFile() error = %v", err)
	}

	// Verify destination exists and has correct content
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if string(data) != content {
		t.Errorf("copyFile() destination content = %q, want %q", string(data), content)
	}
}

func TestExtractBinary(t *testing.T) {
	t.Run("unsupported format", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.unknown")
		os.WriteFile(testFile, []byte("test"), 0644)

		_, err := extractBinary(testFile, tmpDir)
		if err == nil {
			t.Error("extractBinary() with unsupported format should return error")
		}
	})
}
