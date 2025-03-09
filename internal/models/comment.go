package models

import (
	"github.com/google/uuid"
)

type CreateCommentDTO struct {
	ProjectID uuid.UUID `json:"projectId"`
	TaskID    uuid.UUID `json:"taskId"`
	CreatedBy uuid.UUID `json:"createdBy"`
	Comment   string    `json:"comment"`
}

type CommentResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"projectId"`
	TaskID    uuid.UUID `json:"taskId"`
	CreatedBy uuid.UUID `json:"createdBy"`
	Comment   string    `json:"comment"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}

type UpdateCommentDTO struct {
	ProjectID uuid.UUID `json:"projectId"`
	CommentID uuid.UUID `json:"commentId"`
	Comment   string    `json:"comment"`
}

type DeleteCommentDTO struct {
	ProjectID uuid.UUID `json:"projectId"`
	CommentID uuid.UUID `json:"commentId"`
}
