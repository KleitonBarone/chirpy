package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) validateChirpHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	type returnErr struct {
		Error string `json:"error"`
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respBody := returnErr{
			Error: "Invalid JSON",
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			res.WriteHeader(500)
			return
		}
		res.WriteHeader(500)
		res.Write(data)
		return
	}

	if params.Body == "" {
		respBody := returnErr{
			Error: "Empty Chirp",
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			res.WriteHeader(500)
			return
		}
		res.WriteHeader(400)
		res.Write(data)
		return
	}

	if len([]rune(params.Body)) > 141 {
		respBody := returnErr{
			Error: "Chirp is too long",
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			res.WriteHeader(500)
			return
		}
		res.WriteHeader(400)
		res.Write(data)
		return
	}

	blockedWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanedBody := getCleanedBody(params.Body, blockedWords)

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	respBody := returnVals{
		CleanedBody: cleanedBody,
	}
	data, err := json.Marshal(respBody)
	if err != nil {

		res.WriteHeader(500)
		return
	}

	res.WriteHeader(200)
	res.Write(data)
}

func getCleanedBody(body string, badWords []string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		for _, badWord := range badWords {
			if strings.Contains(loweredWord, badWord) {
				words[i] = "****"
			}
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
