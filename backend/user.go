package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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

type LoginResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
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
	rows, err := uh.db.Query("SELECT id, user_id, title, background FROM boards WHERE user_id = $1", userID)
	if err != nil {
		log.Printf("Could not get user board for user with id %v: %v\n", userID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		board := Board{}
		rows.Scan(&board.ID, &board.UserID, &board.Title, &board.Background)
		boards = append(boards, board)
	}
	rows.Close()

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
			err = uh.db.Get(&container, "SELECT id, board_id, title FROM containers WHERE id = $1", containerID)
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

func (s *UserHandler) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the updated data.
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}

	// log.Println(data)
	// Get the board ID from the updated data.
	boardID, ok := data["boardId"].(float64)
	if !ok {
		http.Error(w, "invalid board ID", http.StatusBadRequest)
		return
	}

	// Get the containers and tasks from the updated data.
	containers, ok := data["containers"].([]interface{})
	if !ok {
		http.Error(w, "invalid containers", http.StatusBadRequest)
		return
	}

	tasks, ok := data["tasks"].([]interface{})
	if !ok {
		http.Error(w, "invalid tasks", http.StatusBadRequest)
		return
	}

	err = s.updateContainers(int(boardID), containers)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update containers: %v", err), http.StatusInternalServerError)
		return
	}

	err = s.updateTasks(int(boardID), tasks)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update tasks: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *UserHandler) updateContainers(boardID int, containers []interface{}) error {
	s.db.Exec("DELETE FROM containers WHERE board_id = $1", boardID)
	for _, c := range containers {
		container, ok := c.(map[string]interface{})
		if !ok {
			return errors.New("invalid container")
		}

		log.Println(container)

		id := int(container["id"].(float64))
		title := container["title"].(string)

		_, err := s.db.Exec(`
			INSERT INTO containers (id, board_id, title)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO
			UPDATE SET title=EXCLUDED.title
		`, id, boardID, title)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UserHandler) updateTasks(boardID int, tasks []interface{}) error {
	s.db.Exec(`
		DELETE FROM tasks
		WHERE container_id IN (
			SELECT id
			FROM containers
			WHERE board_id = $1
		)
	`, boardID)
	for _, t := range tasks {
		task, ok := t.(map[string]interface{})
		if !ok {
			return errors.New("invalid task")
		}

		id := int(task["id"].(float64))
		containerID := int(task["container_id"].(float64))
		title, _ := task["title"].(string)
		description, _ := task["description"].(string)

		_, err := s.db.Exec(`
				INSERT INTO tasks (id, container_id, title, description)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (id) DO
				UPDATE SET container_id=EXCLUDED.container_id, title=EXCLUDED.title, description=EXCLUDED.description
			`, id, containerID, title, description)
		if err != nil {
			return err
		}
	}

	return nil
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

	user := User{
		Username: data.Username,
		Email:    data.Email,
		Password: string(hashedPassword),
	}

	err = s.db.QueryRow("INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id", data.Username, data.Email, hashedPassword).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, "userID", user.ID)
	ctx = context.WithValue(ctx, "background", "img-3.jpg")
	ctx = context.WithValue(ctx, "title", data.Username+"'s Board")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours.
		},
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("Failed to create token: %v", err)
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    tokenString,
	}

	s.db.QueryRow("INSERT INTO boards (user_id, title, background) VALUES ($1, $2, $3) RETURNING id", user.ID, ctx.Value("title"), ctx.Value("background"))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours.
		},
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
