package models

import (
	"github.com/google/uuid"
)

type CreateCommentDTO struct {
	ProjectID uuid.UUID `json:"project_id"`
	TaskID    uuid.UUID `json:"task_id"`
	CreatedBy uuid.UUID `json:"created_by"`
	Comment   string    `json:"comment"`
}

type CommentDTO struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	TaskID    uuid.UUID `json:"task_id"`
	CreatedBy uuid.UUID `json:"created_by"`
	Comment   string    `json:"comment"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type UpdateCommentDTO struct {
	ProjectID uuid.UUID `json:"project_id"`
	CommentID uuid.UUID `json:"comment_id"`
	Comment   string    `json:"comment"`
}

type DeleteCommentDTO struct {
	ProjectID uuid.UUID `json:"project_id"`
	CommentID uuid.UUID `json:"comment_id"`
}
