package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const secretKey = "nXkAMA8Mbk4eRzQnwH17beh3YarW3cEK4XefyeiO6hM="

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("auth header not set")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusBadRequest)
			return
		}

		token := authParts[1]

		userID, err := getUserIDFromToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserIDFromToken(token string) (int, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, err
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user ID not found in token")
	}

	return int(userID), nil
}

func main() {
	db, err := sqlx.Open("postgres", "postgres://postgres:postgres@localhost:5432/taskapp?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tm := NewTaskManager(db)
	uh := NewUserHandler(db, tm)

	r := mux.NewRouter()

	// r.Use(authMiddleware)
	r.Use(CORSMiddleware)

	r.HandleFunc("/boards", authMiddleware(tm.GetBoardsHandler))                       //.Methods("GET")
	r.HandleFunc("/boards", authMiddleware(tm.CreateBoardHandler))                     //.Methods("POST")
	r.HandleFunc("/boards/{id}", authMiddleware(tm.UpdateBoardHandler))                //.Methods("PUT")
	r.HandleFunc("/boards/{id}", authMiddleware(tm.DeleteBoardHandler))                //.Methods("DELETE")
	r.HandleFunc("/boards/{id}/containers", authMiddleware(tm.GetContainersHandler))   //.Methods("GET")
	r.HandleFunc("/boards/{id}/containers", authMiddleware(tm.CreateContainerHandler)) //.Methods("POST")
	r.HandleFunc("/containers/{id}", authMiddleware(tm.UpdateContainerHandler))        //.Methods("PUT")
	r.HandleFunc("/containers/{id}", authMiddleware(tm.DeleteContainerHandler))        //.Methods("DELETE")
	r.HandleFunc("/containers/{id}/tasks", authMiddleware(tm.GetTasksHandler))         //.Methods("GET")
	r.HandleFunc("/containers/{id}/tasks", authMiddleware(tm.CreateTaskHandler))       //.Methods("POST")
	r.HandleFunc("/tasks/{id}", authMiddleware(tm.UpdateTaskHandler))                  //.Methods("PUT")
	r.HandleFunc("/tasks/{id}", authMiddleware(tm.DeleteTaskHandler))                  //.Methods("DELETE")

	r.HandleFunc("/signup", uh.signupHandler)                            //.Methods("POST")
	r.HandleFunc("/login", uh.loginHandler)                              //.Methods("POST")
	r.HandleFunc("/logout", uh.logoutHandler)                            //.Methods("POST")
	r.HandleFunc("/user-data", authMiddleware(uh.GetUserData))           //.Methods("GET")
	r.HandleFunc("/update-user-data", authMiddleware(uh.UpdateUserData)) //.Methods("POST")

	log.Println("Server listening on port 8000")
	log.Fatal(http.ListenAndServeTLS("localhost:8000", "ssl/certificate.crt", "ssl/private.key", r))
}
