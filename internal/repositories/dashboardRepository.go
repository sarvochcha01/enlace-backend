package repositories

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
)

type DashboardRepository interface {
	GetRecentlyAssignedTasks(userID uuid.UUID, limit int) ([]models.TaskResponseDTO, error)
	GetInProgressTasks(userID uuid.UUID, limit int) ([]models.TaskResponseDTO, error)
	GetApproachingDeadlineTasks(userID uuid.UUID, limit int) ([]models.TaskResponseDTO, error)
}

type dashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetRecentlyAssignedTasks(userID uuid.UUID, limit int) ([]models.TaskResponseDTO, error) {
	queryString := `
		SELECT 
			t.id,
			t.title, 
			p.id AS project_id,
			p.name AS project_name, 
			t.priority, 
			t.due_date, 
			t.status
		FROM tasks t
		INNER JOIN projects p ON t.project_id = p.id
		LEFT JOIN project_members at_pm ON t.assigned_to = at_pm.id
		WHERE at_pm.user_id = $1
		ORDER BY t.updated_at DESC
		LIMIT $2
	`

	return r.fetchTasks(queryString, userID, limit)
}

func (r *dashboardRepository) GetInProgressTasks(userID uuid.UUID, limit int) ([]models.TaskResponseDTO, error) {
	query := `
		SELECT 
			t.id,
			t.title, 
			p.id AS project_id,
			p.name AS project_name, 
			t.priority, 
			t.due_date, 
			t.status
		FROM tasks t
		INNER JOIN projects p ON t.project_id = p.id
		LEFT JOIN project_members at_pm ON t.assigned_to = at_pm.id
		WHERE at_pm.user_id = $1
		AND t.status = 'in-progress'
		ORDER BY t.updated_at DESC
		LIMIT $2
	`

	return r.fetchTasks(query, userID, limit)
}

func (r *dashboardRepository) GetApproachingDeadlineTasks(userID uuid.UUID, limit int) ([]models.TaskResponseDTO, error) {
	query := `
		SELECT 
			t.id,
			t.title, 
			p.id AS project_id,
			p.name AS project_name, 
			t.priority, 
			t.due_date, 
			t.status
		FROM tasks t
		INNER JOIN projects p ON t.project_id = p.id
		LEFT JOIN project_members at_pm ON t.assigned_to = at_pm.id
		WHERE at_pm.user_id = $1
		AND t.due_date IS NOT NULL
		AND t.due_date >= NOW()
		AND t.due_date <= NOW() + INTERVAL '3 days'
		ORDER BY t.due_date ASC
		LIMIT $2
	`

	return r.fetchTasks(query, userID, limit)
}

func (r *dashboardRepository) fetchTasks(query string, userID uuid.UUID, limit int) ([]models.TaskResponseDTO, error) {
	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.TaskResponseDTO

	for rows.Next() {
		var task models.TaskResponseDTO
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.ProjectID,
			&task.ProjectName,
			&task.Priority,
			&task.DueDate,
			&task.Status,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
