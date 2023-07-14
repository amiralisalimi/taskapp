package main

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	db *sqlx.DB
}

func NewUserHandler(db *sqlx.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (s *UserHandler) signupHandler(w http.ResponseWriter, r *http.Request) {
	var data SignupData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", data.Username, data.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *UserHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	err = s.db.QueryRow("SELECT id, username, email, password FROM users WHERE username = $1", creds.Username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *UserHandler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logout logic, e.g. delete session cookies, clear authentication tokens, etc.

	w.WriteHeader(http.StatusOK)
}

func (s *UserHandler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement authentication logic, e.g. check for session cookies, validate authentication tokens, etc.

		next.ServeHTTP(w, r)
	})
}
