package models

import (
	"database/sql"
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
	ProjectID   uuid.UUID     `json:"projectId"`
	CreatedBy   uuid.UUID     `json:"createdBy"`
	UpdatedBy   uuid.UUID     `json:"updatedBy"`
	AssignedTo  uuid.NullUUID `json:"assignedTo"`
	Title       string        `json:"title"`
	Description *string       `json:"description"`
	Status      TaskStatus    `json:"status"`
	Priority    TaskPriority  `json:"priority"`
	DueDate     sql.NullTime  `json:"dueDate"`
}

type TaskResponseDTO struct {
	ID          uuid.UUID                 `json:"id"`
	ProjectId   uuid.UUID                 `json:"projectId"`
	CreatedBy   ProjectMemberResponseDTO  `json:"createdBy"`
	UpdatedBy   ProjectMemberResponseDTO  `json:"updatedBy"`
	AssignedTo  *ProjectMemberResponseDTO `json:"assignedTo"`
	Title       string                    `json:"title"`
	TaskNumber  int                       `json:"taskNumber"`
	Description sql.NullString            `json:"description"`
	Status      TaskStatus                `json:"status"`
	Priority    TaskPriority              `json:"priority"`
	DueDate     sql.NullTime              `json:"dueDate"`
	CreatedAt   time.Time                 `json:"createdAt"`
	UpdatedAt   time.Time                 `json:"updatedAt"`
}
