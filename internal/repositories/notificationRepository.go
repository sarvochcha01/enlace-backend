package repositories

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
)

type NotificationRepository interface {
	CreateNotification(createNotificationDTO models.CreateNotificationDTO) (*models.NotificationResponseDTO, error)
	GetNotification(notificationID uuid.UUID) (*models.NotificationResponseDTO, error)
	GetAllNotificationsForUser(userID uuid.UUID) ([]models.NotificationResponseDTO, error)
	MarkNotificationAsRead(notificationID uuid.UUID) error
}

type notificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) CreateNotification(createNotificationDTO models.CreateNotificationDTO) (*models.NotificationResponseDTO, error) {
	var notification models.NotificationResponseDTO

	queryString := `
		INSERT INTO notifications
		(user_id, type, content, related_project_id, related_task_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, type, content, related_project_id, related_task_id, status, created_at
	`

	err := r.db.QueryRow(
		queryString,
		createNotificationDTO.UserID,
		createNotificationDTO.Type,
		createNotificationDTO.Content,
		createNotificationDTO.ProjectID,
		createNotificationDTO.TaskID,
	).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Type,
		&notification.Content,
		&notification.ProjectID,
		&notification.TaskID,
		&notification.Status,
		&notification.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &notification, nil
}

func (r *notificationRepository) GetAllNotificationsForUser(userID uuid.UUID) ([]models.NotificationResponseDTO, error) {
	notifications := []models.NotificationResponseDTO{}

	queryString := `
		SELECT 
			n.id, 
			n.user_id, 
			n.type, 
			n.content, 
			n.related_project_id, 
			n.related_task_id, 
			n.status, 
			n.created_at,
			i.id as invitation_id
		FROM notifications n
		LEFT JOIN invitations i
			ON n.type = 'project_invitation' 
			AND i.project_id = n.related_project_id 
			AND i.invited_user_id = n.user_id
		WHERE n.user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(queryString, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification models.NotificationResponseDTO
		if err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&notification.Content,
			&notification.ProjectID,
			&notification.TaskID,
			&notification.Status,
			&notification.CreatedAt,
			&notification.InvitationID,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *notificationRepository) GetNotification(notificationID uuid.UUID) (*models.NotificationResponseDTO, error) {

	var notification models.NotificationResponseDTO

	queryString := `
		SELECT id, user_id, type, content, related_project_id, related_task_id, status, created_at
		FROM notifications
		WHERE id = $1
		ORDER BY created_at DESC
	`

	if err := r.db.QueryRow(queryString, notificationID).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Type,
		&notification.Content,
		&notification.ProjectID,
		&notification.TaskID,
		&notification.Status,
		&notification.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &notification, nil
}

func (r *notificationRepository) MarkNotificationAsRead(notificationID uuid.UUID) error {
	queryString := `
		UPDATE notifications
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(queryString, models.NotificationStatusRead, notificationID)
	if err != nil {
		return err
	}

	return nil
}
