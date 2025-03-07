package models

import (
	"github.com/google/uuid"
)

type CreateProjectDTO struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Key         string    `json:"key"`
	CreatedBy   uuid.UUID `json:"createdBy"`
}

type ProjectResponseDTO struct {
	ID             uuid.UUID                  `json:"id"`
	Name           string                     `json:"name"`
	Description    string                     `json:"description"`
	Key            string                     `json:"key"`
	CreatedBy      UserResponseDTO            `json:"createdBy"`
	ProjectMembers []ProjectMemberResponseDTO `json:"projectMembers"`
	Tasks          []TaskResponseDTO          `json:"tasks"`
	CreatedAt      string                     `json:"createdAt"`
	UpdatedAt      string                     `json:"updatedAt"`
}
