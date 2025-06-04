package services

import (
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
)

type DashboardService interface {
	GetRecentlyAssignedTasks(firebaseUID string, limit int) ([]models.TaskResponseDTO, error)
	GetInProgressTasks(firebaseUID string, limit int) ([]models.TaskResponseDTO, error)
	GetApproachingDeadlineTasks(firebaseUID string, limit int) ([]models.TaskResponseDTO, error)
}

type dashboardService struct {
	dashboardRepository repositories.DashboardRepository
	userService         UserService
}

func NewDashboardService(dr repositories.DashboardRepository, us UserService) DashboardService {
	return &dashboardService{
		dashboardRepository: dr,
		userService:         us,
	}
}

func (s *dashboardService) GetRecentlyAssignedTasks(firebaseUID string, limit int) ([]models.TaskResponseDTO, error) {
	userID, err := s.userService.GetUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}
	return s.dashboardRepository.GetRecentlyAssignedTasks(userID, limit)
}

func (s *dashboardService) GetInProgressTasks(firebaseUID string, limit int) ([]models.TaskResponseDTO, error) {
	userID, err := s.userService.GetUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}
	return s.dashboardRepository.GetInProgressTasks(userID, limit)
}

func (s *dashboardService) GetApproachingDeadlineTasks(firebaseUID string, limit int) ([]models.TaskResponseDTO, error) {
	userID, err := s.userService.GetUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}
	return s.dashboardRepository.GetApproachingDeadlineTasks(userID, limit)
}
