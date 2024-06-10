package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
    fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func (w http.ResponseWriter, req *http.Request) {
        cfg.fileserverHits++
        next.ServeHTTP(w, req)
    })
}

func main() {
	const filepathRoot = "."
	const port = "8080"
    apiCfg := apiConfig{ fileserverHits: 0 }

	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))


    //* Server ping route
    mux.HandleFunc("GET /api/healthz", handleReadiness)
    mux.HandleFunc("GET /api/metrics", apiCfg.handleMetrics)
    mux.HandleFunc("/api/reset", apiCfg.handleResetFileHits)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handleReadiness(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) handleResetFileHits(w http.ResponseWriter, req *http.Request) {
    // Reset the file server hits
    cfg.fileserverHits = 0
    
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Hits %d", cfg.fileserverHits)))
}