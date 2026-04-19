package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

// @Summary Stream live events
// @Description Server-Sent Events stream of all app events (counter, requests, etc). Subscribes to NATS events.> subject.
// @Tags events
// @Produce text/event-stream
// @Param timeout query int false "Stream timeout in seconds (default 30)"
// @Success 200 "Event stream"
// @Router /events [get]
func (a *App) handleEvents(w http.ResponseWriter, r *http.Request) {
	timeoutParam := r.URL.Query().Get("timeout")
	timeoutSec := 30

	if timeoutParam != "" {
		if t, err := strconv.Atoi(timeoutParam); err == nil {
			timeoutSec = t
		}
	}

	events := make(chan *nats.Msg, 64)
	sub, err := a.natsConn.ChanSubscribe("events.>", events)
	if err != nil {
		http.Error(w, "Failed to subscribe to events", http.StatusInternalServerError)
		return
	}
	defer sub.Unsubscribe()

	w.WriteHeader(http.StatusOK)
	for {
		select {
		case msg, ok := <-events:
			if !ok {
				fmt.Fprintf(w, "END: No more events expected\n")
				return
			}
			w.Write(msg.Data)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-time.After(time.Duration(timeoutSec) * time.Second):
			return
		}
	}
}
