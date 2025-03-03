package repositories

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
)

type UserRepository interface {
	CreateUser(userDTO *models.CreateUserDTO) error
	FindUserIDByFirebaseUID(firebaseUID string) (uuid.UUID, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(userDTO *models.CreateUserDTO) error {
	queryString := `INSERT INTO users (firebase_uid, name, email) VALUES ($1, $2, $3)`

	_, err := r.db.Exec(queryString, userDTO.FirebaseUID, userDTO.Name, userDTO.Email)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) FindUserIDByFirebaseUID(firebaseUID string) (uuid.UUID, error) {
	var userID uuid.UUID
	query := "SELECT id FROM users WHERE firebase_uid = $1"

	err := r.db.QueryRow(query, firebaseUID).Scan(&userID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}
