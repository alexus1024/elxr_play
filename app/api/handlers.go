package main

import (
	"fmt"
	"net/http"
)

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>elxr API</title></head>
<body>
  <h1>elxr API</h1>
  <p>A learning project — Go API running on k3s on a Raspberry Pi 5.</p>
  <p><a href="/swagger/">Swagger UI</a></p>
</body>
</html>`)
}

func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok"}`)
}
