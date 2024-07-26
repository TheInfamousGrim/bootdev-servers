package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type chirp struct {
	id   int
	body string
}

var savedChirps = []chirp{}

func (cfg *ApiConfig) HandleCreateChirps(w http.ResponseWriter, req *http.Request) {
	// Structs
	type parameters struct {
		Body string `json:"body"`
	}
	type errResp struct {
		Error string `json:"error"`
	}
	type succResp struct {
		Id   int    `json:"id"`
		Body string `json:"body"`
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

	// Validate chrip length
	if len(params.Body) > 140 {
		w.WriteHeader(400)
		respBody := errResp{
			Error: "Chirp is too long",
		}
		dat, _ := json.Marshal(respBody)
		w.Write(dat)
		return
	}

	// Validate profane words
	splitWords := strings.Split(params.Body, " ")
	astxWord := "****"
	for idx, word := range splitWords {
		lowerCaseWord := strings.ToLower(word)
		switch {
		case lowerCaseWord == "kerfuffle":
			splitWords[idx] = astxWord
			continue
		case lowerCaseWord == "sharbert":
			splitWords[idx] = astxWord
			continue
		case lowerCaseWord == "fornax":
			splitWords[idx] = astxWord
			continue
		default:
			continue
		}
	}
	cleanedWords := strings.Join(splitWords, " ")

	// Save the chirps in memory
	newId := len(savedChirps) + 1
	savedChirps = append(savedChirps, chirp{
		id:   newId,
		body: cleanedWords,
	})

	w.WriteHeader(http.StatusCreated)
	respBody := succResp{
		Id:   newId,
		Body: cleanedWords,
	}
	dat, _ := json.Marshal(respBody)
	w.Write(dat)
}
