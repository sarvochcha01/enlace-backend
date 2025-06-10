package services

import (
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
)

type DashboardService interface {
	GetRecentlyAssignedTasks(firebaseUID string, limit int) ([]models.TaskResponseDTO, error)
	GetInProgressTasks(firebaseUID string, limit int) ([]models.TaskResponseDTO, error)
	GetApproachingDeadlineTasks(firebaseUID string, limit int) ([]models.TaskResponseDTO, error)

	Search(firebaseUID string, query string) (*models.SearchResult, error)
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

func (s *dashboardService) Search(firebaseUID string, query string) (*models.SearchResult, error) {

	userID, err := s.userService.GetUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}

	projectsChan := make(chan []models.ProjectSearchResult, 1)
	projectsErrChan := make(chan error, 1)

	go func() {
		projects, err := s.dashboardRepository.SearchProjects(userID, query)
		if err != nil {
			projectsErrChan <- err
			return
		}
		projectsChan <- projects
	}()

	tasksChan := make(chan []models.TaskResponseDTO, 1)
	tasksErrChan := make(chan error, 1)

	go func() {
		tasks, err := s.dashboardRepository.SearchTasks(userID, query)
		if err != nil {
			tasksErrChan <- err
			return
		}
		tasksChan <- tasks
	}()

	var projects []models.ProjectSearchResult
	var tasks []models.TaskResponseDTO

	select {
	case projects = <-projectsChan:
	case err = <-projectsErrChan:
		if err != nil {
			return nil, err
		}
	}

	select {
	case tasks = <-tasksChan:
	case err = <-tasksErrChan:
		if err != nil {
			return nil, err
		}
	}

	result := &models.SearchResult{
		Projects: projects,
		Tasks:    tasks,
	}

	// If either slice is nil, initialize as empty slice for consistent JSON response
	if result.Projects == nil {
		result.Projects = []models.ProjectSearchResult{}
	}
	if result.Tasks == nil {
		result.Tasks = []models.TaskResponseDTO{}
	}

	return result, nil
}
