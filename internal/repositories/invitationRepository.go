package repositories

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
)

type InvitationRepository interface {
	CreateInvitation(createInvitationDTO *models.CreateInvitationDTO) error
	GetInvitations(userID uuid.UUID) ([]models.InvitationResponseDTO, error)
	HasInvitation(userID uuid.UUID, projectID uuid.UUID) bool
	EditInvitation(editInvitationDTO models.EditInvitationDTO) error
}

type invitationRepository struct {
	db *sql.DB
}

func NewInvitationRepository(db *sql.DB) InvitationRepository {
	return &invitationRepository{db: db}
}

func (r *invitationRepository) CreateInvitation(createInvitationDTO *models.CreateInvitationDTO) error {

	queryString := `
		INSERT INTO invitations (invited_by, invited_user_id, project_id)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(queryString, createInvitationDTO.InvitedBy, createInvitationDTO.InvitedUserID, createInvitationDTO.ProjectID)

	return err
}

func (r *invitationRepository) GetInvitations(userID uuid.UUID) ([]models.InvitationResponseDTO, error) {
	invitations := []models.InvitationResponseDTO{}
	invitationsQuery := `
		SELECT i.id, i.invited_by, i.invited_user_id, i.project_id, i.status, i.created_at, u.name, u.email, p.name
		FROM invitations i
		LEFT JOIN users u ON u.id = i.invited_by
		LEFT JOIN projects p ON p.id = i.project_id
		WHERE i.invited_user_id = $1
	`
	rows, err := r.db.Query(invitationsQuery, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var invitation models.InvitationResponseDTO
		err := rows.Scan(&invitation.ID, &invitation.InvitedBy, &invitation.InvitedUserID, &invitation.ProjectID, &invitation.Status, &invitation.InvitedAt, &invitation.Name, &invitation.Email, &invitation.ProjectName)
		if err != nil {
			return nil, err
		}

		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

func (r *invitationRepository) HasInvitation(userID uuid.UUID, projectID uuid.UUID) bool {
	queryString := `
		SELECT EXISTS (
			SELECT 1 FROM invitations 
			WHERE invited_user_id = $1 
			AND project_id = $2
		)
	`
	var exists bool
	err := r.db.QueryRow(queryString, userID, projectID).Scan(&exists)
	if err != nil {
		log.Println("Error checking invitation:", err)
		return false
	}
	return exists
}

func (r *invitationRepository) EditInvitation(editInvitationDTO models.EditInvitationDTO) error {
	queryString := `
		UPDATE invitations
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(queryString, editInvitationDTO.Status, editInvitationDTO.InvitationID)
	if err != nil {
		return err
	}

	return nil
}
