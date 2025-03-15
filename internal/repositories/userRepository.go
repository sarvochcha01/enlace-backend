package repositories

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
)

type UserRepository interface {
	CreateUser(userDTO *models.CreateUserDTO) error
	GetUserIDByFirebaseUID(firebaseUID string) (uuid.UUID, error)
	GetUserByFirebaseUID(firebaseUID string) (*models.UserResponseDTO, error)
	SearchUsers(searchQuery string) ([]models.UserResponseDTO, error)
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

func (r *userRepository) GetUserIDByFirebaseUID(firebaseUID string) (uuid.UUID, error) {
	var userID uuid.UUID
	query := "SELECT id FROM users WHERE firebase_uid = $1"

	err := r.db.QueryRow(query, firebaseUID).Scan(&userID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

func (r *userRepository) GetUserByFirebaseUID(firebaseUID string) (*models.UserResponseDTO, error) {

	var user models.UserResponseDTO
	queryString := `
		SELECT id, email, name
		FROM users
		WHERE firebase_uid = $1
	`

	err := r.db.QueryRow(queryString, firebaseUID).Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (r *userRepository) SearchUsers(searchQuery string) ([]models.UserResponseDTO, error) {

	queryString := `
		SELECT id, email, name
		FROM users
		WHERE email ILIKE $1 OR name ILIKE $1
		ORDER BY name ASC
		LIMIT 5
	`

	rows, err := r.db.Query(queryString, "%"+searchQuery+"%")
	if err != nil {
		return nil, errors.New("failed to search users")
	}

	defer rows.Close()

	users := []models.UserResponseDTO{}
	for rows.Next() {
		var user models.UserResponseDTO
		if err := rows.Scan(&user.ID, &user.Email, &user.Name); err != nil {
			return nil, errors.New("error scanning user data")
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New("error iterating over search results")
	}

	return users, nil

}
