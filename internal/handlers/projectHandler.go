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

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectID := chi.URLParam(r, "projectID")

	parsedProjectID, err := uuid.Parse(projectID)

	if err != nil {
		log.Println("Invalid project ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid project ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	project, err := h.projectService.GetProjectByID(parsedProjectID, user.UID)

	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)

}

func (h *ProjectHandler) GetProjectName(w http.ResponseWriter, r *http.Request) {
	_, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectID := chi.URLParam(r, "projectID")

	parsedProjectID, err := uuid.Parse(projectID)

	if err != nil {
		log.Println("Invalid project ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid project ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	var projectNameResponse struct {
		Name string `json:"projectName"`
	}
	projectNameResponse.Name, err = h.projectService.GetProjectName(parsedProjectID)

	if err != nil {
		log.Println("Failed to get project name: ", err)
		http.Error(w, "Failed to get project", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projectNameResponse)
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

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {

	projectID := chi.URLParam(r, "projectID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		log.Println("Invalid project ID:", err)
		http.Error(w, "Invalid projedt ID", http.StatusBadRequest)
		return
	}

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updateProjectDTO models.EditProjectDTO
	if err = json.NewDecoder(r.Body).Decode(&updateProjectDTO); err != nil {
		log.Println("Invalid request body: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err = h.projectService.EditProject(user.UID, parsedProjectID, &updateProjectDTO); err != nil {
		log.Println("Failed to create project: ", err)
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Project updated successfully"))

}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		log.Println("Invalid project id:", err)
		http.Error(w, "Invalid project id", http.StatusBadRequest)
		return
	}

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.projectService.DeleteProject(user.UID, parsedProjectID)
	if err != nil {
		log.Println("Failed to delete project:", err)
		http.Error(w, "failed to delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Project Deleted"))
}

func (h *ProjectHandler) JoinProject(w http.ResponseWriter, r *http.Request) {

	projectID := chi.URLParam(r, "projectID")
	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		log.Println("Invalid project id:", err)
		http.Error(w, "Invalid project id", http.StatusBadRequest)
		return
	}

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.projectService.JoinProject(parsedProjectID, user.UID)
	if err != nil {
		log.Println("Failed to join project:", err)
		http.Error(w, "Failed to join project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Project joined successfully"))

}

func (h *ProjectHandler) LeaveProject(w http.ResponseWriter, r *http.Request) {

	projectID := chi.URLParam(r, "projectID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		log.Println("Invalid project ID:", err)
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.projectService.LeaveProject(parsedProjectID, user.UID)
	if err != nil {
		log.Println("Failed to leave project: ", err)
		http.Error(w, "Failed to leave project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Project left successfully"))

}
