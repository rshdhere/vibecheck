package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

type archiveEntry struct {
	name string
	body string // empty for directories
	dir  bool
}

func writeTarGz(t *testing.T, path string, entries ...archiveEntry) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	for _, e := range entries {
		hdr := &tar.Header{Name: e.name, Mode: 0755, Size: int64(len(e.body)), Typeflag: tar.TypeReg}
		if e.dir {
			hdr.Typeflag, hdr.Size = tar.TypeDir, 0
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(e.body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
}

func writeZip(t *testing.T, path string, entries ...archiveEntry) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	for _, e := range entries {
		name := e.name
		if e.dir {
			name += "/"
		}
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(e.body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestExtractTarGz(t *testing.T) {
	t.Run("finds binary", func(t *testing.T) {
		dir := t.TempDir()
		archive := filepath.Join(dir, "release.tar.gz")
		writeTarGz(t, archive,
			archiveEntry{name: "docs", dir: true},
			archiveEntry{name: "docs/README.md", body: "readme"},
			archiveEntry{name: "vibecheck", body: "new binary"},
		)

		got, err := extractBinary(archive, dir)
		if err != nil {
			t.Fatalf("extractBinary() error = %v", err)
		}
		if got != filepath.Join(dir, "vibecheck") || readFile(t, got) != "new binary" {
			t.Errorf("extracted %q with unexpected contents", got)
		}
	})

	t.Run("missing binary", func(t *testing.T) {
		dir := t.TempDir()
		archive := filepath.Join(dir, "release.tar.gz")
		writeTarGz(t, archive, archiveEntry{name: "README.md", body: "readme"})

		if _, err := extractBinary(archive, dir); err == nil || !strings.Contains(err.Error(), "not found") {
			t.Errorf("err = %v, want binary not found", err)
		}
	})

	t.Run("not gzip", func(t *testing.T) {
		dir := t.TempDir()
		archive := filepath.Join(dir, "release.tar.gz")
		os.WriteFile(archive, []byte("plain text"), 0644)

		if _, err := extractBinary(archive, dir); err == nil {
			t.Error("expected an error for a non-gzip archive")
		}
	})

	t.Run("missing archive", func(t *testing.T) {
		if _, err := extractTarGz(filepath.Join(t.TempDir(), "nope.tar.gz"), t.TempDir()); err == nil {
			t.Error("expected an error for a missing archive")
		}
	})
}

func TestExtractZip(t *testing.T) {
	t.Run("finds binary", func(t *testing.T) {
		dir := t.TempDir()
		archive := filepath.Join(dir, "release.zip")
		writeZip(t, archive,
			archiveEntry{name: "docs", dir: true},
			archiveEntry{name: "docs/README.md", body: "readme"},
			archiveEntry{name: "vibecheck.exe", body: "new binary"},
		)

		got, err := extractBinary(archive, dir)
		if err != nil {
			t.Fatalf("extractBinary() error = %v", err)
		}
		if got != filepath.Join(dir, "vibecheck.exe") || readFile(t, got) != "new binary" {
			t.Errorf("extracted %q with unexpected contents", got)
		}
	})

	t.Run("missing binary", func(t *testing.T) {
		dir := t.TempDir()
		archive := filepath.Join(dir, "release.zip")
		writeZip(t, archive, archiveEntry{name: "README.md", body: "readme"})

		if _, err := extractBinary(archive, dir); err == nil || !strings.Contains(err.Error(), "not found") {
			t.Errorf("err = %v, want binary not found", err)
		}
	})

	t.Run("not a zip", func(t *testing.T) {
		dir := t.TempDir()
		archive := filepath.Join(dir, "release.zip")
		os.WriteFile(archive, []byte("plain text"), 0644)

		if _, err := extractBinary(archive, dir); err == nil {
			t.Error("expected an error for an invalid zip")
		}
	})
}

func TestReplaceBinary(t *testing.T) {
	t.Run("replaces and keeps mode", func(t *testing.T) {
		dir := t.TempDir()
		oldPath, newPath := filepath.Join(dir, "vibecheck"), filepath.Join(dir, "new")
		os.WriteFile(oldPath, []byte("old"), 0700)
		os.WriteFile(newPath, []byte("new"), 0644)

		if err := replaceBinary(oldPath, newPath); err != nil {
			t.Fatalf("replaceBinary() error = %v", err)
		}
		if readFile(t, oldPath) != "new" {
			t.Error("binary was not replaced")
		}
		if info, _ := os.Stat(oldPath); info.Mode().Perm() != 0700 {
			t.Errorf("mode = %v, want 0700", info.Mode().Perm())
		}
		if _, err := os.Stat(oldPath + ".backup"); !os.IsNotExist(err) {
			t.Error("backup was not removed")
		}
	})

	t.Run("missing current binary", func(t *testing.T) {
		dir := t.TempDir()
		if err := replaceBinary(filepath.Join(dir, "vibecheck"), filepath.Join(dir, "new")); err == nil {
			t.Error("expected an error when the current binary is missing")
		}
	})

	t.Run("restores backup when copy fails", func(t *testing.T) {
		dir := t.TempDir()
		oldPath := filepath.Join(dir, "vibecheck")
		os.WriteFile(oldPath, []byte("old"), 0755)

		err := replaceBinary(oldPath, filepath.Join(dir, "missing"))
		if err == nil || !strings.Contains(err.Error(), "failed to install new binary") {
			t.Fatalf("err = %v, want install failure", err)
		}
		if readFile(t, oldPath) != "old" {
			t.Error("original binary was not restored")
		}
	})
}

func TestCopyFileErrors(t *testing.T) {
	dir := t.TempDir()
	if err := copyFile(filepath.Join(dir, "missing"), filepath.Join(dir, "dst")); err == nil {
		t.Error("expected an error for a missing source")
	}
	src := filepath.Join(dir, "src")
	os.WriteFile(src, []byte("x"), 0644)
	if err := copyFile(src, filepath.Join(dir, "no-such-dir", "dst")); err == nil {
		t.Error("expected an error for an unwritable destination")
	}
}

func TestDownloadFile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/missing" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte("payload"))
	}))
	defer srv.Close()
	dir := t.TempDir()

	dst := filepath.Join(dir, "asset")
	if err := downloadFile(dst, srv.URL+"/asset"); err != nil {
		t.Fatalf("downloadFile() error = %v", err)
	}
	if readFile(t, dst) != "payload" {
		t.Error("downloaded contents do not match")
	}

	if err := downloadFile(dst, srv.URL+"/missing"); err == nil || !strings.Contains(err.Error(), "status 404") {
		t.Errorf("err = %v, want status 404", err)
	}
	if err := downloadFile(filepath.Join(dir, "no-such-dir", "asset"), srv.URL+"/asset"); err == nil {
		t.Error("expected an error for an unwritable destination")
	}
	if err := downloadFile(dst, "://bad-url"); err == nil {
		t.Error("expected an error for an invalid url")
	}
}

