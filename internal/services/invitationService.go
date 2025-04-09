package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
	"github.com/sarvochcha01/enlace-backend/internal/utils"
)

type InvitationService interface {
	CreateInvitation(firebaseUID string, createInvitationDTO *models.CreateInvitationDTO) error
	GetInvitations(firebaseUID string) ([]models.InvitationResponseDTO, error)
	HasInvitation(userID uuid.UUID, projectID uuid.UUID) bool
	HasInvitationFirebaseUID(firebaseUID string, projectID uuid.UUID) bool
	EditInvitation(firebaseUID string, editInvitationDTO models.EditInvitationDTO) error
}

type invitationService struct {
	invitationRepository repositories.InvitationRepository
	userService          UserService
	projectService       ProjectService
	projectMemberService ProjectMemberService
	notificationService  NotificationService
}

func NewInvitationService(ir repositories.InvitationRepository, us UserService, ps ProjectService, pms ProjectMemberService, ns NotificationService) InvitationService {
	return &invitationService{invitationRepository: ir, userService: us, projectService: ps, projectMemberService: pms, notificationService: ns}
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

	err = s.invitationRepository.CreateInvitation(createInvitationDTO)
	if err == nil {
		projectName, err := s.projectService.GetProjectName(createInvitationDTO.ProjectID)
		if err == nil {
			notification := &models.CreateNotificationDTO{
				UserID:    createInvitationDTO.InvitedUserID,
				Type:      models.NotificationTypeProjectInvitation,
				Content:   fmt.Sprintf("You have been invited to join the project: %s", projectName),
				ProjectID: createInvitationDTO.ProjectID,
				TaskID:    nil,
			}

			err = s.notificationService.CreateNotification(*notification)
			if err != nil {
				fmt.Println("Failed to create notification:", err)
			}
		}
	}

	return err
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

func (s *invitationService) HasInvitationFirebaseUID(firebaseUID string, projectID uuid.UUID) bool {

	userID, err := s.userService.GetUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return false
	}

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
