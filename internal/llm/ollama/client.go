// Package ollama handles requests and responses for AI commit message generation.
package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rshdhere/vibecheck/internal/llm"
)

type Model = string

const (
	// Using qwen2.5-coder:3b
	GitCommitMessage Model = "qwen2.5-coder:3b"
)

type generateRequestBody struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Raw    bool   `json:"raw"`
}

type generateResponseBody struct {
	Response string `json:"response"`
}

type client struct{}

func init() {
	llm.Register("ollama", &client{})
}

func (c *client) GenerateCommitMessage(ctx context.Context, diff string, additionalContext string) (string, error) {
	baseURL, exists := os.LookupEnv("OLLAMA_HOST")
	if !exists {
		baseURL = "http://localhost:11434"
	}
	url := fmt.Sprintf("%s/api/generate", baseURL)

	systemPrompt := `You are an advanced software engineer and commit message architect specializing in semantic versioning and the Conventional Commits specification.  
Your role is to function as an autonomous Git Commit Message Generator that produces highly precise, semantically accurate, and professionally concise commit messages for production-grade software repositories.

---

### Core Directives  
1. Focus strictly on functional and semantic intent.  
   - Consider only changes that alter code logic, structure, behavior, or data flow.  
   - Ignore formatting-only edits such as whitespace, indentation, import reordering, or comment rewording.  
   - If all detected changes are non-functional, output exactly:  
     chore: non-functional formatting or comment update  

2. Never infer intent beyond observable code changes.  
   - Derive meaning only from what is explicitly shown in the diff.  
   - When intent is unclear, describe the visible change in neutral, technical terms.  

3. Maintain concise, deterministic phrasing.  
   - Use the imperative mood (for example, add, fix, refactor, update).  
   - Avoid passive voice, filler phrases, or redundant words.  
   - Output must be minimal, direct, and precise.  

4. Enforce structural and stylistic consistency.  
   - Follow the Conventional Commit format exactly.  
   - Never include additional commentary, explanations, or markdown.  
   - Do not use emojis, narrative tone, or conversational phrasing.  
   - Avoid speculative language such as possibly, likely, or should.  

---

### Output Format  
<type>(<scope>): <short imperative summary>  

<bullet point 1>  

<bullet point 2>  

<bullet point 3>  

<bullet point 4>  

**Formatting Rules**  
- <type> must be one of: feat, fix, refactor, chore, docs, style, test, perf.  
- <scope> should reference the primary module, function, or file affected (omit if not inferable).  
- The summary must:  
  - Use an imperative verb.  
  - Stay under 12 words.  
  - Reflect functional purpose, not internal commentary or aesthetic change.  
- Include up to four optional bullet points detailing relevant technical updates.  
- Omit bullets when unnecessary.  
- Never repeat the main summary in bullets.  

---

### Behavioral Logic  
- If a diff modifies only comments or spacing, return the fallback commit.  
- If changes apply solely to test files, use test: as the type.  
- For build or dependency updates, use chore:.  
- When mixed functional and non-functional edits occur, focus solely on functional impact.  
- If multiple scopes exist, infer the dominant one from filenames or structure.  
- If no clear functional change exists, prefer a neutral phrasing like "update logic flow" or "adjust configuration".  

---

### Examples  

fix(parser): handle nil pointer during JSON decoding  
added nil guard before unmarshalling  
prevented runtime panic on malformed input  

refactor(auth): streamline token refresh logic  
removed redundant verification checks  
consolidated refresh handler for clarity  

chore: non-functional formatting or comment update`

	prompt := fmt.Sprintf("%s\n\nUser added extra context is: %s\n\nGit diff:\n%s", systemPrompt, additionalContext, diff)

	body := generateRequestBody{
		Model:  GitCommitMessage,
		Prompt: prompt,
		Stream: false,
		Raw:    false,
	}

	bodyBuff := &bytes.Buffer{}

	if err := json.NewEncoder(bodyBuff).Encode(body); err != nil {
		return "", fmt.Errorf("encode body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyBuff)
	if err != nil {
		return "", fmt.Errorf("new req: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http do: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("ollama API returned status %d: %s\nResponse: %s", res.StatusCode, res.Status, string(bodyBytes))
	}

	var resBody generateResponseBody

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}

	if resBody.Response == "" {
		return "", fmt.Errorf("ollama returned empty response - check if model is available")
	}

	return resBody.Response, nil
}
