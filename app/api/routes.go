package main

import (
	"fmt"
	"net/http"
	"time"

	_ "app/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (a *App) setupRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/counter/next", a.handleCounterNext)
	mux.HandleFunc("GET /api/counter/next", a.handleCounterNext)
	mux.HandleFunc("POST /api/llm/chat", a.handleLLMChat)
	mux.HandleFunc("POST /api/kill", a.handleKill)
	mux.HandleFunc("GET /api/events", a.handleEvents)

	mux.HandleFunc("GET /health", a.handleHealth)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("/", a.handleIndex)

	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", a.config.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
