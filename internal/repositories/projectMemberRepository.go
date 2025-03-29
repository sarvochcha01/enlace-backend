package repositories

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
)

type ProjectMemberRepository interface {
	CreateProjectMember(*models.CreateProjectMemberDTO) error
	CreateProjectMemberTx(*sql.Tx, *models.CreateProjectMemberDTO) (uuid.UUID, error)
	GetUserID(projectMemberID uuid.UUID) (uuid.UUID, error)
	GetProjectMemberID(uuid.UUID, uuid.UUID) (uuid.UUID, error)
	GetProjectMember(uuid.UUID) (*models.ProjectMemberResponseDTO, error)
	UpdateProjectMemberStatus(uuid.UUID, models.ProjectMemberStatus) error
	UpdateProjectMemberRole(projectMemberID uuid.UUID, newRole models.ProjectMemberRole) error
	GetProjectMemberByUserID(userID uuid.UUID, projectID uuid.UUID) (*models.ProjectMemberResponseDTO, error)
}

type projectMemberRepository struct {
	db *sql.DB
}

func NewProjectMemberRepository(db *sql.DB) ProjectMemberRepository {
	return &projectMemberRepository{db: db}
}

func (r *projectMemberRepository) CreateProjectMember(projectMember *models.CreateProjectMemberDTO) error {

	queryString := `INSERT INTO project_members (user_id, project_id, role) VALUES ($1, $2, $3) RETURNING id`

	var newID uuid.UUID

	err := r.db.QueryRow(queryString, projectMember.UserID, projectMember.ProjectID, projectMember.Role).Scan(&newID)

	if err != nil {
		return err
	}

	return nil
}

func (r *projectMemberRepository) CreateProjectMemberTx(tx *sql.Tx, projectMember *models.CreateProjectMemberDTO) (uuid.UUID, error) {

	queryString := `INSERT INTO project_members (user_id, project_id, role) VALUES ($1, $2, $3) RETURNING id`

	var newID uuid.UUID

	err := tx.QueryRow(queryString, projectMember.UserID, projectMember.ProjectID, projectMember.Role).Scan(&newID)

	if err != nil {
		return uuid.Nil, err
	}

	return newID, nil
}

func (r *projectMemberRepository) GetProjectMemberByUserID(userID uuid.UUID, projectID uuid.UUID) (*models.ProjectMemberResponseDTO, error) {

	var projectMember models.ProjectMemberResponseDTO

	queryString := `
		SELECT id, user_id, project_id, role, joined_at, status
		FROM project_members
		WHERE user_id = $1 AND project_id = $2
	`

	err := r.db.QueryRow(queryString, userID, projectID).Scan(
		&projectMember.ID,
		&projectMember.UserID,
		&projectMember.ProjectID,
		&projectMember.Role,
		&projectMember.JoinedAt,
		&projectMember.Status)

	if err != nil {
		return nil, err
	}

	return &projectMember, nil
}

func (r *projectMemberRepository) GetProjectMember(projectMemberID uuid.UUID) (*models.ProjectMemberResponseDTO, error) {

	var projectMember models.ProjectMemberResponseDTO

	queryString := `
		SELECT id, user_id, project_id, role, joined_at, status
		FROM project_members
		WHERE id = $1
	`

	err := r.db.QueryRow(queryString, projectMemberID).Scan(
		&projectMember.ID,
		&projectMember.UserID,
		&projectMember.ProjectID,
		&projectMember.Role,
		&projectMember.JoinedAt,
		&projectMember.Status)

	if err != nil {
		return nil, err
	}

	return &projectMember, nil
}

func (r *projectMemberRepository) GetUserID(projectMemberID uuid.UUID) (uuid.UUID, error) {

	var userID uuid.UUID

	queryString := `
		SELECT user_id FROM project_members
		WHERE id = $1
	`

	err := r.db.QueryRow(queryString, projectMemberID).Scan(&userID)

	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil

}

func (r *projectMemberRepository) GetProjectMemberID(userID uuid.UUID, projectID uuid.UUID) (uuid.UUID, error) {

	var projectMemberID uuid.UUID

	queryString := `
		SELECT id FROM project_members
		WHERE user_id = $1 AND project_id = $2
	`

	err := r.db.QueryRow(queryString, userID, projectID).Scan(&projectMemberID)

	if err != nil {
		return uuid.Nil, err
	}

	return projectMemberID, nil
}

func (r *projectMemberRepository) UpdateProjectMemberStatus(projectMemberID uuid.UUID, newStatus models.ProjectMemberStatus) error {

	queryString := `
		UPDATE project_members
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(queryString, newStatus, projectMemberID)
	if err != nil {
		return err
	}

	return nil
}

func (r *projectMemberRepository) UpdateProjectMemberRole(projectMemberID uuid.UUID, newRole models.ProjectMemberRole) error {

	queryString := `
		UPDATE project_members
		SET role = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(queryString, newRole, projectMemberID)
	if err != nil {
		return err
	}

	return nil
}
