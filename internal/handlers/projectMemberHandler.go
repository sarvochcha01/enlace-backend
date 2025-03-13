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

type ProjectMemberHandler struct {
	projectMemberService services.ProjectMemberService
}

func NewProjectMemberHandler(pms services.ProjectMemberService) *ProjectMemberHandler {
	return &ProjectMemberHandler{projectMemberService: pms}
}

func (h *ProjectMemberHandler) CreateProjectMember(w http.ResponseWriter, r *http.Request) {
	var err error
	var createProjectMemberDTO models.CreateProjectMemberDTO
	createProjectMemberDTO.Role = models.RoleViewer

	projectID := chi.URLParam(r, "projectID")

	parsedProjectID, err := uuid.Parse(projectID)

	if err != nil {
		log.Println("Invalid project ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid project ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}
	createProjectMemberDTO.ProjectID = parsedProjectID

	var user *auth.Token

	user, err = middlewares.GetFirebaseUser(r)

	if err != nil {
		log.Println("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err = h.projectMemberService.CreateProjectMember(&createProjectMemberDTO, user.UID); err != nil {
		log.Println("Failed to join project: ", err)
		http.Error(w, "Failed to join project", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Project joined successfully"))
}

func (h *ProjectMemberHandler) GetProjectMemberID(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")

	parsedProjectID, err := uuid.Parse(projectID)

	if err != nil {
		log.Println("Invalid project ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid project ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	var user *auth.Token

	user, err = middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var projectMemberResponse struct {
		ID uuid.UUID `json:"id"`
	}

	if projectMemberResponse.ID, err = h.projectMemberService.GetProjectMemberIDByFirebaseUID(user.UID, parsedProjectID); err != nil {
		log.Println("Failed to get project member: ", err)
		http.Error(w, "Failed to get project member", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(projectMemberResponse)
}

func (h *ProjectMemberHandler) UpdateProjectMember(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")

	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		log.Println("Invalid project id:", err)
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	projectMemberID := chi.URLParam(r, "projectMemberID")
	parsedProjectMemberID, err := uuid.Parse(projectMemberID)
	if err != nil {
		log.Println("Invalid member id:", err)
		http.Error(w, "invalid member id", http.StatusBadRequest)
		return
	}

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized:", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var updateProjectMemberDTO models.UpdateProjectMemberDTO
	if err = json.NewDecoder(r.Body).Decode(&updateProjectMemberDTO); err != nil {
		log.Println("Invalid request body: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateProjectMemberDTO.ID = parsedProjectMemberID
	updateProjectMemberDTO.ProjectID = parsedProjectID

	if err = h.projectMemberService.UpdateProjectMemberRole(user.UID, &updateProjectMemberDTO); err != nil {
		log.Println("Failed to update project member role:", err)
		http.Error(w, "failed to update project member role:", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Project member role updated successfully"))

}
