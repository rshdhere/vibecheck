package qwen

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	envVar    = "QWEN_API_KEY"
	wantModel = "qwen-turbo"
)

// isolate keeps the developer's real keys file and environment out of the test.
func isolate(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv(envVar, "")
}

// serve points apiURL at a local server for the duration of the test.
func serve(t *testing.T, h http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	setAPIURL(t, srv.URL)
}

func setAPIURL(t *testing.T, url string) {
	t.Helper()
	old := apiURL
	apiURL = url
	t.Cleanup(func() { apiURL = old })
}

func TestGenerateCommitMessageMissingKey(t *testing.T) {
	isolate(t)

	_, err := (&client{}).GenerateCommitMessage(context.Background(), "diff", "")
	if err == nil || err.Error() != envVar+" environment variable not set" {
		t.Fatalf("err = %v, want missing %s error", err, envVar)
	}
}

func TestGenerateCommitMessageSuccess(t *testing.T) {
	isolate(t)
	t.Setenv(envVar, "test-key")

	serve(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization = %q, want Bearer test-key", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Model != wantModel {
			t.Errorf("model = %q, want %q", req.Model, wantModel)
		}
		if len(req.Messages) != 3 || req.Messages[0].Role != "system" {
			t.Fatalf("messages = %+v, want system + 2 user messages", req.Messages)
		}
		if !strings.Contains(req.Messages[1].Content, "extra context") {
			t.Errorf("context message = %q, want it to include the user context", req.Messages[1].Content)
		}
		if req.Messages[2].Content != "the diff" {
			t.Errorf("diff message = %q, want %q", req.Messages[2].Content, "the diff")
		}
		w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"feat: add thing"}}]}`))
	})

	got, err := (&client{}).GenerateCommitMessage(context.Background(), "the diff", "extra context")
	if err != nil {
		t.Fatalf("GenerateCommitMessage() error = %v", err)
	}
	if got != "feat: add thing" {
		t.Errorf("GenerateCommitMessage() = %q, want %q", got, "feat: add thing")
	}
}

func TestGenerateCommitMessageDashScopeOutputFormat(t *testing.T) {
	isolate(t)
	t.Setenv(envVar, "test-key")
	serve(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"output":{"choices":[{"message":{"role":"assistant","content":"fix: handle nil"}}]}}`))
	})

	got, err := (&client{}).GenerateCommitMessage(context.Background(), "diff", "")
	if err != nil {
		t.Fatalf("GenerateCommitMessage() error = %v", err)
	}
	if got != "fix: handle nil" {
		t.Errorf("GenerateCommitMessage() = %q, want %q", got, "fix: handle nil")
	}
}

func TestGenerateCommitMessageResponseErrors(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantErr string
	}{
		{"non-200 status", http.StatusUnauthorized, "bad key", "status 401: bad key"},
		{"invalid json", http.StatusOK, "not json", "decode response"},
		{"no choices", http.StatusOK, `{"choices":[]}`, "no response choices"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolate(t)
			t.Setenv(envVar, "test-key")
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
		url     string
		wantErr string
	}{
		{"invalid url", "://bad-url", "create request"},
		{"unreachable server", closed.URL, "send request"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolate(t)
			t.Setenv(envVar, "test-key")
			setAPIURL(t, tt.url)

			_, err := (&client{}).GenerateCommitMessage(context.Background(), "diff", "")
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}
