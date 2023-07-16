package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Board struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	Title        string `json:"title"`
	Background   string `json:"background"`
	ContainerIDs []int  `json:"container_ids"`
}

type Container struct {
	ID      int    `json:"id" db:"id"`
	BoardID int    `json:"board_id" db:"board_id"`
	Title   string `json:"title" db:"title"`
	TaskIDs []int  `json:"task_ids"`
}

type Task struct {
	ID          int    `json:"id"`
	ContainerID int    `json:"container_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type TaskManager struct {
	db *sqlx.DB
}

func NewTaskManager(db *sqlx.DB) *TaskManager {
	return &TaskManager{db: db}
}

func (tm *TaskManager) GetBoardsHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID").(int)

	rows, err := tm.db.Query("SELECT * FROM boards WHERE user_id = $1", userID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	boards := []Board{}

	for rows.Next() {
		board := Board{}
		err := rows.Scan(&board.ID, &board.UserID, &board.Title, &board.Background, &board.ContainerIDs)
		if err != nil {
			log.Fatal(err)
		}
		boards = append(boards, board)
	}

	for i := range boards {
		containers, err := tm.getContainersForBoard(boards[i].ID)
		if err != nil {
			log.Fatal(err)
		}
		boards[i].ContainerIDs = containers

		/* 		for j := range containers {
			tasks, err := tm.getTasksForContainer(containers[j])
			if err != nil {
				log.Fatal(err)
			}
			containers[j].TaskIDs = tasks
		} */
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boards)
}

func (tm *TaskManager) CreateBoardHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID").(int)

	board := Board{}
	err := json.NewDecoder(r.Body).Decode(&board)
	if err != nil {
		board.Title = r.Context().Value("title").(string)
		board.Background = r.Context().Value("background").(string)
	}

	var id int
	tm.db.QueryRow("INSERT INTO boards (user_id, title, background) VALUES ($1, $2, $3) RETURNING id", userID, board.Title, board.Background).Scan(&id)

	board.ID = int(id)
	board.UserID = userID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}

func (tm *TaskManager) UpdateBoardHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID").(int)

	vars := mux.Vars(r)
	boardID := vars["id"]

	board := Board{}
	err := json.NewDecoder(r.Body).Decode(&board)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = tm.checkBoardOwnership(userID, boardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	_, err = tm.db.Exec("UPDATE boards SET title = $1, background = $2 WHERE id = $3", board.Title, board.Background, boardID)
	if err != nil {
		log.Fatal(err)
	}

	board.ID, err = strconv.Atoi(boardID)
	board.UserID = userID
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}

func (tm *TaskManager) DeleteBoardHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID").(int)

	vars := mux.Vars(r)
	boardID := vars["id"]

	err := tm.checkBoardOwnership(userID, boardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = tm.deleteContainersForBoard(boardID)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tm.db.Exec("DELETE FROM boards WHERE id = $1", boardID)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (tm *TaskManager) GetContainersHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	boardID := vars["id"]

	rows, err := tm.db.Query("SELECT * FROM containers WHERE board_id = $1", boardID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	containers := []Container{}

	for rows.Next() {
		container := Container{}
		err := rows.Scan(&container.ID, &container.BoardID, &container.Title, &container.TaskIDs)
		if err != nil {
			log.Fatal(err)
		}
		containers = append(containers, container)
	}

	for i := range containers {
		tasks, err := tm.getTasksForContainer(containers[i].ID)
		if err != nil {
			log.Fatal(err)
		}
		containers[i].TaskIDs = tasks
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(containers)
}

func (tm *TaskManager) CreateContainerHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	boardID := vars["id"]

	container := Container{}
	err := json.NewDecoder(r.Body).Decode(&container)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := tm.db.Exec("INSERT INTO containers (board_id, title) VALUES ($1, $2)", boardID, container.Title)
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	container.ID = int(id)
	container.BoardID, err = strconv.Atoi(boardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(container)
}

func (tm *TaskManager) UpdateContainerHandler(w http.ResponseWriter, r *http.Request) {

	var containerData Container
	err := json.NewDecoder(r.Body).Decode(&containerData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = tm.db.Exec("UPDATE containers SET title = $1 WHERE id = $1", containerData.Title, containerData.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (tm *TaskManager) DeleteContainerHandler(w http.ResponseWriter, r *http.Request) {

	containerID := mux.Vars(r)["id"]

	_, err := tm.db.Exec("DELETE FROM tasks WHERE container_id = $1", containerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = tm.db.Exec("DELETE FROM containers WHERE id = $!", containerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (tm *TaskManager) GetTasksHandler(w http.ResponseWriter, r *http.Request) {

	containerID := r.URL.Query().Get("container_id")

	var tasks []Task
	err := tm.db.Select(&tasks, "SELECT * FROM tasks WHERE container_id = $1", containerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []Task
	for _, task := range tasks {
		response = append(response, Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (tm *TaskManager) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var taskData Task
	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := tm.db.Exec("INSERT INTO tasks (title, description, completed, container_id) VALUES ($1, $2, $3, $4)", taskData.Title, taskData.Description, taskData.Completed, taskData.ContainerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Task{
		ID:          int(taskID),
		Title:       taskData.Title,
		Description: taskData.Description,
		Completed:   taskData.Completed,
		ContainerID: taskData.ContainerID,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (tm *TaskManager) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var taskData Task
	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = tm.db.Exec("UPDATE tasks SET title = $1, description = $2, completed = $3 WHERE id = $4", taskData.Title, taskData.Description, taskData.Completed, taskData.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Task{
		ID:          taskData.ID,
		Title:       taskData.Title,
		Description: taskData.Description,
		Completed:   taskData.Completed,
		ContainerID: taskData.ContainerID,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (tm *TaskManager) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	taskID := mux.Vars(r)["id"]

	_, err := tm.db.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (tm *TaskManager) getContainersForBoard(boardID int) ([]int, error) {

	rows, err := tm.db.Query("SELECT id FROM containers WHERE board_id = $1", boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	containerIDs := []int{}

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		containerIDs = append(containerIDs, id)
	}

	return containerIDs, nil
}

func (tm *TaskManager) getTasksForContainer(containerID int) ([]int, error) {

	rows, err := tm.db.Query("SELECT id FROM tasks WHERE container_id = $1", containerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	taskIDs := []int{}

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		taskIDs = append(taskIDs, id)
	}

	return taskIDs, nil
}

func (tm *TaskManager) deleteContainersForBoard(boardID string) error {

	rows, err := tm.db.Query("SELECT id FROM containers WHERE board_id = $1", boardID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return err
		}
		err = tm.deleteTasksForContainer(id)
		if err != nil {
			return err
		}
		_, err = tm.db.Exec("DELETE FROM containers WHERE id = $1", id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tm *TaskManager) deleteTasksForContainer(containerID int) error {
	_, err := tm.db.Exec("DELETE FROM tasks WHERE container_id = $1", containerID)
	if err != nil {
		return err
	}
	return nil
}

func (tm *TaskManager) checkBoardOwnership(userID int, boardID string) error {
	var ownerID int
	err := tm.db.QueryRow("SELECT user_id FROM boards WHERE id = $1", boardID).Scan(&ownerID)
	if err != nil {
		return err
	}
	if userID != ownerID {
		return ErrForbidden
	}
	return nil
}
