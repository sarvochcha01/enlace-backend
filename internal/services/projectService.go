package services

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
	"github.com/sarvochcha01/enlace-backend/internal/utils"
)

type ProjectService interface {
	CreateProject(projectDTO *models.CreateProjectDTO, firebaseUID string) error
	GetProjectByID(uuid.UUID, string) (*models.ProjectResponseDTO, error)
	GetAllProjectsForUser(firebaseUID string) ([]models.ProjectResponseDTO, error)
	EditProject(firebaseUID string, projectID uuid.UUID, projectDTO *models.EditProjectDTO) error
	DeleteProject(firebaseUID string, projectID uuid.UUID) error

	GetProjectName(projectID uuid.UUID) (string, error)

	LeaveProject(projectID uuid.UUID, firebaseUID string) error
	JoinProject(projectID uuid.UUID, firebaseUID string) error
	GetProjectCreatorID(projectID uuid.UUID) (uuid.UUID, error)
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
		Role:      models.RoleOwner,
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

func (s *projectService) GetProjectByID(projectID uuid.UUID, firebaseUID string) (*models.ProjectResponseDTO, error) {

	userID, err := s.userService.FindUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		log.Println("UserID not found: ", err)
		return nil, errors.New("UserID not found: " + err.Error())
	}

	projectMember, err := s.projectMemberService.GetProjectMemberByUserID(userID, projectID)
	if err != nil {
		log.Println("Not a project member: ", err)
		return nil, errors.New("not a project member: " + err.Error())
	}

	if projectMember.Status == models.StatusInactive {
		log.Println("Project member is inactive")
		return nil, errors.New("project member is inactive")
	}

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

func (s *projectService) GetProjectName(projectID uuid.UUID) (string, error) {
	return s.projectRepository.GetProjectName(projectID)
}

func (s *projectService) GetProjectCreatorID(projectID uuid.UUID) (uuid.UUID, error) {
	return s.projectRepository.GetProjectCreatorID(projectID)
}

func (s *projectService) DeleteProject(firebaseUID string, projectID uuid.UUID) error {

	projectMember, err := s.projectMemberService.GetProjectMemberByFirebaseUID(firebaseUID, projectID)
	if err != nil {
		return err
	}

	if !utils.HasDeletePrivilege(projectMember) {
		return errors.New("only owners an delete the project")
	}

	return s.projectRepository.DeleteProject(projectID)
}

// TODO: Only owner and editor should be able to edit the project
func (s *projectService) EditProject(firebaseUID string, projectID uuid.UUID, projectDTO *models.EditProjectDTO) error {
	return s.projectRepository.EditProject(projectID, projectDTO)
}

func (s *projectService) JoinProject(projectID uuid.UUID, firebaseUID string) error {

	userID, err := s.userService.FindUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return err
	}

	projectMember, err := s.projectMemberService.GetProjectMemberByUserID(userID, projectID)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			var projectMemberDTO models.CreateProjectMemberDTO
			projectMemberDTO.ProjectID = projectID
			projectMemberDTO.Role = models.RoleViewer

			return s.projectMemberService.CreateProjectMember(&projectMemberDTO, firebaseUID)
		}

		return err

	}

	if projectMember.Status == models.StatusActive {
		return errors.New("already an active project member")
	}

	return s.projectMemberService.UpdateProjectMemberStatus(projectMember.ID, models.StatusActive)
}

func (s *projectService) LeaveProject(projectID uuid.UUID, firebaseUID string) error {

	userID, err := s.userService.FindUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return err
	}

	projectMemberID, err := s.projectMemberService.GetProjectMemberID(userID, projectID)
	if err != nil {
		return err
	}

	return s.projectMemberService.UpdateProjectMemberStatus(projectMemberID, models.StatusInactive)
}
