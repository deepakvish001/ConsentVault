package main

import (
    "encoding/json"
    "log/slog"
    "net/http"
    "os"
    "time"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
    })
    server := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
    slog.Info("ConsentVault API listening", "address", server.Addr)
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed { slog.Error("server stopped", "error", err); os.Exit(1) }
}
