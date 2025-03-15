package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
)

type InvitationService interface {
	CreateInvitation(firebaseUID string, createInvitationDTO *models.CreateInvitationDTO) error
	GetInvitations(firebaseUID string) ([]models.InvitationResponseDTO, error)
	HasInvitation(userID uuid.UUID, projectID uuid.UUID) bool
	EditInvitation(firebaseUID string, editInvitationDTO models.EditInvitationDTO) error
}

type invitationService struct {
	invitationRepository repositories.InvitationRepository
	userService          UserService
	projectService       ProjectService
}

func NewInvitationService(ir repositories.InvitationRepository, us UserService, ps ProjectService) InvitationService {
	return &invitationService{invitationRepository: ir, userService: us, projectService: ps}
}

func (r *invitationService) CreateInvitation(firebaseUID string, createInvitationDTO *models.CreateInvitationDTO) error {

	userID, err := r.userService.GetUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return err
	}

	createInvitationDTO.InvitedBy = userID

	return r.invitationRepository.CreateInvitation(createInvitationDTO)
}

func (r *invitationService) GetInvitations(firebaseUID string) ([]models.InvitationResponseDTO, error) {

	userID, err := r.userService.GetUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}

	return r.invitationRepository.GetInvitations(userID)
}

func (r *invitationService) HasInvitation(userID uuid.UUID, projectID uuid.UUID) bool {

	return r.invitationRepository.HasInvitation(userID, projectID)
}

func (s *invitationService) EditInvitation(firebaseUID string, editInvitationDTO models.EditInvitationDTO) error {
	userID, err := s.userService.GetUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return nil
	}

	hasInvitation := s.HasInvitation(userID, editInvitationDTO.ProjectID)

	if hasInvitation {
		if editInvitationDTO.Status == string(models.InivtationStatusAccepted) {
			if err = s.projectService.JoinProject(editInvitationDTO.ProjectID, firebaseUID); err != nil {
				return err
			}
		}
		return s.invitationRepository.EditInvitation(editInvitationDTO)
	} else {
		return errors.New("invitation for the user and project doesnt exist")
	}

}
