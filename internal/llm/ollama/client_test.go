package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// isolate keeps the developer's real keys file and environment out of the test.
func isolate(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("OLLAMA_HOST", "")
}

// serve starts a local Ollama stand-in and points OLLAMA_HOST at it.
func serve(t *testing.T, h http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	t.Setenv("OLLAMA_HOST", srv.URL)
}

func TestGenerateCommitMessageSuccess(t *testing.T) {
	isolate(t)
	serve(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/generate" {
			t.Errorf("request = %s %s, want POST /api/generate", r.Method, r.URL.Path)
		}
		var req generateRequestBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Model != GitCommitMessage {
			t.Errorf("model = %q, want %q", req.Model, GitCommitMessage)
		}
		if req.Stream {
			t.Error("stream = true, want false")
		}
		if !strings.Contains(req.Prompt, "extra context") || !strings.Contains(req.Prompt, "the diff") {
			t.Errorf("prompt does not include the context and the diff: %q", req.Prompt)
		}
		w.Write([]byte(`{"response":"feat: add thing"}`))
	})

	got, err := (&client{}).GenerateCommitMessage(context.Background(), "the diff", "extra context")
	if err != nil {
		t.Fatalf("GenerateCommitMessage() error = %v", err)
	}
	if got != "feat: add thing" {
		t.Errorf("GenerateCommitMessage() = %q, want %q", got, "feat: add thing")
	}
}

func TestGenerateCommitMessageResponseErrors(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantErr string
	}{
		{"non-200 status", http.StatusNotFound, `{"error":"model 'gpt-oss:20b' not found"}`, "status 404"},
		{"invalid json", http.StatusOK, "not json", "decode"},
		{"empty response", http.StatusOK, `{"response":""}`, "empty response"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolate(t)
			serve(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.body))
			})

			_, err := (&client{}).GenerateCommitMessage(context.Background(), "diff", "")
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateCommitMessageRequestErrors(t *testing.T) {
	closed := httptest.NewServer(http.NotFoundHandler())
	closed.Close()

	tests := []struct {
		name    string
		host    string
		wantErr string
	}{
		{"invalid host", "://bad-host", "new req"},
		{"unreachable host", closed.URL, "http do"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolate(t)
			t.Setenv("OLLAMA_HOST", tt.host)

			_, err := (&client{}).GenerateCommitMessage(context.Background(), "diff", "")
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}
