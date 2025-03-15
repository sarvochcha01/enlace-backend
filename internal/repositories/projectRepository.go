package repositories

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
)

type ProjectRepository interface {
	BeginTransaction() (*sql.Tx, error)
	CreateProject(*sql.Tx, *models.CreateProjectDTO) (uuid.UUID, error)
	GetAllProjectsForUser(uuid.UUID) ([]models.ProjectResponseDTO, error)
	EditProject(uuid.UUID, *models.EditProjectDTO) error
	DeleteProject(projectID uuid.UUID) error

	GetProjectByID(uuid.UUID) (*models.ProjectResponseDTO, error)
	GetProjectName(uuid.UUID) (string, error)
	GetProjectCreatorID(projectID uuid.UUID) (uuid.UUID, error)
	// JoinProject(userID uuid.UUID, projectID uuid.UUID) error
}

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) BeginTransaction() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *projectRepository) CreateProject(tx *sql.Tx, projectDTO *models.CreateProjectDTO) (uuid.UUID, error) {

	queryString := `INSERT INTO projects (name, description, key, created_by) VALUES ($1, $2, $3, $4) RETURNING id`

	var newProjectID uuid.UUID
	err := tx.QueryRow(queryString, projectDTO.Name, projectDTO.Description, projectDTO.Key, projectDTO.CreatedBy).Scan(&newProjectID)

	if err != nil {
		return uuid.Nil, err
	}

	return newProjectID, nil

}

