package groq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	envVar    = "GROQ_API_KEY"
	wantModel = "llama-3.3-70b-versatile"
)

// isolate keeps the developer's real keys file and environment out of the test.
func isolate(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv(envVar, "")
}

// point directs the client at url for the duration of the test.
func point(t *testing.T, url string) {
	t.Helper()
	old := baseURL
	baseURL = url
	t.Cleanup(func() { baseURL = old })
}

// serve starts an OpenAI-compatible stand-in that replies with status and body.
func serve(t *testing.T, status int, body string) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || !strings.HasSuffix(r.URL.Path, "/chat/completions") {
			t.Errorf("request = %s %s, want POST .../chat/completions", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization = %q, want Bearer test-key", got)
		}
		var req struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Model != wantModel {
			t.Errorf("model = %q, want %q", req.Model, wantModel)
		}
		if len(req.Messages) != 3 || req.Messages[0].Role != "system" || req.Messages[2].Content != "the diff" {
			t.Errorf("messages = %+v, want system + context + diff", req.Messages)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	point(t, srv.URL)
}

func completion(choices string) string {
	return `{"id":"c1","object":"chat.completion","created":0,"model":"` + wantModel + `","choices":` + choices + `}`
}

func TestGenerateCommitMessageMissingKey(t *testing.T) {
	isolate(t)

	_, err := (&client{}).GenerateCommitMessage(context.Background(), "the diff", "")
	if err == nil || err.Error() != envVar+" environment variable not set" {
		t.Fatalf("err = %v, want missing %s error", err, envVar)
	}
}

func TestGenerateCommitMessageSuccess(t *testing.T) {
	isolate(t)
	t.Setenv(envVar, "test-key")
	serve(t, http.StatusOK, completion(`[{"index":0,"finish_reason":"stop","message":{"role":"assistant","content":"feat: add thing"}}]`))

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
		{"api error", http.StatusBadRequest, `{"error":{"message":"bad request"}}`, "400"},
		{"no choices", http.StatusOK, completion(`[]`), "no response choices"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolate(t)
			t.Setenv(envVar, "test-key")
			serve(t, tt.status, tt.body)

			_, err := (&client{}).GenerateCommitMessage(context.Background(), "the diff", "")
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}
