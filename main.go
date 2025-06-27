package main

import (
	"database/sql"
	"fmt"
	"net/http"
<<<<<<< HEAD
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/notsoexpert/gowebserver/internal/admin"
	"github.com/notsoexpert/gowebserver/internal/api"
	"github.com/notsoexpert/gowebserver/internal/database"
)

=======
	"sync/atomic"
	"encoding/json"
	"strings"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) CountRequestsHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(200)
	response.Write(fmt.Appendf([]byte{}, `
	<html>
		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, cfg.fileserverHits.Load()))
}

func (cfg *apiConfig) ResetRequestsHandler(response http.ResponseWriter, request *http.Request) {
	cfg.fileserverHits.Store(0)
}

func ReadinessHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(200)
	response.Write([]byte("OK"))
}

type ChirpRequest struct {
    Body string `json:"body"`
}

type ChirpValidation struct {
	CleanedBody string `json:"cleaned_body,omitempty"`
	Error string `json:"error,omitempty"`
	Valid bool `json:"valid"`
}

func cleanResponseBody(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		switch strings.ToLower(word) {
			case "kerfuffle", "sharbert", "fornax":
				words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func respondWithError(response http.ResponseWriter, code int, msg string) {
	respBody := ChirpValidation{
		Error: msg,
		Valid: false,
	}
	data, encErr := json.Marshal(respBody)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response") // funny infinite error recursion
		return
	}
	response.WriteHeader(code)
	response.Write(data)
}

func respondWithJSON(response http.ResponseWriter, code int, payload interface{}) {
	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(code)
	response.Write(payload.([]byte))
}

func ValidateChirpHandler(response http.ResponseWriter, request *http.Request) {
    decoder := json.NewDecoder(request.Body)
    params := ChirpRequest{}
    err := decoder.Decode(&params)
    if err != nil {
		respondWithError(response, 400, "Something went wrong")
		return
    }

	if len(params.Body) > 140 {
		respondWithError(response, 400, "Chirp is too long")
		return
	}

	respBody := ChirpValidation{
		CleanedBody: cleanResponseBody(params.Body),
		Valid: true,
	}
	data, encErr := json.Marshal(respBody)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response")
		return
	}
	respondWithJSON(response, 200, data)
}

>>>>>>> 37e0c7198bf419c6668ec4b9b138a33d645acc4e
func main() {
	var apiCfg admin.APIConfig
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error: failed to open database")
		return
	}
	apiCfg.DBQueries = database.New(db)

	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
<<<<<<< HEAD
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", api.ReadinessHandler)
	mux.HandleFunc("POST /api/validate_chirp", api.ValidateChirpHandler)
	mux.HandleFunc("POST /api/users", api.CreateUserHandler)
=======
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", ReadinessHandler)
	mux.HandleFunc("POST /api/validate_chirp", ValidateChirpHandler)
>>>>>>> 37e0c7198bf419c6668ec4b9b138a33d645acc4e
	mux.HandleFunc("GET /admin/metrics", apiCfg.CountRequestsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetRequestsHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Running server...")
	ok := server.ListenAndServe()
	if ok != nil {
		fmt.Println(ok.Error())
	}
	fmt.Println("Server shutting down...")
}
