package admin

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/notsoexpert/gowebserver/internal/database"
)

type APIConfig struct {
	DBQueries      *database.Queries
	fileserverHits atomic.Int32
}

func (cfg *APIConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *APIConfig) CountRequestsHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "text/html; charset=utf-8")
	response.Write(fmt.Appendf([]byte{}, `
	<html>
		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, cfg.fileserverHits.Load()))
}

func (cfg *APIConfig) ResetRequestsHandler(response http.ResponseWriter, request *http.Request) {
	cfg.fileserverHits.Store(0)
}
