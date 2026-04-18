package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

type CounterResponse struct {
	Value int `json:"value"`
}

// @Summary Get next counter value
// @Description Returns the next sequential number from memory
// @Tags counter
// @Accept json
// @Produce json
// @Success 200 {object} CounterResponse
// @Router /counter/next [post]
func (a *App) handleCounterNext(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	a.counterMu.Lock()
	a.counter++
	currentValue := a.counter
	a.counterMu.Unlock()

	a.logger.InfoContext(ctx, "counter incremented", slog.Int("value", currentValue))

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"value":%d}`, currentValue)
}
