package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string
type TaskPriority string

const (
	Todo       TaskStatus = "todo"
	InProgress TaskStatus = "in-progress"
	Completed  TaskStatus = "completed"
)

const (
	Low      TaskPriority = "low"
	Medium   TaskPriority = "medium"
	High     TaskPriority = "high"
	Critical TaskPriority = "critical"
)

type CreateTaskDTO struct {
	ProjectID   uuid.UUID    `json:"projectId"`
	CreatedBy   uuid.UUID    `json:"createdBy"`
	UpdatedBy   uuid.UUID    `json:"updatedBy"`
	AssignedTo  *uuid.UUID   `json:"assignedTo"`
	Title       string       `json:"title"`
	Description *string      `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	DueDate     *time.Time   `json:"dueDate"`
}

type TaskResponseDTO struct {
	ID             uuid.UUID                 `json:"id"`
	ProjectID      uuid.UUID                 `json:"projectId"`
	ProjectKey     string                    `json:"projectKey"`
	ProjectName    string                    `json:"projectName"`
	CreatedBy      ProjectMemberResponseDTO  `json:"createdBy"`
	UpdatedBy      ProjectMemberResponseDTO  `json:"updatedBy"`
	AssignedTo     *ProjectMemberResponseDTO `json:"assignedTo"`
	Title          string                    `json:"title"`
	TaskNumber     int                       `json:"taskNumber"`
	Description    *string                   `json:"description"`
	Status         TaskStatus                `json:"status"`
	Priority       TaskPriority              `json:"priority"`
	DueDate        *time.Time                `json:"dueDate"`
	AssignedToName string                    `json:"assignedToName"`
	CreatedAt      time.Time                 `json:"createdAt"`
	UpdatedAt      time.Time                 `json:"updatedAt"`
}

type UpdateTaskDTO struct {
	UpdatedBy   uuid.UUID    `json:"updatedBy"`
	AssignedTo  *uuid.UUID   `json:"assignedTo"`
	Title       string       `json:"title"`
	Description *string      `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	DueDate     *time.Time   `json:"dueDate,omitempty"`
}

type DeleteTaskDTO struct {
	TaskID      uuid.UUID `json:"commentId"`
	ProjectID   uuid.UUID `json:"projectId"`
	FirebaseUID string    `json:"firebaseUID"`
}
