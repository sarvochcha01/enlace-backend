package services

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
)

type ProjectMemberService interface {
	CreateProjectMember(createProjectMemberDTO *models.CreateProjectMemberDTO, firebaseUID string) error
	CreateProjectMemberTx(*sql.Tx, *models.CreateProjectMemberDTO) (uuid.UUID, error)
	GetProjectMemberID(userID uuid.UUID, projectID uuid.UUID) (uuid.UUID, error)
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

func (s *projectMemberService) GetProjectMemberID(userID uuid.UUID, projectID uuid.UUID) (uuid.UUID, error) {
	return s.projectMemberRepository.GetProjectMemberID(userID, projectID)
}
