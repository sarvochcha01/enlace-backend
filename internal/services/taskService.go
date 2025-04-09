package services

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
)

type TaskService interface {
	CreateTask(taskDTo *models.CreateTaskDTO, firebaseUID string) (uuid.UUID, error)
	GetTaskByID(fireabseUID string, projectID uuid.UUID, taskID uuid.UUID) (*models.TaskResponseDTO, error)
	EditTask(uuid.UUID, uuid.UUID, string, *models.UpdateTaskDTO) error
	DeleteTask(*models.DeleteTaskDTO) error
}

type taskService struct {
	taskRepository       repositories.TaskRepository
	userService          UserService
	projectMemberService ProjectMemberService
	notificationService  NotificationService
}

func NewTaskService(tr repositories.TaskRepository, us UserService, pms ProjectMemberService, ns NotificationService) TaskService {
	return &taskService{taskRepository: tr, userService: us, projectMemberService: pms, notificationService: ns}
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

func (s *taskService) GetTaskByID(firebaseUID string, projectID uuid.UUID, taskID uuid.UUID) (*models.TaskResponseDTO, error) {
	_, err := s.projectMemberService.GetProjectMemberByFirebaseUID(firebaseUID, projectID)
	if err != nil {
		return nil, nil
	}
	return s.taskRepository.GetTaskByID(taskID)
}

func (s *taskService) GetTaskByIDNoAuth(taskID uuid.UUID) (*models.TaskResponseDTO, error) {
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

	if updateTaskDTO.AssignedTo != nil {
		var assignedToUserID uuid.UUID

		assignedToUserID, err = s.projectMemberService.GetUserID(*updateTaskDTO.AssignedTo)
		if err == nil {
			notification := &models.CreateNotificationDTO{
				UserID:    assignedToUserID,
				Type:      models.NotificationTypeTaskAssigned,
				Content:   fmt.Sprintf("You have been assigned to task: %s", updateTaskDTO.Title),
				ProjectID: projectID,
				TaskID:    &taskID,
			}
			err = s.notificationService.CreateNotification(*notification)
			if err != nil {
				fmt.Println("Failed to create notification:", err)
			}
		}

	}

	return s.taskRepository.EditTask(taskID, updateTaskDTO)
}

func (s *taskService) DeleteTask(deleteTaskDTO *models.DeleteTaskDTO) error {

	projectMember, err := s.projectMemberService.GetProjectMemberByFirebaseUID(deleteTaskDTO.FirebaseUID, deleteTaskDTO.ProjectID)
	if err != nil {
		return errors.New("failed to get project member")
	}

	task, err := s.GetTaskByIDNoAuth(deleteTaskDTO.TaskID)
	if err != nil {
		return errors.New("failed to get task")
	}

	if projectMember.Role != models.RoleOwner && task.CreatedBy.ID != projectMember.ID {
		return errors.New("only the owner or task creator can delete this task")
	}
	return s.taskRepository.DeleteTask(deleteTaskDTO.TaskID)
}
