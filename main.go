package main

import (
	"log"
	"motadataAssignment/ai"
	"motadataAssignment/api"
	"motadataAssignment/kb"
	"motadataAssignment/store"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	k := kb.New()
	aiClient := ai.NewSimulated()
	st := store.NewInMemoryStore()

	h := api.NewHandler(k, aiClient, st)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("starting server on %s", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
