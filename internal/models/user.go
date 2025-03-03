package models

import "github.com/google/uuid"

type CreateUserDTO struct {
	FirebaseUID string `json:"firebaseUID"`
	Name        string `json:"name"`
	Email       string `json:"email"`
}

type UserResponseDTO struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
