package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/notsoexpert/gowebserver/internal/admin"
	"github.com/notsoexpert/gowebserver/internal/api"
	"github.com/notsoexpert/gowebserver/internal/database"
)

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
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", api.ReadinessHandler)
	mux.HandleFunc("POST /api/validate_chirp", api.ValidateChirpHandler)
	mux.HandleFunc("POST /api/users", api.CreateUserHandler)
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
