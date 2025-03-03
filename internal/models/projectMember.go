package models

import (
	"github.com/google/uuid"
)

type ProjectRole string

const (
	Owner  ProjectRole = "owner"
	Editor ProjectRole = "editor"
	Viewer ProjectRole = "viewer"
)

type CreateProjectMemberDTO struct {
	UserID    uuid.UUID   `json:"user_id"`
	ProjectID uuid.UUID   `json:"project_id"`
	Role      ProjectRole `json:"role"`
}

type ProjectMemberResponseDTO struct {
	ID       uuid.UUID   `json:"id"`
	UserID   uuid.UUID   `json:"userId"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Role     ProjectRole `json:"role"`
	JoinedAt string      `json:"joinedAt"`
}
