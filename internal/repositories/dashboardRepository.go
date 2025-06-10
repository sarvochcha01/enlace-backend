package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
)

type DashboardRepository interface {
	GetRecentlyAssignedTasks(userID uuid.UUID, limit int) ([]models.TaskResponseDTO, error)
	GetInProgressTasks(userID uuid.UUID, limit int) ([]models.TaskResponseDTO, error)
	GetApproachingDeadlineTasks(userID uuid.UUID, limit int) ([]models.TaskResponseDTO, error)

	SearchProjects(userID uuid.UUID, query string) ([]models.ProjectSearchResult, error)
	SearchTasks(userID uuid.UUID, query string) ([]models.TaskResponseDTO, error)
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
			t.task_number,
			p.id AS project_id,
			p.key AS project_key,
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
			t.task_number,
			p.id AS project_id,
			p.key AS project_key,
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
			t.task_number,
			p.id AS project_id,
			p.key AS project_key,
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

func (r *dashboardRepository) SearchProjects(userID uuid.UUID, query string) ([]models.ProjectSearchResult, error) {
	baseQuery := `
		SELECT
			p.id,
			p.name,
			p.description,
			p.key,
			COUNT(t.id) as total_tasks,
			COUNT(CASE WHEN t.status = 'completed' THEN 1 END) AS completed_tasks,
			COUNT(CASE
				WHEN t.status IN ('todo', 'in-progress') AND assigned_pm.user_id = $1
				THEN 1
			END) AS active_tasks_assigned_to_user
		FROM projects p
		JOIN project_members pm ON p.id = pm.project_id
		LEFT JOIN tasks t ON t.project_id = p.id
		LEFT JOIN project_members assigned_pm ON t.assigned_to = assigned_pm.id
		WHERE pm.user_id = $1
		AND pm.status = 'active'`

	var args []interface{}
	args = append(args, userID)
	argIndex := 2

	if query != "" {
		searchCondition := fmt.Sprintf(" AND (LOWER(p.name) LIKE LOWER($%d) OR LOWER(p.description) LIKE LOWER($%d) OR LOWER(p.key) LIKE LOWER($%d))", argIndex, argIndex, argIndex)
		baseQuery += searchCondition
		searchTerm := "%" + strings.ToLower(query) + "%"
		args = append(args, searchTerm)
	}

	baseQuery += `
		GROUP BY p.id, p.name, p.description, p.key
		ORDER BY MAX(p.updated_at) DESC`

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.ProjectSearchResult
	for rows.Next() {
		var project models.ProjectSearchResult
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.Key,
			&project.TotalTasks,
			&project.CompletedTasks,
			&project.ActiveTasksAssignedToUser,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (r *dashboardRepository) SearchTasks(userID uuid.UUID, query string) ([]models.TaskResponseDTO, error) {
	baseQuery := `
		SELECT
			t.id,
			t.title,
			t.task_number,
			p.id AS project_id,
			p.key AS project_key,
			p.name AS project_name,
			t.priority,
			t.due_date,
			t.status,
			t.assigned_to AS assigned_project_member_id, -- project_members.id for assignee
			assigned_pm.user_id AS assigned_user_id,      -- users.id for assignee
			assigned_user.name AS assigned_user_name     -- users.name for assignee
		FROM tasks t
		INNER JOIN projects p ON t.project_id = p.id
		INNER JOIN project_members pm ON p.id = pm.project_id -- For checking user's access
		LEFT JOIN project_members assigned_pm ON t.assigned_to = assigned_pm.id -- For assignee details
		LEFT JOIN users assigned_user ON assigned_pm.user_id = assigned_user.id
		WHERE pm.user_id = $1
		AND pm.status = 'active'`

	var args []interface{}
	args = append(args, userID)
	argIndex := 2

	if query != "" {
		// Ensure search terms apply to relevant fields. Original query searched p.name and p.key in tasks.
		// Adjusted to search t.title and also project context (p.name, p.key, or task_number with project_key prefix).
		// For simplicity, sticking to the original search scope for query, but this could be enhanced.
		// Example: (LOWER(t.title) LIKE LOWER($%d) OR (p.key || '-' || t.task_number::TEXT) LIKE LOWER($%d))
		searchCondition := fmt.Sprintf(" AND (LOWER(t.title) LIKE LOWER($%d) OR LOWER(p.name) LIKE LOWER($%d) OR LOWER(p.key) LIKE LOWER($%d))", argIndex, argIndex, argIndex)
		baseQuery += searchCondition
		searchTerm := "%" + strings.ToLower(query) + "%"
		args = append(args, searchTerm)
	}

	baseQuery += `
		ORDER BY 
			CASE WHEN assigned_pm.user_id = $1 THEN 0 ELSE 1 END, -- Tasks assigned to the current user first
			t.due_date ASC NULLS LAST,
			t.updated_at DESC`

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search tasks query: %w", err)
	}
	defer rows.Close()

	var tasks []models.TaskResponseDTO
	for rows.Next() {
		var task models.TaskResponseDTO
		var assignedProjectMemberIDNullable uuid.NullUUID // Use uuid.NullUUID for nullable UUID
		var assignedUserIDNullable uuid.NullUUID
		var assignedUserNameNullable sql.NullString // Use sql.NullString for nullable string

		// Note: The TaskResponseDTO has other fields like CreatedBy, UpdatedBy, Description, CreatedAt, UpdatedAt.
		// These are not selected in the current SQL query and will remain as zero values.
		// If they are needed, the SQL query and this Scan must be updated.
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.TaskNumber,
			&task.ProjectID,
			&task.ProjectKey,
			&task.ProjectName,
			&task.Priority, // Assumes TaskPriority is string/int compatible or implements Scanner
			&task.DueDate,  // Nullable timestamp
			&task.Status,   // Assumes TaskStatus is string/int compatible or implements Scanner
			&assignedProjectMemberIDNullable,
			&assignedUserIDNullable,
			&assignedUserNameNullable,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task row: %w", err)
		}

		if assignedProjectMemberIDNullable.Valid {
			task.AssignedTo = &models.ProjectMemberResponseDTO{
				ID: assignedProjectMemberIDNullable.UUID,
			}
			if assignedUserIDNullable.Valid {
				task.AssignedTo.UserID = assignedUserIDNullable.UUID
			}
			if assignedUserNameNullable.Valid {
				task.AssignedTo.Name = assignedUserNameNullable.String
				task.AssignedToName = assignedUserNameNullable.String // Populate the separate name field
			} else {
				// Name might be null if user somehow has no name, or if user was deleted
				// but project_member record still exists (though FKs should prevent this unless user name is nullable)
				task.AssignedToName = ""
			}
		} else {
			task.AssignedTo = nil
			task.AssignedToName = "" // No assignee, so name is empty or appropriately handled
		}

		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating task rows: %w", err)
	}

	return tasks, nil
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
			&task.TaskNumber,
			&task.ProjectID,
			&task.ProjectKey,
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
