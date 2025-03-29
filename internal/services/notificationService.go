package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
	"github.com/sarvochcha01/enlace-backend/internal/websockets"
)

type NotificationService interface {
	CreateNotification(createNotificationDTO models.CreateNotificationDTO) error
	GetAllNotificationsForUser(firebaseUID string) ([]models.NotificationResponseDTO, error)
	GetNotification(notificationID uuid.UUID) (*models.NotificationResponseDTO, error)
	MarkNotificationAsRead(firebaseUID string, notificationID uuid.UUID) error
}

type notificationService struct {
	notificationRepository repositories.NotificationRepository
	wsHub                  *websockets.WebSocketHub
	userService            UserService
}

func NewNotificationService(nr repositories.NotificationRepository, wsHub *websockets.WebSocketHub, us UserService) NotificationService {
	return &notificationService{notificationRepository: nr, wsHub: wsHub, userService: us}
}

func (s *notificationService) CreateNotification(createNotificationDTO models.CreateNotificationDTO) error {
	notification, err := s.notificationRepository.CreateNotification(createNotificationDTO)
	if err != nil {
		return err
	}

	s.wsHub.SendNotificationToUser(*notification)

	return nil
}

func (s *notificationService) GetAllNotificationsForUser(firebaseUID string) ([]models.NotificationResponseDTO, error) {
	userID, err := s.userService.GetUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}

	return s.notificationRepository.GetAllNotificationsForUser(userID)
}

func (s *notificationService) GetNotification(notificationID uuid.UUID) (*models.NotificationResponseDTO, error) {
	return s.notificationRepository.GetNotification(notificationID)
}

func (s *notificationService) MarkNotificationAsRead(fireabseUID string, notificationID uuid.UUID) error {
	userID, err := s.userService.GetUserIDByFirebaseUID(fireabseUID)
	if err != nil {
		return err
	}

	notification, err := s.GetNotification(notificationID)
	if err != nil {
		return err
	}

	if userID != notification.UserID {
		return errors.New("only the receiver can update the notification status")
	}

	return s.notificationRepository.MarkNotificationAsRead(notificationID)

}
