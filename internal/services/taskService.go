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
	EditTask(uuid.UUID, uuid.UUID, string, *models.UpdateTaskDTO) error
	DeleteTask(*models.DeleteTaskDTO) error
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

	userID, err := s.userService.GetUserIDByFirebaseUID(firebaseUID)

	if err != nil {
		log.Println("UserID not found: ", err)
		return uuid.Nil, errors.New("UserID not found: " + err.Error())
	}

	var projectMemberID uuid.UUID
	projectMemberID, err = s.projectMemberService.GetProjectMemberID(userID, taskDTO.ProjectID)

	if err != nil {
		log.Println("Project Member not found: ", err)
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

func (s *taskService) EditTask(taskID uuid.UUID, projectID uuid.UUID, firebaseUID string, updateTaskDTO *models.UpdateTaskDTO) error {

	userID, err := s.userService.GetUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return errors.New("user not found")
	}

	projectMemberID, err := s.projectMemberService.GetProjectMemberID(userID, projectID)
	if err != nil {
		return errors.New("Project Member not found: " + err.Error())
	}

	updateTaskDTO.UpdatedBy = projectMemberID

	return s.taskRepository.EditTask(taskID, updateTaskDTO)
}

// TODO: Add logic so that the project owners and editors can delete it too
func (s *taskService) DeleteTask(deleteTaskDTO *models.DeleteTaskDTO) error {

	userID, err := s.userService.GetUserIDByFirebaseUID(deleteTaskDTO.FirebaseUID)
	if err != nil {
		log.Println("Failed to get UserID from firebaseUID")
		return errors.New("failted to  get UserID from firebaseUID")
	}

	projectMemberID, err := s.projectMemberService.GetProjectMemberID(userID, deleteTaskDTO.ProjectID)
	if err != nil {
		log.Println("Failed to get project member id")
		return errors.New("failed to get project member id")
	}

	var taskResponseDTO *models.TaskResponseDTO
	taskResponseDTO, err = s.GetTaskByID(deleteTaskDTO.TaskID)
	if err != nil {
		log.Println("Failed to get task")
		return errors.New("failed to get task")
	}

	if taskResponseDTO.CreatedBy.ID == projectMemberID {
		return s.taskRepository.DeleteTask(deleteTaskDTO.TaskID)
	} else {
		log.Println("Not the creator of comment. failed to delete")
		return errors.New("not the creator of comment. failed to delete")
	}

}
