package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type LLMRequest struct {
	Prompt string `json:"prompt"`
}

// @Summary Fake LLM endpoint
// @Description Takes a prompt and streams response tokens via Server-Sent Events
// @Tags llm
// @Accept json
// @Produce text/event-stream
// @Param request body LLMRequest true "Prompt text"
// @Success 200 "Event stream"
// @Router /llm/chat [post]
func (a *App) handleLLMChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req LLMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.WarnContext(ctx, "failed to parse request", slog.String("error", err.Error()))
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	a.logger.InfoContext(ctx, "llm request received", slog.Int("prompt_length", len(req.Prompt)))

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		a.logger.ErrorContext(ctx, "streaming not supported")
		http.Error(w, `{"error":"streaming not supported"}`, http.StatusInternalServerError)
		return
	}

	tokens := []string{
		"This", " is", " a", " fake", " LLM", " response", " to", " your", " prompt", ".", " It", " demonstrates",
		" Server", "-", "Sent", " Events", " streaming", ".", " Each", " token", " is", " sent", " separately", ".",
	}

	for _, token := range tokens {
		select {
		case <-ctx.Done():
			a.logger.InfoContext(ctx, "streaming cancelled")
			return
		default:
		}

		fmt.Fprintf(w, "data: %s\n\n", token)
		flusher.Flush()
		time.Sleep(50 * time.Millisecond)
	}

	a.logger.InfoContext(ctx, "llm response completed", slog.Int("token_count", len(tokens)))
}
