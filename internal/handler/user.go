package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/KleitonBarone/chirpy/internal/auth"
	"github.com/KleitonBarone/chirpy/internal/config"
	"github.com/KleitonBarone/chirpy/internal/database"
)

// CreateUserHandler handles user creation.
func CreateUserHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		type returnErr struct {
			Error string `json:"error"`
		}

		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
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

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			log.Println(err)
			res.WriteHeader(500)
			return
		}

		user, err := cfg.DbQueries.CreateUser(req.Context(), database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hashedPassword,
		})
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
}

// LoginHandler handles user login.
func LoginHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		type returnErr struct {
			Error string `json:"error"`
		}

		type parameters struct {
			Email            string `json:"email"`
			Password         string `json:"password"`
			ExpiresInSeconds *int   `json:"expires_in_seconds"`
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

		user, err := cfg.DbQueries.GetUserByEmail(req.Context(), params.Email)
		if err != nil {
			log.Println(err)
			res.WriteHeader(401)
			return
		}

		if auth.CheckPasswordHash(params.Password, user.HashedPassword) != nil {
			res.WriteHeader(401)
			return
		}

		if params.ExpiresInSeconds == nil {
			params.ExpiresInSeconds = new(int)
		}

		if *params.ExpiresInSeconds < 1 || *params.ExpiresInSeconds > 3600 {
			*params.ExpiresInSeconds = 3600
		}

		token, err := auth.MakeJWT(user.ID, cfg.JwtSecret, time.Duration(*params.ExpiresInSeconds)*time.Second)
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
			Token     string `json:"token"`
		}

		respBody := returnVals{
			Id:        user.ID.String(),
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
			Token:     token,
		}
		data, err := json.Marshal(respBody)
		if err != nil {

			res.WriteHeader(500)
			return
		}

		res.WriteHeader(200)
		res.Write(data)
	}
}
