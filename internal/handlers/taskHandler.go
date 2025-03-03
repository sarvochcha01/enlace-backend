package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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
		log.Fatal("Invalid project ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid project ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		log.Fatal("Invalid request body: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Fatal("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var userID uuid.UUID
	if userID, err = h.taskService.CreateTask(&taskDTO, user.UID); err != nil {
		log.Fatal("Failed to create project: ", err)
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Task created successfully with ID: " + userID.String()))

}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {

}
