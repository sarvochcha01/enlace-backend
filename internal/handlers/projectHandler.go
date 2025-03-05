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

type ProjectHandler struct {
	projectService services.ProjectService
}

func NewProjectHandler(projectService services.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var projectDTO models.CreateProjectDTO

	if err := json.NewDecoder(r.Body).Decode(&projectDTO); err != nil {
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

	if err := h.projectService.CreateProject(&projectDTO, user.UID); err != nil {
		log.Println("Failed to create project: ", err)
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Project created successfully"))
}

func (h *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {

	projectID := chi.URLParam(r, "projectID")

	parsedProjectID, err := uuid.Parse(projectID)

	if err != nil {
		log.Println("Invalid project ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid project ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	project, err := h.projectService.GetProjectByID(parsedProjectID)

	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)

}

func (h *ProjectHandler) GetAllProjectsForUser(w http.ResponseWriter, r *http.Request) {
	var user *auth.Token
	var err error
	user, err = middlewares.GetFirebaseUser(r)

	if err != nil {
		log.Println("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var projectResponseDTO []models.ProjectResponseDTO

	if projectResponseDTO, err = h.projectService.GetAllProjectsForUser(user.UID); err != nil {
		log.Println("Failed to get Projects: ", err)
		http.Error(w, "Failed to get Projects", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(projectResponseDTO); err != nil {
		log.Println("Failed to encode response: ", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
