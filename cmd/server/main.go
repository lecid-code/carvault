package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/lecid-code/carvault/internal/database"
	"github.com/lecid-code/carvault/internal/handlers"
)

func main() {
	// 1. Server port configurable via environment variable
	port := os.Getenv("CARVAULT_PORT")
	if port == "" {
		port = "8080"
	}

	// 2. Database DSN configurable via environment variable
	dbDSN := os.Getenv("CARVAULT_DB_DSN")
	if dbDSN == "" {
		dbDSN = "carvault.db" // fallback, e.g. SQLite file
	}

	// 3. Log level configurable via environment variable
	logLevel := strings.ToLower(os.Getenv("CARVAULT_LOG_LEVEL"))
	if logLevel == "" {
		logLevel = "info"
	}

	// 4. Session secret configurable via environment variable
	sessionSecret := os.Getenv("CARVAULT_SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("CARVAULT_SESSION_SECRET must be set")
	}

	// Example: log level usage (simple)
	if logLevel == "debug" {
		log.Printf("[DEBUG] Using DB DSN: %s", dbDSN)
	}

	// Initialize DB
	dbi, err := database.New(dbDSN)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer dbi.Close()
	db := dbi.Conn()

	// Auth handler
	auth := &handlers.AuthHandler{
		DB:            db,
		SessionSecret: sessionSecret,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("CarVault server is running."))
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			auth.LoginForm(w, r)
		case http.MethodPost:
			auth.Login(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/logout", auth.Logout)

	log.Printf("Starting server on :%s...", port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
