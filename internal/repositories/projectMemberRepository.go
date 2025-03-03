package repositories

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
)

type ProjectMemberRepository interface {
	CreateProjectMember(*models.CreateProjectMemberDTO) error
	CreateProjectMemberTx(*sql.Tx, *models.CreateProjectMemberDTO) (uuid.UUID, error)
	GetProjectMemberID(uuid.UUID, uuid.UUID) (uuid.UUID, error)
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

func (s *projectMemberRepository) GetProjectMemberID(userID uuid.UUID, projectID uuid.UUID) (uuid.UUID, error) {

	var projectMemberID uuid.UUID

	queryString := `
		SELECT id FROM project_members
		WHERE user_id = $1 AND project_id = $2
	`

	err := s.db.QueryRow(queryString, userID, projectID).Scan(&projectMemberID)

	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, fmt.Errorf("no project member found for user %s in project %s", userID, projectID)
		}
		return uuid.Nil, err
	}

	return projectMemberID, nil
}
