package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "postgres://postgres:postgres@localhost:5432/taskman?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tm := NewTaskManager(db)
	uh := NewUserHandler(db)

	r := mux.NewRouter()

	r.Use(authMiddleware)

	r.HandleFunc("/boards", tm.GetBoardsHandler).Methods("GET")
	r.HandleFunc("/boards", tm.CreateBoardHandler).Methods("POST")
	r.HandleFunc("/boards/{id}", tm.UpdateBoardHandler).Methods("PUT")
	r.HandleFunc("/boards/{id}", tm.DeleteBoardHandler).Methods("DELETE")
	r.HandleFunc("/boards/{id}/containers", tm.GetContainersHandler).Methods("GET")
	r.HandleFunc("/boards/{id}/containers", tm.CreateContainerHandler).Methods("POST")
	r.HandleFunc("/containers/{id}", tm.UpdateContainerHandler).Methods("PUT")
	r.HandleFunc("/containers/{id}", tm.DeleteContainerHandler).Methods("DELETE")
	r.HandleFunc("/containers/{id}/tasks", tm.GetTasksHandler).Methods("GET")
	r.HandleFunc("/containers/{id}/tasks", tm.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", tm.UpdateTaskHandler).Methods("PUT")
	r.HandleFunc("/tasks/{id}", tm.DeleteTaskHandler).Methods("DELETE")

	r.HandleFunc("/signup", uh.signupHandler).Methods("POST")
	r.HandleFunc("/login", uh.loginHandler).Methods("POST")
	r.HandleFunc("/logout", uh.authMiddleware(uh.logoutHandler)).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
