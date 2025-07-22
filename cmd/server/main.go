package main

import (
  "log"
  "net/http"
  "os"
  "strings"
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

	   mux := http.NewServeMux()
	   mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			   w.Write([]byte("CarVault server is running."))
	   })

	   log.Printf("Starting server on :%s...", port)
	   err := http.ListenAndServe(":"+port, mux)
	   if err != nil {
			   log.Fatalf("Server failed: %v", err)
	   }
}
