package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/lecid-code/carvault/internal/database"
	"github.com/lecid-code/carvault/internal/handlers"
)

type AppConfig struct {
	Addr          string
	DBDSN         string
	SessionSecret string
	TemplateDir   string
}

type Application struct {
	Config AppConfig
	DB     *sqlx.DB
	Logger *slog.Logger
}

func newApplication() Application {
	addr := flag.String("addr", ":8080", "HTTP network address")
	dbDSN := flag.String("db", "./data/carvault.db", "Database DSN")
	sessionSecret := flag.String("session", "", "Session secret for authentication")
	templateDir := flag.String("templates", "./templates", "Directory for HTML templates")

	flag.Parse()

	// Validate required flags
	if *sessionSecret == "" {
		log.Fatal("Session secret must be provided. Use the -session flag.")
	}

	// Initialize DB
	dbi, err := database.New(*dbDSN)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	db := dbi.Conn()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return Application{
		Config: AppConfig{
			Addr:          *addr,
			DBDSN:         *dbDSN,
			SessionSecret: *sessionSecret,
			TemplateDir:   *templateDir,
		},

		DB:     db,
		Logger: logger,
	}
}

// main function initializes the application and starts the server.
func main() {
	app := newApplication()
	defer app.DB.Close()

	// Auth handler
	auth := &handlers.AuthHandler{
		DB:            app.DB,
		SessionSecret: app.Config.SessionSecret,
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

	app.Logger.Info("Starting server...", "address", app.Config.Addr)
	err := http.ListenAndServe(app.Config.Addr, mux)
	if err != nil {
		app.Logger.Error("Server failed", "error", err)
	}
}
