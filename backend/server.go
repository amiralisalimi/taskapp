package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for the presence of an authorization header
		// authHeader := r.Header.Get("Authorization")
		// if authHeader == "" {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		// Set the user ID as a value in the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", 1)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

	r.Use(authMiddleware)
	r.Use(CORSMiddleware)

	r.HandleFunc("/boards", tm.GetBoardsHandler)                       //.Methods("GET")
	r.HandleFunc("/boards", tm.CreateBoardHandler)                     //.Methods("POST")
	r.HandleFunc("/boards/{id}", tm.UpdateBoardHandler)                //.Methods("PUT")
	r.HandleFunc("/boards/{id}", tm.DeleteBoardHandler)                //.Methods("DELETE")
	r.HandleFunc("/boards/{id}/containers", tm.GetContainersHandler)   //.Methods("GET")
	r.HandleFunc("/boards/{id}/containers", tm.CreateContainerHandler) //.Methods("POST")
	r.HandleFunc("/containers/{id}", tm.UpdateContainerHandler)        //.Methods("PUT")
	r.HandleFunc("/containers/{id}", tm.DeleteContainerHandler)        //.Methods("DELETE")
	r.HandleFunc("/containers/{id}/tasks", tm.GetTasksHandler)         //.Methods("GET")
	r.HandleFunc("/containers/{id}/tasks", tm.CreateTaskHandler)       //.Methods("POST")
	r.HandleFunc("/tasks/{id}", tm.UpdateTaskHandler)                  //.Methods("PUT")
	r.HandleFunc("/tasks/{id}", tm.DeleteTaskHandler)                  //.Methods("DELETE")

	r.HandleFunc("/signup", uh.signupHandler)                    //.Methods("POST")
	r.HandleFunc("/login", uh.loginHandler)                      //.Methods("POST")
	r.HandleFunc("/logout", uh.authMiddleware(uh.logoutHandler)) //.Methods("POST")
	r.HandleFunc("/user-data", uh.GetUserData)                   //.Methods("GET")

	log.Println("Server listening on port 8000")
	log.Fatal(http.ListenAndServeTLS("localhost:8000", "ssl/certificate.crt", "ssl/private.key", r))
}
