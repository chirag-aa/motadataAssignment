package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"motadataAssignment/models"
	"net/http"
	"os"
	"strings"
	"time"
)

// OpenAIClient implements ai.Client using OpenAI API
type OpenAIClient struct {
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

// NewOpenAIClient creates a new client.
// It reads API key from OPENAI_API_KEY if apiKey == "".
func NewOpenAIClient(apiKey string, model string) *OpenAIClient {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if model == "" {
		model = "gpt-4o-mini" // default, can be overridden
	}
	return &OpenAIClient{
		APIKey:     apiKey,
		Model:      model,
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// request/response structs for OpenAI API
type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Summarize calls the OpenAI API
func (c *OpenAIClient) Summarize(ctx context.Context, query string, articles []models.Article) (string, []string, error) {
	if c.APIKey == "" {
		return "", nil, fmt.Errorf("missing OPENAI_API_KEY")
	}

	// Prepare context: concatenate query + articles
	var b strings.Builder
	fmt.Fprintf(&b, "User query: %s\n\n", query)
	fmt.Fprintf(&b, "Knowledge Base Articles:\n")
	for _, a := range articles {
		fmt.Fprintf(&b, "- ID: %s\nTitle: %s\nContent: %s\n\n", a.ID, a.Title, a.Content)
	}
	fmt.Fprintf(&b, "Task: Provide a concise helpful answer to the user's query based ONLY on the KB articles above. Also output the relevant Article IDs as a JSON array.")

	reqBody := chatRequest{
		Model: c.Model,
		Messages: []chatMessage{
			{Role: "system", Content: "You are an IT help assistant."},
			{Role: "user", Content: b.String()},
		},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("openai api error: %s", resp.Status)
	}

	var cr chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return "", nil, err
	}
	if len(cr.Choices) == 0 {
		return "", nil, fmt.Errorf("no choices returned")
	}

	answer := cr.Choices[0].Message.Content

	// naive parsing: expect model to produce something like:
	// "Answer: ... \nRelevantIDs: [\"id1\", \"id2\"]"
	ids := extractIDs(answer)

	return answer, ids, nil
}

// extractIDs tries to parse article IDs from model output (very naive JSON detection)
func extractIDs(s string) []string {
	start := strings.Index(s, "[")
	end := strings.Index(s, "]")
	if start == -1 || end == -1 || start > end {
		return nil
	}
	sub := s[start : end+1]
	var ids []string
	_ = json.Unmarshal([]byte(sub), &ids)
	return ids
}
