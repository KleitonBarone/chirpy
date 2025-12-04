package handler

import (
	"fmt"
	"net/http"

	"github.com/KleitonBarone/chirpy/internal/config"
)

// MiddlewareMetricsInc increments the fileserver hits counter.
func MiddlewareMetricsInc(cfg *config.ApiConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}

// FileServerHitsHandler returns the number of fileserver hits.
func FileServerHitsHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
	`, cfg.FileserverHits.Load())))
	}
}

// FileServerHitsResetHandler resets the fileserver hits counter.
func FileServerHitsResetHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if cfg.Platform != "dev" {
			res.WriteHeader(http.StatusForbidden)
			return
		}

		err := cfg.DbQueries.RemoveAllUsers(req.Context())
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		cfg.FileserverHits.Store(0)
		res.WriteHeader(http.StatusOK)
	}
}