func TestFetchLatestRelease(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantErr string
	}{
		{"success", http.StatusOK, `{"tag_name":"v9.9.9","assets":[{"name":"a","browser_download_url":"u"}]}`, ""},
		{"non-200", http.StatusForbidden, "rate limited", "status 403"},
		{"invalid json", http.StatusOK, "not json", "invalid character"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serveRelease(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.body))
			})

			release, err := fetchLatestRelease()
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			if err != nil || release.TagName != "v9.9.9" || len(release.Assets) != 1 {
				t.Fatalf("release = %+v, err = %v", release, err)
			}
		})
	}

	t.Run("unreachable", func(t *testing.T) {
		setVar(t, &githubAPIURL, "://bad-url")
		if _, err := fetchLatestRelease(); err == nil {
			t.Error("expected an error for an invalid url")
		}
	})
}

func setVar[T any](t *testing.T, v *T, val T) {
	t.Helper()
	old := *v
	*v = val
	t.Cleanup(func() { *v = old })
}

// serveRelease points githubAPIURL at a local server and returns its base URL.
func serveRelease(t *testing.T, h http.HandlerFunc) string {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	setVar(t, &githubAPIURL, srv.URL+"/releases/latest")
	return srv.URL
}

// platformAsset is the release asset name findAssetForPlatform expects on this machine.
func platformAsset() string {
	osName := map[string]string{"darwin": "Darwin", "linux": "Linux", "windows": "Windows"}[runtime.GOOS]
	arch := map[string]string{"amd64": "x86_64", "386": "i386", "arm64": "arm64"}[runtime.GOARCH]
	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	return "vibecheck_" + osName + "_" + arch + ext
}

