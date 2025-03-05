package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/middlewares"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/services"
)

type TaskHandler struct {
	taskService services.TaskService
}

func NewTaskHandler(taskService services.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO models.CreateTaskDTO
	var err error

	projectID := chi.URLParam(r, "projectID")

	taskDTO.ProjectID, err = uuid.Parse(projectID)
	if err != nil {
		log.Println("Invalid project ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid project ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		log.Println("Invalid request body: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if _, err = h.taskService.CreateTask(&taskDTO, user.UID); err != nil {
		log.Println("Failed to create task: ", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Task created successfully "))

}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")

	parsedTaskID, err := uuid.Parse(taskID)

	if err != nil {
		log.Println("Invalid task ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid task ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.GetTaskByID(parsedTaskID)

	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) EditTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")

	parsedTaskID, err := uuid.Parse(taskID)

	if err != nil {
		log.Println("Invalid task ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid task ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	projectID := chi.URLParam(r, "projectID")

	parsedProjectID, err := uuid.Parse(projectID)

	if err != nil {
		log.Println("Invalid task ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid task ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	var user *auth.Token
	user, err = middlewares.GetFirebaseUser(r)

	if err != nil {
		log.Println("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updateTaskDTO models.UpdateTaskDTO

	if err := json.NewDecoder(r.Body).Decode(&updateTaskDTO); err != nil {
		log.Println("Invalid request body: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.taskService.EditTask(parsedTaskID, parsedProjectID, user.UID, &updateTaskDTO); err != nil {
		log.Println("Failed to update task: ", err)
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Task updated successfully "))
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	var deleteTaskDTO models.DeleteTaskDTO

	parsedTaskID, err := uuid.Parse(taskID)

	if err != nil {
		log.Println("Invalid task ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid task ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}
	deleteTaskDTO.TaskID = parsedTaskID

	projectID := chi.URLParam(r, "projectID")

	parsedProjectID, err := uuid.Parse(projectID)

	if err != nil {
		log.Println("Invalid task ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid task ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}
	deleteTaskDTO.ProjectID = parsedProjectID

	var user *auth.Token
	user, err = middlewares.GetFirebaseUser(r)

	if err != nil {
		log.Println("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	deleteTaskDTO.FirebaseUID = user.UID

	if err := h.taskService.DeleteTask(&deleteTaskDTO); err != nil {
		log.Println("Failed to delete task:", err)
		http.Error(w, "Failed to delete task", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task Deleted"))
}
