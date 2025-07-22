package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	schema := `CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT, password_hash TEXT);`
	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}
	// password: "password" (bcrypt hash generated for test)
	pw := "$2a$10$on2k7Zqbe2ZWhbSAo3Ck0.sGWVLl8fvXWajwBczBIDzSQ3KmmmtGu" // valid hash for "password"
	_, err = db.Exec(`INSERT INTO users (username, password_hash) VALUES (?, ?)`, "testuser", pw)
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}
	return db
}

func TestLoginForm(t *testing.T) {
	db := setupTestDB(t)
	h := &AuthHandler{DB: db, SessionSecret: "secret"}
	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()
	h.LoginForm(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestLoginSuccess(t *testing.T) {
	db := setupTestDB(t)
	h := &AuthHandler{DB: db, SessionSecret: "secret"}
	form := strings.NewReader("username=testuser&password=password")
	r := httptest.NewRequest(http.MethodPost, "/login", form)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h.Login(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("expected redirect, got %d", resp.StatusCode)
	}
	cookie := resp.Cookies()
	found := false
	for _, c := range cookie {
		if c.Name == "jwt" && c.Value != "" {
			found = true
		}
	}
	if !found {
		t.Error("jwt cookie not set on login")
	}
}

func TestLoginFail(t *testing.T) {
	db := setupTestDB(t)
	h := &AuthHandler{DB: db, SessionSecret: "secret"}
	form := strings.NewReader("username=testuser&password=wrong")
	r := httptest.NewRequest(http.MethodPost, "/login", form)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h.Login(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestLogout(t *testing.T) {
	db := setupTestDB(t)
	h := &AuthHandler{DB: db, SessionSecret: "secret"}
	r := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()
	h.Logout(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("expected redirect, got %d", resp.StatusCode)
	}
	cookie := resp.Cookies()
	found := false
	for _, c := range cookie {
		if c.Name == "jwt" && c.Value == "" && c.MaxAge == -1 {
			found = true
		}
	}
	if !found {
		t.Error("jwt cookie not cleared on logout")
	}
}
