package services

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
	"github.com/sarvochcha01/enlace-backend/internal/utils"
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

	projectMember, err := s.projectMemberService.GetProjectMemberByUserID(userID, taskDTO.ProjectID)
	if err != nil {
		return uuid.Nil, errors.New("Project Member not found: " + err.Error())
	}

	if !utils.HasEditPrivileges(projectMember) {
		return uuid.Nil, errors.New("no edit privilege")
	}

	taskDTO.CreatedBy = projectMember.ID
	taskDTO.UpdatedBy = projectMember.ID

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

	projectMember, err := s.projectMemberService.GetProjectMemberByUserID(userID, projectID)
	if err != nil {
		return errors.New("Project Member not found: " + err.Error())
	}

	if !utils.HasEditPrivileges(projectMember) {
		return errors.New("no edit privilege")
	}
	updateTaskDTO.UpdatedBy = projectMember.ID

	if updateTaskDTO.AssignedTo != nil {

		currentTask, err := s.taskRepository.GetTaskByID(taskID)
		if err != nil {
			return errors.New("Failed to get current task: " + err.Error())
		}

		if currentTask.AssignedTo == nil || currentTask.AssignedTo.ID != *updateTaskDTO.AssignedTo {
			var assignedToUserID uuid.UUID

			assignedToUserID, err = s.projectMemberService.GetUserID(*updateTaskDTO.AssignedTo)

			if err == nil {
				if assignedToUserID != userID {
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
		}

	}

	return s.taskRepository.EditTask(taskID, updateTaskDTO)
}

func (s *taskService) DeleteTask(deleteTaskDTO *models.DeleteTaskDTO) error {

	projectMember, err := s.projectMemberService.GetProjectMemberByFirebaseUID(deleteTaskDTO.FirebaseUID, deleteTaskDTO.ProjectID)
	if err != nil {
		return errors.New("failed to get project member")
	}

	if !utils.HasEditPrivileges(projectMember) {
		return errors.New("no edit privilege")
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
