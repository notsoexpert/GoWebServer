package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/notsoexpert/gowebserver/internal/api"
	"github.com/notsoexpert/gowebserver/internal/database"
)

func main() {
	var apiCfg api.APIConfig
	godotenv.Load(".env")
	apiCfg.Platform = os.Getenv("PLATFORM")
	apiCfg.Secret = os.Getenv("SECRET")
	apiCfg.PolkaKey = os.Getenv("POLKA_KEY")
	dbURL := os.Getenv("DB_URL")
	fmt.Println("Connecting to ", dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error: failed to open database")
		return
	}
	apiCfg.DBQueries = database.New(db)

	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", api.ReadinessHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.GetChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirpHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.DeleteChirpHandler)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.PolkaWebhooksHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.PostChirpsHandler)
	mux.HandleFunc("POST /api/users", apiCfg.CreateUserHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.UpdateCredentialsHandler)
	mux.HandleFunc("POST /api/login", apiCfg.LoginHandler)
	mux.HandleFunc("POST /api/refresh", apiCfg.RefreshHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.RevokeHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.CountRequestsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetRequestsHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Running server...")
	if ok := server.ListenAndServe(); ok != nil {
		fmt.Println(ok.Error())
	}
	fmt.Println("Server shutting down...")
}
