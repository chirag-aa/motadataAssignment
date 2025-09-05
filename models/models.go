package models

import "time"

// Article represents a KB article.
type Article struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// SearchRequest is the payload to POST /search-query
type SearchRequest struct {
	Query string `json:"query"`
}

// SearchResponse is returned by POST /search-query
type SearchResponse struct {
	AISummaryAnswer    string   `json:"ai_summary_answer"`
	AIRelevantArticles []string `json:"ai_relevant_articles"` // article IDs or titles
}

// StoredSearch is a record of the search and AI results (for persistence)
type StoredSearch struct {
	ID                 string    `json:"id"`
	Query              string    `json:"query"`
	AISummaryAnswer    string    `json:"ai_summary_answer"`
	AIRelevantArticles []string  `json:"ai_relevant_articles"`
	CreatedAt          time.Time `json:"created_at"`
}
