package anthropic

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
	t.Setenv("ANTHROPIC_API_KEY", "")
}

// serve points the SDK at a local server via ANTHROPIC_BASE_URL.
func serve(t *testing.T, status int, body string) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/messages" {
			t.Errorf("request = %s %s, want POST /v1/messages", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("X-Api-Key"); got != "test-key" {
			t.Errorf("X-Api-Key = %q, want test-key", got)
		}
		var req struct {
			Model    string `json:"model"`
			Messages []struct {
				Content []struct {
					Text string `json:"text"`
				} `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Model != "claude-3-5-haiku-20241022" {
			t.Errorf("model = %q, want claude-3-5-haiku-20241022", req.Model)
		}
		if len(req.Messages) != 1 || len(req.Messages[0].Content) != 1 || !strings.Contains(req.Messages[0].Content[0].Text, "the diff") {
			t.Errorf("messages = %+v, want a single user message containing the diff", req.Messages)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	t.Setenv("ANTHROPIC_BASE_URL", srv.URL)
}

func message(content string) string {
	return `{"id":"msg_1","type":"message","role":"assistant","model":"claude-3-5-haiku-20241022","content":` + content +
		`,"stop_reason":"end_turn","usage":{"input_tokens":1,"output_tokens":1}}`
}

func TestGenerateCommitMessageMissingKey(t *testing.T) {
	isolate(t)

	_, err := (&client{}).GenerateCommitMessage(context.Background(), "the diff", "")
	if err == nil || err.Error() != "ANTHROPIC_API_KEY environment variable not set" {
		t.Fatalf("err = %v, want missing ANTHROPIC_API_KEY error", err)
	}
}

func TestGenerateCommitMessageSuccess(t *testing.T) {
	isolate(t)
	t.Setenv("ANTHROPIC_API_KEY", "test-key")
	serve(t, http.StatusOK, message(`[{"type":"text","text":"feat: add thing"}]`))

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
		{"api error", http.StatusBadRequest, `{"type":"error","error":{"type":"invalid_request_error","message":"bad"}}`, "error while prompting to Anthropic"},
		{"empty content", http.StatusOK, message(`[]`), "no response generated"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolate(t)
			t.Setenv("ANTHROPIC_API_KEY", "test-key")
			serve(t, tt.status, tt.body)

			_, err := (&client{}).GenerateCommitMessage(context.Background(), "the diff", "")
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}
