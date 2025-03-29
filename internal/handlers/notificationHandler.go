package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/middlewares"
	"github.com/sarvochcha01/enlace-backend/internal/services"
)

// NotificationHandler struct and NewNotificationHandler (keep as is)
type NotificationHandler struct {
	notificationService services.NotificationService
	userService         services.UserService
}

func NewNotificationHandler(ns services.NotificationService, us services.UserService) *NotificationHandler {
	return &NotificationHandler{notificationService: ns, userService: us}
}

func (h *NotificationHandler) GetAllNotificationsForUser(w http.ResponseWriter, r *http.Request) {

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notifications, err := h.notificationService.GetAllNotificationsForUser(user.UID)
	if err != nil {
		log.Println("Failed to get notificaitons:", err)
		http.Error(w, "Failed to get notificaitons", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (h *NotificationHandler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {

	notificationIDStr := chi.URLParam(r, "notificationID")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		log.Println("Invalid notification ID:", err)
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	user, err := middlewares.GetFirebaseUser(r)
	if err != nil {
		log.Println("Unauthorized:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.notificationService.MarkNotificationAsRead(user.UID, notificationID)
	if err != nil {
		log.Println("Failed to get notificaitons:", err)
		http.Error(w, "Failed to get notificaitons", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification Read"))
}
