package gemini

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/api/option"
)

// isolate keeps the developer's real keys file and environment out of the test.
func isolate(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("GEMINI_API_KEY", "")
}

// serve points the client at a local server that replies with status and body.
func serve(t *testing.T, status int, body string) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "gemini-2.5-flash") {
			t.Errorf("path = %q, want it to target gemini-2.5-flash", r.URL.Path)
		}
		req, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(req), "the diff") {
			t.Errorf("request body does not include the diff: %s", req)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	old := endpointOptions
	endpointOptions = []option.ClientOption{option.WithEndpoint(srv.URL)}
	t.Cleanup(func() { endpointOptions = old })
}

func TestGenerateCommitMessageMissingKey(t *testing.T) {
	isolate(t)

	_, err := (&client{}).GenerateCommitMessage(context.Background(), "the diff", "")
	if err == nil || err.Error() != "GEMINI_API_KEY environment variable not set" {
		t.Fatalf("err = %v, want missing GEMINI_API_KEY error", err)
	}
}

func TestGenerateCommitMessageSuccess(t *testing.T) {
	isolate(t)
	t.Setenv("GEMINI_API_KEY", "test-key")
	serve(t, http.StatusOK, `{"candidates":[{"content":{"role":"model","parts":[{"text":"feat: add thing"}]},"finishReason":"STOP"}]}`)

	got, err := (&client{}).GenerateCommitMessage(context.Background(), "the diff", "ctx")
	if err != nil {
		t.Fatalf("GenerateCommitMessage() error = %v", err)
	}
	if got != "feat: add thing" {
		t.Errorf("GenerateCommitMessage() = %q, want %q", got, "feat: add thing")
	}
}

func TestGenerateCommitMessageErrors(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantErr string
	}{
		{"api error", http.StatusBadRequest, `{"error":{"code":400,"message":"bad request"}}`, "generate content"},
		{"no candidates", http.StatusOK, `{"candidates":[]}`, "no candidates"},
		{"unexpected finish reason", http.StatusOK, `{"candidates":[{"content":{"role":"model","parts":[{"text":"x"}]},"finishReason":"OTHER"}]}`, "finish reason"},
		{"empty content", http.StatusOK, `{"candidates":[{"finishReason":"STOP"}]}`, "empty content"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolate(t)
			t.Setenv("GEMINI_API_KEY", "test-key")
			serve(t, tt.status, tt.body)

			_, err := (&client{}).GenerateCommitMessage(context.Background(), "the diff", "")
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}
