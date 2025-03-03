package services

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
)

type TaskService interface {
	CreateTask(taskDTo *models.CreateTaskDTO, firebaseUID string) (uuid.UUID, error)
	GetTaskByID(uuid.UUID) (*models.TaskResponseDTO, error)
}

type taskService struct {
	taskRepository       repositories.TaskRepository
	userService          UserService
	projectMemberService ProjectMemberService
}

func NewTaskService(tr repositories.TaskRepository, us UserService, pms ProjectMemberService) TaskService {
	return &taskService{taskRepository: tr, userService: us, projectMemberService: pms}
}

func (s *taskService) CreateTask(taskDTO *models.CreateTaskDTO, firebaseUID string) (uuid.UUID, error) {

	userID, err := s.userService.FindUserIDByFirebaseUID(firebaseUID)

	if err != nil {
		log.Fatal("UserID not found: ", err)
		return uuid.Nil, errors.New("UserID not found: " + err.Error())
	}

	var projectMemberID uuid.UUID
	projectMemberID, err = s.projectMemberService.GetProjectMemberID(userID, taskDTO.ProjectID)

	if err != nil {
		log.Fatal("Project Member not found: ", err)
		return uuid.Nil, errors.New("Project Member not found: " + err.Error())
	}

	taskDTO.CreatedBy = projectMemberID
	taskDTO.UpdatedBy = projectMemberID

	return s.taskRepository.CreateTask(taskDTO)
}

// TODO: Add verification so only the project members can access it
func (s *taskService) GetTaskByID(taskID uuid.UUID) (*models.TaskResponseDTO, error) {
	return s.taskRepository.GetTaskByID(taskID)
}
