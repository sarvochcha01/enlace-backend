package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string
type NotificationStatus string

const (
	NotificationTypeTaskAssigned      NotificationType = "task_assigned"
	NotificationTypeProjectInvitation NotificationType = "project_invitation"
	NotificationTypeCommentAdded      NotificationType = "comment_added"

	NotificationStatusUnread NotificationStatus = "unread"
	NotificationStatusRead   NotificationStatus = "read"
)

type NotificationResponseDTO struct {
	ID        uuid.UUID          `json:"id"`
	UserID    uuid.UUID          `json:"userId"`
	Type      NotificationType   `json:"type"`
	Content   string             `json:"content"`
	ProjectID uuid.UUID          `json:"projectId"`
	TaskID    uuid.UUID          `json:"taskId"`
	Status    NotificationStatus `json:"status"`
	CreatedAt time.Time          `json:"createdAt"`
}

type CreateNotificationDTO struct {
	UserID    uuid.UUID        `json:"user_id"`
	Type      NotificationType `json:"type"`
	Content   string           `json:"content"`
	ProjectID uuid.UUID        `json:"projectId"`
	TaskID    uuid.UUID        `json:"taskId"`
}
