package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

// @Summary Stops the current instance
// @Description Shuts down the application. mode=graceful (default), panic, exit. exit accepts code=N query param.
// @Tags debug
// @Produce json
// @Param mode query string false "Stop mode: graceful, panic, exit"
// @Param code query int false "Exit code (only for mode=exit, default 1)"
// @Success 200 "Application shutting down"
// @Router /kill [post]
func (a *App) handleKill(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "graceful"
	}

	a.logger.Warn("kill endpoint called", slog.String("mode", mode))
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"shutting down","mode":%q}`, mode)

	go func() {
		time.Sleep(1 * time.Second)
		switch mode {
		case "panic":
			panic("kill endpoint: panic mode")
		case "exit":
			code := 1
			if c, err := strconv.Atoi(r.URL.Query().Get("code")); err == nil {
				code = c
			}
			os.Exit(code)
		default: // graceful
			close(a.quit)
		}
	}()
}
