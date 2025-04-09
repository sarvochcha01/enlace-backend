package models

import (
	"time"

	"github.com/google/uuid"
)

type InivtationStatus string

const (
	InivtationStatusPending  InivtationStatus = "pending"
	InivtationStatusAccepted InivtationStatus = "accepted"
	InivtationStatusDeclined InivtationStatus = "declined"
)

type CreateInvitationDTO struct {
	InvitedBy     uuid.UUID `json:"invitedBy"`
	InvitedUserID uuid.UUID `json:"invitedUserId"`
	ProjectID     uuid.UUID `json:"projectId"`
}

type InvitationResponseDTO struct {
	ID            uuid.UUID        `json:"id"`
	InvitedBy     uuid.UUID        `json:"invitedBy"`
	InvitedUserID uuid.UUID        `json:"invitedUserId"`
	ProjectID     uuid.UUID        `json:"projectId"`
	ProjectName   string           `json:"projectName"`
	Status        InivtationStatus `json:"status"`
	InvitedAt     time.Time        `json:"invitedAt"`
	Name          string           `json:"name"`
	Email         string           `json:"email"`
}

type EditInvitationDTO struct {
	Status       string    `json:"status"`
	InvitationID uuid.UUID `json:"id"`
	ProjectID    uuid.UUID `json:"projectId"`
}
