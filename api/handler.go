package api

import (
	"context"
	"encoding/json"
	"log"
	"motadataAssignment/ai"
	"motadataAssignment/kb"
	"motadataAssignment/models"
	"motadataAssignment/store"
	"net/http"
	"time"
)

// Handler contains dependencies
type Handler struct {
	KB    *kb.KB
	AI    ai.Client
	Store store.Store
}

// NewHandler creates a new Handler
func NewHandler(k *kb.KB, aiClient ai.Client, st store.Store) *Handler {
	return &Handler{
		KB:    k,
		AI:    aiClient,
		Store: st,
	}
}

// RegisterRoutes registers HTTP endpoints on the given mux
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/search-query", h.searchQueryHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
}

func (h *Handler) searchQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.SearchRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Query == "" {
		http.Error(w, "query cannot be empty", http.StatusBadRequest)
		return
	}

	// Find relevant KB articles (top 3)
	articles := h.KB.FindRelevant(req.Query, 3)

	// Call AI to summarize - with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	summary, articleIDs, err := h.AI.Summarize(ctx, req.Query, articles)
	if err != nil {
		log.Printf("ai summarize error: %v", err)
		http.Error(w, "internal ai error", http.StatusInternalServerError)
		return
	}

	// Save to store (non-blocking best-effort; but showing simple approach)
	if h.Store != nil {
		if _, err := h.Store.SaveSearch(req.Query, summary, articleIDs); err != nil {
			log.Printf("store save error: %v", err)
			// don't fail the request
		}
	}

	resp := models.SearchResponse{
		AISummaryAnswer:    summary,
		AIRelevantArticles: articleIDs,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
