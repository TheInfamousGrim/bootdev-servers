package main

import (
	"fmt"
	"log"
	"net/http"
)

type ApiConfig struct {
	fileserverHits int
}

func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, req)
	})
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := ApiConfig{fileserverHits: 0}

	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	//* Server ping route
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /api/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("/api/reset", apiCfg.handleResetFileHits)
	//* Chirps route
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandleCreateChirps)

	//* Auth Routes
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleAdminMetrics)

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

func (cfg *ApiConfig) handleMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *ApiConfig) handleResetFileHits(w http.ResponseWriter, req *http.Request) {
	// Reset the file server hits
	cfg.fileserverHits = 0

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits %d", cfg.fileserverHits)))
}

func (cfg *ApiConfig) handleAdminMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
    <html>
    
    <body>
        <h1>Welcome, Chirpy Admin</h1> 
        <p>Chirpy has been visited %d times!</p>
    </body>    

    </html> 
    `, cfg.fileserverHits)))
}
