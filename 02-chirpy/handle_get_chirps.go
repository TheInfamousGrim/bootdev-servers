package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *ApiConfig) handleGetChirp(w http.ResponseWriter, req *http.Request) {
	// Structs
	type chirpType struct {
		Id   int    `json:"id"`
		Body string `json:"body"`
	}
	type succResp []chirpType
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var resp succResp
	for _, chirp := range savedChirps {
		resp = append(resp, chirpType{
			Id:   chirp.id,
			Body: chirp.body,
		})
	}
	dat, _ := json.Marshal(resp)
	w.Write(dat)
}
