package models

import (
	"github.com/google/uuid"
)

type ProjectMemberRole string
type ProjectMemberStatus string

const (
	RoleOwner  ProjectMemberRole = "owner"
	RoleEditor ProjectMemberRole = "editor"
	RoleViewer ProjectMemberRole = "viewer"

	StatusActive   ProjectMemberStatus = "active"
	StatusInactive ProjectMemberStatus = "inactive"
)

type CreateProjectMemberDTO struct {
	UserID    uuid.UUID         `json:"user_id"`
	ProjectID uuid.UUID         `json:"project_id"`
	Role      ProjectMemberRole `json:"role"`
}

type ProjectMemberResponseDTO struct {
	ID        uuid.UUID         `json:"id"`
	UserID    uuid.UUID         `json:"userId"`
	ProjectID uuid.UUID         `json:"projectId"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Status    string            `json:"status"`
	Role      ProjectMemberRole `json:"role"`
	JoinedAt  string            `json:"joinedAt"`
}
