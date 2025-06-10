package models

import "github.com/google/uuid"

type SearchResult struct {
	Projects []ProjectSearchResult `json:"projects"`
	Tasks    []TaskResponseDTO     `json:"tasks"`
}

type ProjectSearchResult struct {
	ID                        uuid.UUID `json:"id"`
	Name                      string    `json:"name"`
	Description               *string   `json:"description"`
	Key                       string    `json:"key"`
	TotalTasks                int       `json:"total_tasks"`
	CompletedTasks            int       `json:"completed_tasks"`
	ActiveTasksAssignedToUser int       `json:"active_tasks_assigned_to_user"`
}
