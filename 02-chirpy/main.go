package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, req)
	})
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{fileserverHits: 0}

	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	//* Server ping route
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /api/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("/api/reset", apiCfg.handleResetFileHits)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handleValidateChirp)

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

func (cfg *apiConfig) handleValidateChirp(w http.ResponseWriter, req *http.Request) {
	// Structs
	type parameters struct {
		Body string `json:"body"`
	}
	type errResp struct {
		Error string `json:"error"`
	}
	type succResp struct {
		Valid bool `json:"valid"`
	}
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(500)
		respBody := errResp{
			Error: "Something went wrong",
		}
		dat, _ := json.Marshal(respBody)
		w.Write(dat)
		return
	}
	if len(params.Body) > 140 {
		w.WriteHeader(400)
		respBody := errResp{
			Error: "Chirp is too long",
		}
		dat, _ := json.Marshal(respBody)
		w.Write(dat)
		return
	}

	w.WriteHeader(200)
	respBody := succResp{
		Valid: true,
	}
	dat, _ := json.Marshal(respBody)
	w.Write(dat)
}

func (cfg *apiConfig) handleAdminMetrics(w http.ResponseWriter, req *http.Request) {
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
