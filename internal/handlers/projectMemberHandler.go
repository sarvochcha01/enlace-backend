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
