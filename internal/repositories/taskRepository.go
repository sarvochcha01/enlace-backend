package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
)

type TaskRepository interface {
	CreateTask(*models.CreateTaskDTO) (uuid.UUID, error)
	GetTaskByID(uuid.UUID) (*models.TaskResponseDTO, error)
	EditTask(uuid.UUID, *models.UpdateTaskDTO) error
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
		log.Println("Failed to insert task:", err)
		return uuid.Nil, err
	}

	return taskID, nil

}

func (r *taskRepository) GetTaskByID(taskID uuid.UUID) (*models.TaskResponseDTO, error) {
	var task models.TaskResponseDTO
	var assignedToID, assignedToUserID sql.NullString
	var assignedToName, assignedToEmail, assignedToRole sql.NullString
	var assignedToJoinedAt sql.NullString

	queryString := `
        SELECT t.id, t.project_id, t.task_number, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
               -- Created by details
               cb_pm.id, cb_u.id, cb_u.name, cb_u.email, cb_pm.role, cb_pm.joined_at,
               -- Updated by details
               ub_pm.id, ub_u.id, ub_u.name, ub_u.email, ub_pm.role, ub_pm.joined_at,
               -- Assigned to details (might be NULL)
               at_pm.id, at_u.id, at_u.name, at_u.email, at_pm.role, at_pm.joined_at
        FROM tasks t
        LEFT JOIN project_members cb_pm ON t.created_by = cb_pm.id
        LEFT JOIN users cb_u ON cb_pm.user_id = cb_u.id
        LEFT JOIN project_members ub_pm ON t.updated_by = ub_pm.id
        LEFT JOIN users ub_u ON ub_pm.user_id = ub_u.id
        LEFT JOIN project_members at_pm ON t.assigned_to = at_pm.id
        LEFT JOIN users at_u ON at_pm.user_id = at_u.id
        WHERE t.id = $1
    `

	err := r.db.QueryRow(queryString, taskID).Scan(
		&task.ID,
		&task.ProjectId,
		&task.TaskNumber,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		// Created by
		&task.CreatedBy.ID,
		&task.CreatedBy.UserID,
		&task.CreatedBy.Name,
		&task.CreatedBy.Email,
		&task.CreatedBy.Role,
		&task.CreatedBy.JoinedAt,
		// Updated by
		&task.UpdatedBy.ID,
		&task.UpdatedBy.UserID,
		&task.UpdatedBy.Name,
		&task.UpdatedBy.Email,
		&task.UpdatedBy.Role,
		&task.UpdatedBy.JoinedAt,
		// Assigned To
		&assignedToID,
		&assignedToUserID,
		&assignedToName,
		&assignedToEmail,
		&assignedToRole,
		&assignedToJoinedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error fetching task: %w", err)
	}

	if assignedToID.Valid {
		var assignedToUUID, assignedToUserUUID uuid.UUID

		if err := assignedToUUID.Scan(assignedToID.String); err != nil {
			return nil, fmt.Errorf("error converting assigned_to ID: %w", err)
		}

		if err := assignedToUserUUID.Scan(assignedToUserID.String); err != nil {
			return nil, fmt.Errorf("error converting assigned_to user ID: %w", err)
		}

		task.AssignedTo = &models.ProjectMemberResponseDTO{
			ID:       assignedToUUID,
			UserID:   assignedToUserUUID,
			Name:     assignedToName.String,
			Email:    assignedToEmail.String,
			Role:     models.ProjectRole(assignedToRole.String),
			JoinedAt: assignedToJoinedAt.String,
		}
	} else {
		task.AssignedTo = nil
	}

	return &task, nil
}

func (r *taskRepository) EditTask(taskID uuid.UUID, updateTaskDTO *models.UpdateTaskDTO) error {
	queryString := `
	    UPDATE tasks
	    SET updated_by = $1,
	        assigned_to = $2,
	        title = $3,
	        description = $4,
	        status = $5,
	        priority = $6,
	        due_date = $7
	    WHERE id = $8
	`

	_, err := r.db.Exec(queryString, updateTaskDTO.UpdatedBy, updateTaskDTO.AssignedTo, updateTaskDTO.Title, updateTaskDTO.Description, updateTaskDTO.Status, updateTaskDTO.Priority, updateTaskDTO.DueDate, taskID)

	return err
}
