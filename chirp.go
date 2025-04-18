package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/KleitonBarone/chirpy/internal/auth"
	"github.com/KleitonBarone/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirpHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	type returnErr struct {
		Error string `json:"error"`
	}

	jwtToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Println(err)
		res.WriteHeader(401)
		return
	}
	userID, err := auth.ValidateJWT(jwtToken, cfg.jwtSecret)
	if err != nil {
		log.Println(err)
		res.WriteHeader(401)
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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

	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	type returnVals struct {
		ID        string `json:"id"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	respBody := returnVals{
		ID:        chirp.ID.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
		CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
	}
	data, err := json.Marshal(respBody)
	if err != nil {
		res.WriteHeader(500)
		return
	}

	res.WriteHeader(201)
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

func (cfg *apiConfig) getChirpsHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	chirps, err := cfg.dbQueries.GetChirps(req.Context())
	if err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	type returnVals struct {
		ID        string `json:"id"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	respBody := []returnVals{}

	for _, chirp := range chirps {
		respBody = append(respBody, returnVals{
			ID:        chirp.ID.String(),
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
			CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
			UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
		})
	}

	data, err := json.Marshal(respBody)
	if err != nil {
		res.WriteHeader(500)
		return
	}

	res.WriteHeader(200)
	res.Write(data)
}

func (cfg *apiConfig) getChirpHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	chirpID := req.PathValue("chirpID")
	chirpIdUUID, err := uuid.Parse(chirpID)
	if err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), chirpIdUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			res.WriteHeader(404)
			return
		}
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	type returnVals struct {
		ID        string `json:"id"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	respBody := returnVals{
		ID:        chirp.ID.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
		CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
	}

	data, err := json.Marshal(respBody)
	if err != nil {
		res.WriteHeader(500)
		return
	}

	res.WriteHeader(200)
	res.Write(data)
}
