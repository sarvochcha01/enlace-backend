package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
	"github.com/sarvochcha01/enlace-backend/internal/utils"
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
	projectMemberService ProjectMemberService
}

func NewInvitationService(ir repositories.InvitationRepository, us UserService, ps ProjectService, pms ProjectMemberService) InvitationService {
	return &invitationService{invitationRepository: ir, userService: us, projectService: ps, projectMemberService: pms}
}

func (s *invitationService) CreateInvitation(firebaseUID string, createInvitationDTO *models.CreateInvitationDTO) error {

	user, err := s.userService.GetUserByFirebaseUID(firebaseUID)
	if err != nil {
		return err
	}

	createInvitationDTO.InvitedBy = user.ID

	projectMember, err := s.projectMemberService.GetProjectMemberByFirebaseUID(firebaseUID, createInvitationDTO.ProjectID)
	if err != nil {
		return err
	}

	if !utils.HasEditPrivileges(projectMember) {
		return errors.New("you dont have the privileges to invite an user")
	}

	return s.invitationRepository.CreateInvitation(createInvitationDTO)
}

func (s *invitationService) GetInvitations(firebaseUID string) ([]models.InvitationResponseDTO, error) {

	userID, err := s.userService.GetUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}

	return s.invitationRepository.GetInvitations(userID)
}

func (s *invitationService) HasInvitation(userID uuid.UUID, projectID uuid.UUID) bool {

	return s.invitationRepository.HasInvitation(userID, projectID)
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
