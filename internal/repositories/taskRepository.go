package repositories

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
)

type TaskRepository interface {
	CreateTask(*models.CreateTaskDTO) (uuid.UUID, error)
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(taskDTO *models.CreateTaskDTO) (uuid.UUID, error) {
	queryString := `
	INSERT INTO tasks 
	(project_id, created_by, updated_by, assigned_to, title, description, status, priority, due_date) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
	RETURNING id
	`

	var taskID uuid.UUID

	err := r.db.QueryRow(queryString, taskDTO.ProjectID, taskDTO.CreatedBy, taskDTO.UpdatedBy, taskDTO.AssignedTo,
		taskDTO.Title, taskDTO.Description, taskDTO.Status, taskDTO.Priority, taskDTO.DueDate,
	).Scan(&taskID)

	if err != nil {
		log.Fatal("Failed to insert task:", err)
		return uuid.Nil, err
	}

	return taskID, nil

}