// fakeInstall makes the upgrade command treat a scratch file as the running binary.
func fakeInstall(t *testing.T) string {
	t.Helper()
	exe := filepath.Join(t.TempDir(), "vibecheck")
	if err := os.WriteFile(exe, []byte("old binary"), 0755); err != nil {
		t.Fatal(err)
	}
	setVar(t, &executablePath, func() (string, error) { return exe, nil })
	setVar(t, &version, "v1.0.0")
	return exe
}

func TestUpgradeCmdInstallsLatestRelease(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("release archive fixture is tar.gz")
	}
	exe := fakeInstall(t)
	archive := filepath.Join(t.TempDir(), platformAsset())
	writeTarGz(t, archive, archiveEntry{name: "vibecheck", body: "new binary"})

	var base string
	base = serveRelease(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/download" {
			http.ServeFile(w, r, archive)
			return
		}
		fmt.Fprintf(w, `{"tag_name":"v2.0.0","assets":[{"name":%q,"browser_download_url":%q}]}`, platformAsset(), base+"/download")
	})

	if err := execute(t, "upgrade"); err != nil {
		t.Fatalf("upgrade error = %v", err)
	}
	if readFile(t, exe) != "new binary" {
		t.Error("binary was not replaced with the release")
	}
}

func TestUpgradeCmdAlreadyLatest(t *testing.T) {
	exe := fakeInstall(t)
	setVar(t, &version, "v2.0.0-3-gabc1234-dirty")
	serveRelease(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"tag_name":"v2.0.0"}`))
	})

	if err := execute(t, "upgrade"); err != nil {
		t.Fatalf("upgrade error = %v", err)
	}
	if readFile(t, exe) != "old binary" {
		t.Error("binary should be untouched when already up to date")
	}
}

func TestUpgradeCmdErrors(t *testing.T) {
	brokenArchive := filepath.Join(t.TempDir(), "broken")
	os.WriteFile(brokenArchive, []byte("not an archive"), 0644)

	tests := []struct {
		name    string
		handler func(base string) http.HandlerFunc
		wantErr string
	}{
		{
			name: "release lookup fails",
			handler: func(string) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) }
			},
			wantErr: "failed to fetch latest release",
		},
		{
			name: "no asset for platform",
			handler: func(string) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"tag_name":"v2.0.0","assets":[]}`)) }
			},
			wantErr: "no compatible release found",
		},
		{
			name: "download fails",
			handler: func(base string) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/download" {
						http.NotFound(w, r)
						return
					}
					fmt.Fprintf(w, `{"tag_name":"v2.0.0","assets":[{"name":%q,"browser_download_url":%q}]}`, platformAsset(), base+"/download")
				}
			},
			wantErr: "failed to download release",
		},
		{
			name: "archive is corrupt",
			handler: func(base string) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/download" {
						http.ServeFile(w, r, brokenArchive)
						return
					}
					fmt.Fprintf(w, `{"tag_name":"v2.0.0","assets":[{"name":%q,"browser_download_url":%q}]}`, platformAsset(), base+"/download")
				}
			},
			wantErr: "failed to extract binary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeInstall(t)
			var base string
			base = serveRelease(t, func(w http.ResponseWriter, r *http.Request) { tt.handler(base)(w, r) })

			err := execute(t, "upgrade")
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("upgrade error = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}

	t.Run("executable path unavailable", func(t *testing.T) {
		setVar(t, &executablePath, func() (string, error) { return "", errors.New("no exe") })
		if err := execute(t, "upgrade"); err == nil || !strings.Contains(err.Error(), "failed to get executable path") {
			t.Fatalf("upgrade error = %v, want executable path failure", err)
		}
	})

	t.Run("executable path unresolvable", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "missing")
		setVar(t, &executablePath, func() (string, error) { return missing, nil })
		if err := execute(t, "upgrade"); err == nil || !strings.Contains(err.Error(), "failed to resolve executable path") {
			t.Fatalf("upgrade error = %v, want resolve failure", err)
		}
	})
}
