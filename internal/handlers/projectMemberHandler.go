package handlers

import (
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
	createProjectMemberDTO.Role = models.Viewer

	projectID := chi.URLParam(r, "projectID")

	parsedProjectID, err := uuid.Parse(projectID)

	if err != nil {
		log.Fatal("Invalid project ID (must be a valid UUID): ", err)
		http.Error(w, "Invalid project ID (must be a valid UUID)", http.StatusBadRequest)
		return
	}
	createProjectMemberDTO.ProjectID = parsedProjectID

	var user *auth.Token

	user, err = middlewares.GetFirebaseUser(r)

	if err != nil {
		log.Fatal("Unauthorized: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err = h.projectMemberService.CreateProjectMember(&createProjectMemberDTO, user.UID); err != nil {
		log.Fatal("Failed to join project: ", err)
		http.Error(w, "Failed to join project", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Project joined successfully"))
}
