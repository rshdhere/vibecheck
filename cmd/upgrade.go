/*
Copyright © 2025 raashed
*/
package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/rshdhere/vibecheck/internal/ui/banner"
)

const (
	githubAPIURL = "https://api.github.com/repos/rshdhere/vibecheck/releases/latest"
	repoURL      = "rshdhere/vibecheck"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Short:   "Upgrade vibecheck to the latest version",
	Long:    `Automatically downloads and installs the latest version of vibecheck from GitHub releases.`,
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get current executable path
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}
		execPath, err = filepath.EvalSymlinks(execPath)
		if err != nil {
			return fmt.Errorf("failed to resolve executable path: %w", err)
		}

		// Check if we need elevated privileges
		if !isWritable(filepath.Dir(execPath)) {
			if runtime.GOOS == "windows" {
				return fmt.Errorf("insufficient permissions to upgrade. Please run as Administrator")
			}

			// Check if we're already running with sudo
			if os.Geteuid() != 0 {
				fmt.Println("🔐 Elevated privileges required. Re-running with sudo...")
				return rerunWithSudo()
			}
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithColor("cyan"))
		s.Suffix = " Checking for updates..."
		s.Start()

		// Fetch latest release info
		release, err := fetchLatestRelease()
		if err != nil {
			s.Stop()
			return fmt.Errorf("failed to fetch latest release: %w", err)
		}

		latestVersion := strings.TrimPrefix(release.TagName, "v")
		currentVersion := strings.TrimPrefix(version, "v")

		// Remove git describe suffix for comparison (e.g., -3-ge034ae7-dirty)
		if idx := strings.Index(currentVersion, "-"); idx != -1 {
			currentVersion = currentVersion[:idx]
		}

		s.Stop()

		if currentVersion == latestVersion {
			fmt.Printf("✅ Already on the latest version: %s\n", version)
			fmt.Println()
			banner.Print()
			return nil
		}

		fmt.Printf("Current version: %s\n", version)
		fmt.Printf("Latest version: %s\n", release.TagName)
		fmt.Printf("\n")

		// Find the appropriate asset for current OS/arch
		assetName, downloadURL := findAssetForPlatform(release)
		if assetName == "" {
			return fmt.Errorf("no compatible release found for %s/%s", runtime.GOOS, runtime.GOARCH)
		}

		s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithColor("cyan"))
		s.Suffix = fmt.Sprintf(" Downloading %s...", assetName)
		s.Start()

		// Download the asset
		tmpDir, err := os.MkdirTemp("", "vibecheck-upgrade-*")
		if err != nil {
			s.Stop()
			return fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer os.RemoveAll(tmpDir)

		tmpFile := filepath.Join(tmpDir, assetName)
		if err := downloadFile(tmpFile, downloadURL); err != nil {
			s.Stop()
			return fmt.Errorf("failed to download release: %w", err)
		}

		s.Suffix = " Extracting..."

		// Extract the binary
		binaryPath, err := extractBinary(tmpFile, tmpDir)
		if err != nil {
			s.Stop()
			return fmt.Errorf("failed to extract binary: %w", err)
		}

		s.Suffix = " Installing..."

		// Replace current binary
		if err := replaceBinary(execPath, binaryPath); err != nil {
			s.Stop()
			return fmt.Errorf("failed to replace binary: %w", err)
		}

		s.Stop()
		fmt.Printf("Successfully upgraded to version %s!\n", release.TagName)
		fmt.Printf("   Run 'vibecheck --version' to verify.\n\n")
		banner.Print()

		return nil
	},
}

func fetchLatestRelease() (*GitHubRelease, error) {
	resp, err := http.Get(githubAPIURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func findAssetForPlatform(release *GitHubRelease) (string, string) {
	osName := runtime.GOOS
	archName := runtime.GOARCH

	// Map Go OS names to release naming
	osMap := map[string]string{
		"darwin":  "Darwin",
		"linux":   "Linux",
		"windows": "Windows",
	}

	// Map Go arch names to release naming
	archMap := map[string]string{
		"amd64": "x86_64",
		"386":   "i386",
		"arm64": "arm64",
	}

	releaseOS := osMap[osName]
	releaseArch := archMap[archName]

	// Expected format: vibecheck_OS_ARCH.tar.gz or .zip
	var expectedExt string
	if osName == "windows" {
		expectedExt = ".zip"
	} else {
		expectedExt = ".tar.gz"
	}

	expectedName := fmt.Sprintf("vibecheck_%s_%s%s", releaseOS, releaseArch, expectedExt)

	for _, asset := range release.Assets {
		if asset.Name == expectedName {
			return asset.Name, asset.BrowserDownloadURL
		}
	}

	return "", ""
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractBinary(archivePath, destDir string) (string, error) {
	if strings.HasSuffix(archivePath, ".tar.gz") {
		return extractTarGz(archivePath, destDir)
	} else if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, destDir)
	}
	return "", fmt.Errorf("unsupported archive format")
}

func extractTarGz(archivePath, destDir string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	var binaryPath string
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		target := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return "", err
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return "", err
			}
			outFile.Close()

			// Find the vibecheck binary
			if filepath.Base(header.Name) == "vibecheck" {
				binaryPath = target
			}
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("vibecheck binary not found in archive")
	}

	return binaryPath, nil
}

func extractZip(archivePath, destDir string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var binaryPath string
	for _, f := range r.File {
		target := filepath.Join(destDir, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return "", err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return "", err
		}

		outFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return "", err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return "", err
		}

		// Find the vibecheck binary
		baseName := filepath.Base(f.Name)
		if baseName == "vibecheck.exe" || baseName == "vibecheck" {
			binaryPath = target
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("vibecheck binary not found in archive")
	}

	return binaryPath, nil
}

func replaceBinary(oldPath, newPath string) error {
	// Get the permissions of the old binary
	info, err := os.Stat(oldPath)
	if err != nil {
		return err
	}

	// Backup the old binary
	backupPath := oldPath + ".backup"
	if err := os.Rename(oldPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Copy new binary to the old location
	if err := copyFile(newPath, oldPath); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, oldPath)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Set the same permissions as the old binary
	if err := os.Chmod(oldPath, info.Mode()); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Remove backup on success
	os.Remove(backupPath)

	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// isWritable checks if the directory is writable by the current user
func isWritable(dir string) bool {
	// Try to create a temporary file in the directory
	testFile := filepath.Join(dir, ".vibecheck-write-test")
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return false
	}
	f.Close()
	os.Remove(testFile)
	return true
}

// rerunWithSudo re-executes the current command with sudo on Unix-like systems
func rerunWithSudo() error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("automatic privilege elevation not supported on Windows")
	}

	// Get the path to sudo
	sudoPath, err := exec.LookPath("sudo")
	if err != nil {
		return fmt.Errorf("sudo not found. Please run: sudo vibecheck upgrade")
	}

	// Prepare arguments: sudo vibecheck upgrade [original args]
	args := append([]string{os.Args[0]}, os.Args[1:]...)

	// Execute sudo with the current command
	cmd := exec.Command(sudoPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Preserve the exit code from the sudo command
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		return fmt.Errorf("failed to run with sudo: %w", err)
	}

	// Exit successfully after sudo command completes
	os.Exit(0)
	return nil
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
