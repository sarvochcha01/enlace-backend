package services

import (
	"errors"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
)

type ProjectService interface {
	CreateProject(projectDTO *models.CreateProjectDTO, firebaseUID string) error
	GetProjectByID(uuid.UUID) (*models.ProjectResponseDTO, error)
	GetAllProjectsForUser(firebaseUID string) ([]models.ProjectResponseDTO, error)
	// JoinProject(firebaseUID string, projectID uuid.UUID) error
}

type projectService struct {
	projectRepository    repositories.ProjectRepository
	userService          UserService
	projectMemberService ProjectMemberService
}

func NewProjectService(pr repositories.ProjectRepository, us UserService, ps ProjectMemberService) ProjectService {
	return &projectService{projectRepository: pr, userService: us, projectMemberService: ps}
}

func (s *projectService) CreateProject(projectDTO *models.CreateProjectDTO, firebaseUID string) error {

	userID, err := s.userService.FindUserIDByFirebaseUID(firebaseUID)

	if err != nil {
		log.Println("UserID not found: ", err)
		return errors.New("UserID not found: " + err.Error())
	}

	projectDTO.Key = strings.ToUpper(projectDTO.Key)
	projectDTO.CreatedBy = userID

	var projectID uuid.UUID

	tx, err := s.projectRepository.BeginTransaction()
	if err != nil {
		return err
	}

	projectID, err = s.projectRepository.CreateProject(tx, projectDTO)

	if err != nil {
		tx.Rollback()
		log.Println("Failed to create project: ", err)
		return err
	}

	projectMemberDTO := &models.CreateProjectMemberDTO{
		UserID:    userID,
		ProjectID: projectID,
		Role:      models.Owner,
	}

	_, err = s.projectMemberService.CreateProjectMemberTx(tx, projectMemberDTO)
	if err != nil {
		tx.Rollback()
		log.Println("Failed to add creator as project member: ", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Failed to commit transaction:", err)
		return err
	}

	return nil
}

// TODO: Add verification so that only the project members can access it
func (s *projectService) GetProjectByID(projectID uuid.UUID) (*models.ProjectResponseDTO, error) {
	return s.projectRepository.GetProjectByID(projectID)
}

func (s *projectService) GetAllProjectsForUser(firebaseUID string) ([]models.ProjectResponseDTO, error) {

	userID, err := s.userService.FindUserIDByFirebaseUID(firebaseUID)

	if err != nil {
		log.Println("UserID not found: ", err)
		return nil, errors.New("UserID not found: " + err.Error())
	}

	return s.projectRepository.GetAllProjectsForUser(userID)
}

// TODO: Modify this so only the invited users can access it
// func (s *projectService) JoinProject(firebaseUID string, projectID uuid.UUID) error {
// 	userID, err := s.userService.FindUserIDByFirebaseUID(firebaseUID)

// 	if err != nil {
// 		log.Println("UserID not found: ", err)
// 		return errors.New("UserID not found: " + err.Error())
// 	}

// 	return nil
// }
