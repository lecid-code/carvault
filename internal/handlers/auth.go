package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB            *sqlx.DB
	SessionSecret string
}

func (h *AuthHandler) LoginForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := loadLoginTemplate()
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.LoginFormWithError(w, r, "Invalid form data")
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	var hash string
	var id int
	err := h.DB.QueryRowx("SELECT id, password_hash FROM users WHERE username = ?", username).Scan(&id, &hash)
	if err == sql.ErrNoRows || err != nil {
		h.LoginFormWithError(w, r, "Invalid username or password")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		h.LoginFormWithError(w, r, "Invalid username or password")
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	jwtString, err := token.SignedString([]byte(h.SessionSecret))
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    jwtString,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to destination if provided
	dest := r.URL.Query().Get("next")
	if dest == "" {
		dest = "/"
	}
	http.Redirect(w, r, dest, http.StatusSeeOther)
}

func (h *AuthHandler) LoginFormWithError(w http.ResponseWriter, r *http.Request, errMsg string) {
	tmpl, err := loadLoginTemplate()
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]string{"Error": errMsg})
}

func loadLoginTemplate() (*template.Template, error) {
	path := filepath.Join("templates", "login.html")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return template.New("login").ParseFiles(path)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "jwt",
		Value: "",
	})
}
