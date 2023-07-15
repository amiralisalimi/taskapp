package main

import (
	"context"
	"encoding/json"
	"log"
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
	tm *TaskManager
}

func NewUserHandler(db *sqlx.DB, tm *TaskManager) *UserHandler {
	return &UserHandler{db: db, tm: tm}
}

func (uh *UserHandler) GetUserData(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request context
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the user's preferred background and theme from the database
	var background string
	err := uh.db.Get(&background, "SELECT background FROM users WHERE id = $1", userID)
	if err != nil {
		log.Println("Could not get user background")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the user's boards, containers, and tasks from the database
	var boards []Board
	err = uh.db.Select(&boards, "SELECT id FROM boards WHERE user_id = $1", userID)
	if err != nil {
		log.Printf("Could not get user board for user with id %v: %v\n", userID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var containers []Container
	var tasks []Task

	for _, board := range boards {
		// Get the containers for the board
		containerIDs, err := uh.tm.getContainersForBoard(board.ID)
		if err != nil {
			log.Println("Could not get board containers")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, containerID := range containerIDs {
			// Get the tasks for the container
			taskIDs, err := uh.tm.getTasksForContainer(containerID)
			if err != nil {
				log.Println("Could not get container tasks")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for _, taskID := range taskIDs {
				// Get the task from the database
				var task Task
				err = uh.db.Get(&task, "SELECT * FROM tasks WHERE id = $1", taskID)
				if err != nil {
					log.Println("Could not get task data")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				tasks = append(tasks, task)
			}

			// Get the container from the database
			var container Container
			err = uh.db.Get(&container, "SELECT * FROM containers WHERE id = $1", containerID)
			if err != nil {
				log.Println("Could not get container data")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			containers = append(containers, container)
		}

		// Add the container IDs to the board
		board.ContainerIDs = containerIDs
	}

	// Construct the response object
	response := map[string]interface{}{
		"boards":     boards,
		"containers": containers,
		"tasks":      tasks,
		"background": background,
	}

	// Send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *UserHandler) signupHandler(w http.ResponseWriter, r *http.Request) {
	var data SignupData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("Signup Data invalid: %v\n", data)
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

	ctx := r.Context()
	ctx = context.WithValue(ctx, "background", "img-3.jpg")
	ctx = context.WithValue(ctx, "title", data.Username+"' Board")

	s.tm.CreateBoardHandler(w, r.WithContext(ctx))

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
