package services

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
	"github.com/sarvochcha01/enlace-backend/internal/utils"
)

type ProjectMemberService interface {
	CreateProjectMember(createProjectMemberDTO *models.CreateProjectMemberDTO, firebaseUID string) error
	CreateProjectMemberTx(*sql.Tx, *models.CreateProjectMemberDTO) (uuid.UUID, error)
	GetProjectMemberID(userID uuid.UUID, projectID uuid.UUID) (uuid.UUID, error)
	GetProjectMemberIDByFirebaseUID(firebaseUID string, projectID uuid.UUID) (uuid.UUID, error)
	GetProjectMemberByFirebaseUID(firebaseUID string, projectID uuid.UUID) (*models.ProjectMemberResponseDTO, error)
	GetProjectMemberByUserID(userID uuid.UUID, projectID uuid.UUID) (*models.ProjectMemberResponseDTO, error)
	GetProjectMember(uuid.UUID) (*models.ProjectMemberResponseDTO, error)

	UpdateProjectMemberStatus(projectMemberID uuid.UUID, newStatus models.ProjectMemberStatus) error
	UpdateProjectMemberRole(firebaseUID string, updateProjectMemberDTO *models.UpdateProjectMemberDTO) error
}

type projectMemberService struct {
	projectMemberRepository repositories.ProjectMemberRepository
	userService             UserService
}

func NewProjectMemberService(pr repositories.ProjectMemberRepository, us UserService) ProjectMemberService {
	return &projectMemberService{projectMemberRepository: pr, userService: us}
}

func (s *projectMemberService) CreateProjectMember(createProjectMemberDTO *models.CreateProjectMemberDTO, firebaseUID string) error {

	userID, err := s.userService.FindUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		log.Println("UserID not found: ", err)
		return errors.New("UserID not found: " + err.Error())
	}

	createProjectMemberDTO.UserID = userID

	return s.projectMemberRepository.CreateProjectMember(createProjectMemberDTO)
}

func (s *projectMemberService) CreateProjectMemberTx(tx *sql.Tx, createProjectMemberDTO *models.CreateProjectMemberDTO) (uuid.UUID, error) {
	return s.projectMemberRepository.CreateProjectMemberTx(tx, createProjectMemberDTO)
}

func (s *projectMemberService) GetProjectMemberByUserID(userID uuid.UUID, projectID uuid.UUID) (*models.ProjectMemberResponseDTO, error) {
	return s.projectMemberRepository.GetProjectMemberByUserID(userID, projectID)
}

func (s *projectMemberService) GetProjectMemberID(userID uuid.UUID, projectID uuid.UUID) (uuid.UUID, error) {
	return s.projectMemberRepository.GetProjectMemberID(userID, projectID)
}

func (s *projectMemberService) GetProjectMember(projectMemberID uuid.UUID) (*models.ProjectMemberResponseDTO, error) {
	return s.projectMemberRepository.GetProjectMember(projectMemberID)
}

func (s *projectMemberService) GetProjectMemberIDByFirebaseUID(firebaseUID string, projectID uuid.UUID) (uuid.UUID, error) {

	userID, err := s.userService.FindUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		log.Println("Failed to get userID")
		return uuid.Nil, errors.New("failed to get userID")
	}

	return s.projectMemberRepository.GetProjectMemberID(userID, projectID)
}

func (s *projectMemberService) GetProjectMemberByFirebaseUID(firebaseUID string, projectID uuid.UUID) (*models.ProjectMemberResponseDTO, error) {
	userID, err := s.userService.FindUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		log.Println("Failed to get userID")
		return nil, errors.New("failed to get userID")
	}

	return s.projectMemberRepository.GetProjectMemberByUserID(userID, projectID)
}

func (s *projectMemberService) UpdateProjectMemberStatus(projectMemberID uuid.UUID, newStatus models.ProjectMemberStatus) error {
	return s.projectMemberRepository.UpdateProjectMemberStatus(projectMemberID, newStatus)
}

func (s *projectMemberService) UpdateProjectMemberRole(firebaseUID string, updateProjectMemberDTO *models.UpdateProjectMemberDTO) error {

	projectMemberWhoRequested, err := s.GetProjectMemberByFirebaseUID(firebaseUID, updateProjectMemberDTO.ProjectID)
	if err != nil {
		return err
	}

	if projectMemberWhoRequested.ID == updateProjectMemberDTO.ID {
		return errors.New("no you can't change your own role")
	}

	if !utils.HasEditPrivileges(projectMemberWhoRequested) {
		return errors.New("no edit privileges")
	}

	projectMemberToUpdate, err := s.GetProjectMember(updateProjectMemberDTO.ID)
	if err != nil {
		return err
	}

	if projectMemberToUpdate.Role == models.RoleOwner {
		return errors.New("owners can't be edited")
	}

	if projectMemberWhoRequested.Role == models.RoleEditor && updateProjectMemberDTO.Role == models.RoleOwner {
		return errors.New("lmao, editors can't make themselves owner")
	}

	return s.projectMemberRepository.UpdateProjectMemberRole(projectMemberToUpdate.ID, models.ProjectMemberRole(updateProjectMemberDTO.Role))
}
