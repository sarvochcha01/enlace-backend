package services

import (
	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
)

type UserService interface {
	CreateUser(*models.CreateUserDTO) error
	FindUserIDByFirebaseUID(string) (uuid.UUID, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) UserService {
	return &userService{userRepository: r}
}

func (s *userService) CreateUser(userDTO *models.CreateUserDTO) error {
	return s.userRepository.CreateUser(userDTO)
}

func (s *userService) FindUserIDByFirebaseUID(firebaseUID string) (uuid.UUID, error) {
	return s.userRepository.FindUserIDByFirebaseUID(firebaseUID)
}
