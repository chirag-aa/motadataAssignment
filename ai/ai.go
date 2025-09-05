package ai

import (
	"context"
	"fmt"
	"motadataAssignment/models"
	"strings"
	"time"
)

// Client is the AI integration interface
type Client interface {
	Summarize(ctx context.Context, query string, articles []models.Article) (string, []string, error)
}

// SimulatedAI is an implementation used for testing/demo without real API calls
type SimulatedAI struct{}

// NewSimulated returns a simulated AI client
func NewSimulated() *SimulatedAI {
	return &SimulatedAI{}
}

// Summarize returns a short summary and list of article titles (IDs) considered relevant.
func (s *SimulatedAI) Summarize(ctx context.Context, query string, articles []models.Article) (string, []string, error) {
	// simulate latency
	select {
	case <-time.After(20 * time.Millisecond):
	case <-ctx.Done():
		return "", nil, ctx.Err()
	}

	if len(articles) == 0 {
		return fmt.Sprintf("I couldn't find any KB articles relevant to \"%s\". Try rephrasing or adding details.", query),
			[]string{}, nil
	}

	// generate a concise answer by taking first sentences (naive)
	var snippets []string
	var ids []string
	for _, a := range articles {
		ids = append(ids, a.ID)
		first := firstSentence(a.Content)
		if first != "" {
			snippets = append(snippets, fmt.Sprintf("%s: %s", a.Title, first))
		} else {
			snippets = append(snippets, fmt.Sprintf("%s", a.Title))
		}
	}

	answer := fmt.Sprintf("Based on %d article(s): %s", len(articles), strings.Join(snippets, " | "))
	// keep answer concise
	if len(answer) > 700 {
		answer = answer[:700] + "..."
	}
	return answer, ids, nil
}

func firstSentence(s string) string {
	s = strings.TrimSpace(s)
	for _, sep := range []string{".", "!", "?"} {
		if idx := strings.Index(s, sep); idx != -1 {
			return strings.TrimSpace(s[:idx+1])
		}
	}
	return s
}
