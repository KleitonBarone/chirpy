package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (cfg *apiConfig) createUserHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	type returnErr struct {
		Error string `json:"error"`
	}

	type parameters struct {
		Email string `json:"email"`
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

	if params.Email == "" {
		respBody := returnErr{
			Error: "Empty Email",
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

	user, err := cfg.dbQueries.CreateUser(req.Context(), params.Email)
	if err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	type returnVals struct {
		Id        string `json:"id"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	respBody := returnVals{
		Id:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
	data, err := json.Marshal(respBody)
	if err != nil {

		res.WriteHeader(500)
		return
	}

	res.WriteHeader(201)
	res.Write(data)
}