func (r *projectRepository) GetProjectByID(projectID uuid.UUID) (*models.ProjectResponseDTO, error) {
	var projectDTO models.ProjectResponseDTO

	// Query basic project information with creator details
	queryString := `
        SELECT p.id, p.name, p.description, p.key, u.id, u.name, u.email, p.created_at, p.updated_at
        FROM projects p
        JOIN users u ON p.created_by = u.id
        WHERE p.id = $1
    `
	err := r.db.QueryRow(queryString, projectID).Scan(
		&projectDTO.ID,
		&projectDTO.Name,
		&projectDTO.Description,
		&projectDTO.Key,
		&projectDTO.CreatedBy.ID,
		&projectDTO.CreatedBy.Name,
		&projectDTO.CreatedBy.Email,
		&projectDTO.CreatedAt,
		&projectDTO.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Query project members
	projectDTO.ProjectMembers = []models.ProjectMemberResponseDTO{}
	membersQuery := `
        SELECT pm.id, pm.user_id, pm.project_id, pm.role, pm.joined_at, pm.status, u.name, u.email
        FROM project_members pm
        JOIN users u ON pm.user_id = u.id
        WHERE pm.project_id = $1
		AND pm.status = 'active'
    `
	rows, err := r.db.Query(membersQuery, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member models.ProjectMemberResponseDTO
		err := rows.Scan(&member.ID, &member.UserID, &member.ProjectID, &member.Role, &member.JoinedAt, &member.Status, &member.Name, &member.Email)
		if err != nil {
			return nil, err
		}
		projectDTO.ProjectMembers = append(projectDTO.ProjectMembers, member)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Query invitations
	projectDTO.Invitations = []models.InvitationResponseDTO{}
	invitationsQuery := `
		SELECT i.id, i.invited_by, i.invited_user_id, i.project_id, i.status, i.created_at, u.name, u.email
		FROM invitations i
		LEFT JOIN users u ON u.id = i.invited_user_id
		WHERE project_id = $1
	`
	invitationRows, err := r.db.Query(invitationsQuery, projectID)
	if err != nil {
		return nil, err
	}

	defer invitationRows.Close()

	for invitationRows.Next() {
		var invitation models.InvitationResponseDTO
		err := invitationRows.Scan(&invitation.ID, &invitation.InvitedBy, &invitation.InvitedUserID, &invitation.ProjectID, &invitation.Status, &invitation.InvitedAt, &invitation.Name, &invitation.Email)
		if err != nil {
			return nil, err
		}

		projectDTO.Invitations = append(projectDTO.Invitations, invitation)
	}

	// Query tasks with creator, updater, and assignee details
	projectDTO.Tasks = []models.TaskResponseDTO{}
	tasksQuery := `
        SELECT t.id, t.project_id, t.task_number, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
               -- Created by details
               cb_pm.id, cb_u.id, cb_u.name, cb_u.email, cb_pm.role, cb_pm.joined_at,
               -- Updated by details
               ub_pm.id, ub_u.id, ub_u.name, ub_u.email, ub_pm.role, ub_pm.joined_at,
               -- Assigned to details
               at_pm.id, at_u.id, at_u.name, at_u.email, at_pm.role, at_pm.joined_at
        FROM tasks t
        LEFT JOIN project_members cb_pm ON t.created_by = cb_pm.id 
        LEFT JOIN users cb_u ON cb_pm.user_id = cb_u.id
        LEFT JOIN project_members ub_pm ON t.updated_by = ub_pm.id
        LEFT JOIN users ub_u ON ub_pm.user_id = ub_u.id
        LEFT JOIN project_members at_pm ON t.assigned_to = at_pm.id
        LEFT JOIN users at_u ON at_pm.user_id = at_u.id
        WHERE t.project_id = $1
        ORDER BY t.task_number ASC
    `

	taskRows, err := r.db.Query(tasksQuery, projectID)
	if err != nil {
		return nil, err
	}
	defer taskRows.Close()

	for taskRows.Next() {
		var task models.TaskResponseDTO

		// Variables for nullable project member fields
		var createdByID, createdByUserID, createdByName, createdByEmail, createdByRole, createdByJoined sql.NullString
		var updatedByID, updatedByUserID, updatedByName, updatedByEmail, updatedByRole, updatedByJoined sql.NullString
		var assignedToID, assignedToUserID, assignedToName, assignedToEmail, assignedToRole, assignedToJoined sql.NullString

		err := taskRows.Scan(
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
			// Created by fields
			&createdByID,
			&createdByUserID,
			&createdByName,
			&createdByEmail,
			&createdByRole,
			&createdByJoined,
			// Updated by fields
			&updatedByID,
			&updatedByUserID,
			&updatedByName,
			&updatedByEmail,
			&updatedByRole,
			&updatedByJoined,
			// Assigned to fields
			&assignedToID,
			&assignedToUserID,
			&assignedToName,
			&assignedToEmail,
			&assignedToRole,
			&assignedToJoined,
		)
		if err != nil {
			return nil, err
		}

		// Set CreatedBy if not null
		if createdByID.Valid {
			var createdByUUID, createdByUserUUID uuid.UUID
			if err := createdByUUID.Scan(createdByID.String); err != nil {
				return nil, err
			}
			if err := createdByUserUUID.Scan(createdByUserID.String); err != nil {
				return nil, err
			}

			task.CreatedBy = models.ProjectMemberResponseDTO{
				ID:       createdByUUID,
				UserID:   createdByUserUUID,
				Name:     createdByName.String,
				Email:    createdByEmail.String,
				Role:     models.ProjectMemberRole(createdByRole.String),
				JoinedAt: createdByJoined.String,
			}
		}

		// Set UpdatedBy if not null
		if updatedByID.Valid {
			var updatedByUUID, updatedByUserUUID uuid.UUID
			if err := updatedByUUID.Scan(updatedByID.String); err != nil {
				return nil, err
			}
			if err := updatedByUserUUID.Scan(updatedByUserID.String); err != nil {
				return nil, err
			}

			task.UpdatedBy = models.ProjectMemberResponseDTO{
				ID:       updatedByUUID,
				UserID:   updatedByUserUUID,
				Name:     updatedByName.String,
				Email:    updatedByEmail.String,
				Role:     models.ProjectMemberRole(updatedByRole.String),
				JoinedAt: updatedByJoined.String,
			}
		}

		task.AssignedTo = nil

		// Set AssignedTo if not null
		if assignedToID.Valid {
			var assignedToUUID, assignedToUserUUID uuid.UUID
			if err := assignedToUUID.Scan(assignedToID.String); err != nil {
				return nil, err
			}
			if err := assignedToUserUUID.Scan(assignedToUserID.String); err != nil {
				return nil, err
			}

			task.AssignedTo = &models.ProjectMemberResponseDTO{
				ID:       assignedToUUID,
				UserID:   assignedToUserUUID,
				Name:     assignedToName.String,
				Email:    assignedToEmail.String,
				Role:     models.ProjectMemberRole(assignedToRole.String),
				JoinedAt: assignedToJoined.String,
			}
		}

		projectDTO.Tasks = append(projectDTO.Tasks, task)
	}

	if err := taskRows.Err(); err != nil {
		return nil, err
	}

	return &projectDTO, nil
}

func (r *projectRepository) GetAllProjectsForUser(userID uuid.UUID) ([]models.ProjectResponseDTO, error) {
	var projects []models.ProjectResponseDTO

	queryString := `
        SELECT p.id, p.name, p.description, p.key, u.id, u.name, u.email, p.created_at, p.updated_at
        FROM projects p
        JOIN project_members pm ON p.id = pm.project_id
		JOIN users u ON p.created_by = u.id											
        WHERE pm.user_id = $1
		AND pm.status = 'active'
    `

	rows, err := r.db.Query(queryString, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var project models.ProjectResponseDTO
		var createdBy models.UserResponseDTO
		err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.Key, &createdBy.ID, &createdBy.Name, &createdBy.Email, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, err
		}
		project.CreatedBy = createdBy
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) DeleteProject(projectID uuid.UUID) error {
	querystring := `
		DELETE FROM projects WHERE id = $1
	`
	_, err := r.db.Exec(querystring, projectID)

	return err
}

func (r *projectRepository) GetProjectName(projectID uuid.UUID) (string, error) {
	queryString := `
		SELECT name FROM projects WHERE id = $1
	`

	var name string
	err := r.db.QueryRow(queryString, projectID).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (r *projectRepository) GetProjectCreatorID(projectID uuid.UUID) (uuid.UUID, error) {

	var creatorID uuid.UUID
	queryString := `
		SELECT created_by
		FROM projects
		WHERE id = $1
	`

	if err := r.db.QueryRow(queryString, projectID).Scan(&creatorID); err != nil {
		return uuid.Nil, err
	}

	return creatorID, nil
}

func (r *projectRepository) EditProject(projectID uuid.UUID, projectDTO *models.EditProjectDTO) error {
	queryString := `
		UPDATE projects
		SET name = $1,
			description = $2
		WHERE id = $3
	`

	result, err := r.db.Exec(queryString, projectDTO.Name, projectDTO.Description, projectID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("project not found")
	}

	return nil
}
