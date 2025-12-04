package config

import (
	"sync/atomic"

	"github.com/KleitonBarone/chirpy/internal/database"
)

// ApiConfig holds the application configuration and shared state.
type ApiConfig struct {
	FileserverHits atomic.Int32
	DbQueries      *database.Queries
	Platform       string
	JwtSecret      string
}

// NewApiConfig creates a new ApiConfig instance.
func NewApiConfig(dbQueries *database.Queries, platform, jwtSecret string) *ApiConfig {
	return &ApiConfig{
		FileserverHits: atomic.Int32{},
		DbQueries:      dbQueries,
		Platform:       platform,
		JwtSecret:      jwtSecret,
	}
}
